package main

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"slices"
	"strings"

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
	filenames := make([]string, 0, len(directories))
	for _, dir := range directories {
		filenames = append(filenames, fmt.Sprintf("%s.tar.gz", dir))
	}
	defer Remove(filenames...)

	src, err := Archive(directories...)
	defer Remove(src)
	if err != nil {
		return err
	}

	if err := c.Put(src, sftp.Join(remoteBaseDir, src)); err != nil {
		return err
	}

	return nil
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
