package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os/exec"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`rm|shutdown|poweroff|init|mkfs|dd|mv|curl`)

// Command represents an external command being prepared or run.
type Command struct {
	Name string
	Args []string
}

type Result struct {
	Stdout string
	Stderr string
	Err    error
}

type ShellExecutor struct{}

func (*ShellExecutor) Exec(args Command, reply *Result) error {
	if match := re.FindString(args.Name); match != "" {
		return fmt.Errorf("%q is not allowed", args.Name)
	}

	cmd := exec.Command(args.Name, args.Args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		reply.Err = err
		return err
	}

	reply.Stdout = stdout.String()
	reply.Stderr = stderr.String()

	return nil
}

func main() {
	// security risk: anyone can execute restricted commands
	executor := new(ShellExecutor) // object aka service
	if err := rpc.Register(executor); err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	go http.Serve(ln, nil)

	client, err := rpc.DialHTTP("tcp", ":8090")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// args := Command{"init", []string{"0"}}
	args := Command{Name: "w"}
	result := new(Result)
	// if err := client.Call("ShellExecutor.Exec", args, result); err != nil {
	// 	log.Fatal("call ShellExecutor.Exec failed: ", err)
	// }
	call1 := client.Go("ShellExecutor.Exec", args, result, nil)
	call2 := <-call1.Done
	if call2.Error != nil {
		log.Fatal(call2.Error)
	}

	log.Println("stdout: ", result.Stdout)
	log.Println("sterr: ", result.Stderr)
}
