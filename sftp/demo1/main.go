package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	BACKUP = iota
	DOWNLOAD
)

// storage server
var (
	host     string
	port     string
	username string
	password string
)

func init() {
	flag.StringVar(&host, "h", "172.20.40.129", "storage server's hostname")
	flag.StringVar(&port, "P", "22", "storage server's ssh port")
	flag.StringVar(&username, "u", "root", "storage server's ssh name")
	flag.StringVar(&password, "p", "123456", "storage server's ssh password")
}

func main() {
	flag.Parse()

	logFile, err := os.OpenFile("backuplog.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// h := slog.NewTextHandler(logFile, &slog.HandlerOptions{AddSource: true})
	h := slog.NewTextHandler(logFile, nil)
	slog.SetDefault(slog.New(h))

	prompt := `------------------------------------------------------------
请勿手动关闭该程序, 运行时间较长, 请耐心等待
当此终端退出时，请检查当前目录下的日志文件 'backuplog.log'"
如果最后一行包含 'Done.' 字样，视为备份成功，否则请联系运维"

请选择操作序号:
0:     备份
1:     下载备份
------------------------------------------------------------`
	fmt.Println(prompt)

	var index int
	if _, err := fmt.Scanf("%d\n", &index); err != nil {
		log.Fatal(err)
	}
	if index != BACKUP && index != DOWNLOAD {
		log.Fatalf("invalid seq number: %d", index)
	}
	fmt.Printf("正在备份 {%s} 到目标服务器 %s\n", strings.Join(directories, ","), host)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	myclient := &Client{
		SSHC:  conn,
		SFTPC: client,
	}

	if err := myclient.Prune(); err != nil {
		log.Fatal(err)
	}

	switch index {
	case BACKUP:
		if err := Backup(myclient); err != nil {
			log.Fatal(err)
		}
	case DOWNLOAD:
		if err := Download(myclient); err != nil {
			log.Fatal(err)
		}
	}

	// successful flag
	log.Println("Done.")
}
