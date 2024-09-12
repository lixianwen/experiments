package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	directories      = []string{"AssetBundles", "GameBundles", "Release"}
	remoteBaseDir    = "/data/pack_bak"
	archiveThreshold = 30
)

type Client struct {
	SSHC  *ssh.Client
	SFTPC *sftp.Client
}

// Put puts the specified file to the remote server.
func (c *Client) Put(localPath string, remotePath string) error {
	lf, err := os.Open(localPath)
	if err != nil {
		slog.Error("os.Open", "err", err)
		return err
	}
	defer lf.Close()

	rf, err := c.SFTPC.Create(remotePath)
	if err != nil {
		slog.Error("client.Create", "err", err)
		return err
	}
	defer rf.Close()

	if _, err := rf.ReadFrom(lf); err != nil {
		slog.Error("*sftp.File.ReadFrom", "err", err)
		return err
	}

	return nil
}

// Prune removes the leftmost files or directories from the remote base directory
// which exceeded the archive threshold.
func (c *Client) Prune() error {
	sl, err := c.SFTPC.ReadDir(remoteBaseDir)
	if err != nil {
		return err
	}
	if sl == nil {
		return nil
	}

	slices.SortFunc(sl, func(a, b fs.FileInfo) int {
		return a.ModTime().Compare(b.ModTime())
	})

	threshold := len(sl) - archiveThreshold
	if threshold > 0 {
		for _, info := range sl[:threshold] {
			leaf := sftp.Join(remoteBaseDir, info.Name())
			slog.Info("Delete", "name", leaf)
			if err := c.SFTPC.RemoveAll(leaf); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run runs cmd on the remote host.
func (c *Client) Run(cmd string) (*ssh.Session, error) {
	session, err := c.SSHC.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var b1, b2 bytes.Buffer
	session.Stdout = &b1
	session.Stderr = &b2

	if err := session.Run(cmd); err != nil {
		return session, err
	}

	return session, nil
}

// Get copies from src to dst until either EOF is reached on src or an error occurs.
func (c *Client) Get(dst, src string) error {
	rf, err := c.SFTPC.Open(src)
	if err != nil {
		return err
	}
	defer rf.Close()

	lf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer lf.Close()

	// equal to call rf.WriteTo(lf)
	if _, err := io.Copy(lf, rf); err != nil {
		return err
	}

	return nil
}

func Backup(c *Client) error {
	prefix := sftp.Join(remoteBaseDir, time.Now().Format("20060102150405"))

	for _, bkitem := range directories {
		err := filepath.WalkDir(bkitem, func(p string, d fs.DirEntry, err error) error {
			// returning the error will cause WalkDir to stop walking the entire tree
			if err != nil {
				slog.Error("err param", "err", err)
				return err
			}

			newp := strings.ReplaceAll(p, "\\", "/")
			remotePath := sftp.Join(prefix, newp)
			if d.IsDir() {
				slog.Info("remote", "directory", remotePath)
				if err := c.SFTPC.MkdirAll(remotePath); err != nil {
					slog.Error("client.MkdirAll", "err", err)
					return err
				}
			} else {
				slog.Info("remote", "filename", remotePath)
				if err := c.Put(p, remotePath); err != nil {
					slog.Error("c.Put", "err", err)
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	parent, child := sftp.Split(prefix)
	if _, err := c.Run(fmt.Sprintf("cd %s && tar -zcf %s.tar.gz %s", parent, child, child)); err != nil {
		slog.Error("archive", "err", err)
		return err
	}
	slog.Info("cleanup", "rm", prefix)
	if err := c.SFTPC.RemoveAll(prefix); err != nil {
		slog.Error("cleanup", "rm", prefix, "err", err)
		return err
	}
	return nil
}

// Note: This concurrent version may be slower than the sequential version above.
func Backup2(c *Client) error {
	var wg sync.WaitGroup
	errCH := make(chan error)
	done := make(chan struct{}, len(directories))
	prefix := sftp.Join(remoteBaseDir, time.Now().Format("20060102150405"))

	for _, bkitem := range directories {
		wg.Add(1)
		go func(bkitem string) {
			defer wg.Done()

			err := filepath.WalkDir(bkitem, func(p string, d fs.DirEntry, err error) error {
				// returning the error will cause WalkDir to stop walking the entire tree
				if err != nil {
					slog.Error("err param", "err", err)
					return err
				}

				newp := strings.ReplaceAll(p, "\\", "/")
				remotePath := sftp.Join(prefix, newp)
				if d.IsDir() {
					slog.Info("remote", "directory", remotePath)
					if err := c.SFTPC.MkdirAll(remotePath); err != nil {
						slog.Error("client.MkdirAll", "err", err)
						return err
					}
				} else {
					slog.Info("remote", "filename", remotePath)
					if err := c.Put(p, remotePath); err != nil {
						slog.Error("c.Put", "err", err)
						return err
					}
				}

				return nil
			})

			if err != nil {
				errCH <- err
			} else {
				done <- struct{}{}
			}
		}(bkitem)
	}

	go func() {
		wg.Wait()
		rm := func() {
			slog.Info("cleanup", "rm", prefix)
			if err := c.SFTPC.RemoveAll(prefix); err != nil {
				slog.Warn("cleanup", "rm", prefix, "err", err)
			}
		}

		// one failed, others succeeded
		if len(done) != cap(done) {
			rm()
			return
		}

		parent, child := sftp.Split(prefix)
		if _, err := c.Run(fmt.Sprintf("cd %s && tar -zcf %s.tar.gz %s", parent, child, child)); err != nil {
			slog.Error("archive", "err", err)
			rm()
			errCH <- err
			return
		}

		rm()
		close(errCH)
	}()

	// short circuit asap when error arrived
	return <-errCH
}

func Download(c *Client) error {
	sl, err := c.SFTPC.ReadDir(remoteBaseDir)
	if err != nil {
		slog.Error("client.ReadDir", "err", err)
		return err
	}
	if sl == nil {
		return nil
	}

	var archives []fs.FileInfo
	for _, info := range sl {
		if !info.IsDir() {
			archives = append(archives, info)
		}
	}
	if len(archives) == 0 {
		slog.Info("没有可供下载的文件")
		return nil
	}
	slices.SortFunc(archives, func(a, b fs.FileInfo) int {
		return a.ModTime().Compare(b.ModTime())
	})

	count := func(x int) (c int) {
		if x == 0 {
			c++
			return
		}

		for ; x > 0; x /= 10 {
			c++
		}

		return
	}

	fmt.Println("请选择想要下载的备份序号:")
	for index, info := range archives {
		fmt.Printf("%d:%s%s\n", index, strings.Repeat(" ", count(len(archives))), info.Name())
	}

	var index int
	if _, err := fmt.Scanf("%d\n", &index); err != nil {
		slog.Error("fmt.Scanf", "err", err)
		return err
	}
	if index < 0 || index >= len(archives) {
		return fmt.Errorf("invalid seq number: %d", index)
	}

	name := archives[index].Name()
	slog.Info("download", "name", name)
	if err := c.Get(name, sftp.Join(remoteBaseDir, name)); err != nil {
		slog.Error("c.Get", "err", err)
		return err
	}

	return nil
}
