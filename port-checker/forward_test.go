package main

import (
	"net"
	"testing"
)

// This test would require a legitimate, functioning SSH server and MySQL server to connect to,
// and it's closer to an integration test because it uses a real external resource.
func TestProxy(t *testing.T) {
	config, err := NewClientConfigForKey("/root/.ssh/id_rsa", "root")
	if err != nil {
		t.Error(err)
	}

	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		t.Error(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		go forwarder(conn, "localhost:22", "localhost:3306", config)
	}()

	db, err := NewDB("mysql", "localhost", "8090", "root", "123456", "ht", false)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	// arrange data to be tested based on you purpose
	disconnected, err := checkConn(db, "27001_27004")
	if err != nil {
		t.Error(err)
	}
	t.Log("disconnected", disconnected)
	// modify this assertion based on your data
	if disconnected == nil {
		t.Error("It should be disconnected")
	}
}
