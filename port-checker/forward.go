package main

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

// forwarder establishes a ssh tunnel between a remote server(proxy) and internal server.
func forwarder(client net.Conn, remoteAddr, privateAddr string, config *ssh.ClientConfig) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover panic[%s] and exit", err)
		}
		client.Close()
	}()

	remoteConn, err := ssh.Dial("tcp", remoteAddr, config)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	rw, err := remoteConn.Dial("tcp", privateAddr)
	if err != nil {
		return
	}
	defer rw.Close()

	done := make(chan bool)

	go func() {
		defer func() { done <- true }()
		if _, err := io.Copy(rw, client); err != nil {
			log.Println("client -> rw", err)
		}
	}()

	go func() {
		defer func() { done <- true }()
		if _, err := io.Copy(client, rw); err != nil {
			log.Println("rw -> client", err)
		}
	}()

	<-done
}
