package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nExample: %s -dbip 127.0.0.1 -mip 172.20.40.129 -sids 1001:1002\n", os.Args[0])
	}
	dbip := flag.String("dbip", "", "API database server private ip address")
	mip := flag.String("mip", "", "Master server public ip address")
	sids := flag.String("sids", "", "service id list")
	flag.Parse()
	if flag.NFlag() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	config, err := NewClientConfigForKey("/root/.ssh/id_rsa", "root")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	errCH := make(chan error)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCH <- err
			return
		}
		go forwarder(conn, net.JoinHostPort(*mip, "22"), net.JoinHostPort(*dbip, "3306"), config)
	}()

	db, err := NewDB("mysql", "localhost", "8090", "root", "123456", "ht", false)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	select {
	case err := <-errCH:
		panic(err)
	default:
		disconnected, err := checkConn(db, *sids)
		if err != nil {
			panic(err)
		}
		if err := json.NewEncoder(os.Stdout).Encode(disconnected); err != nil {
			panic(err)
		}
	}
}
