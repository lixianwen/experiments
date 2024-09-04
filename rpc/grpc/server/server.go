// It implements the shell executor service whose definition can be found in proto/shell_executor.proto.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"

	pb "demo/proto"

	"google.golang.org/grpc"
)

var re = regexp.MustCompile(`rm|shutdown|poweroff|init|mkfs|dd|mv|curl`)

type server struct {
	pb.UnimplementedShellExecutorServer
}

func (*server) Exec(ctx context.Context, args *pb.Command) (*pb.Result, error) {
	if match := re.FindString(args.Name); match != "" {
		return nil, fmt.Errorf("%q is not allowed", args.Name)
	}

	reply := new(pb.Result)

	cmd := exec.Command(args.Name, args.Args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		reply.Err = err.Error()
		return nil, err
	}

	reply.Stdout = stdout.String()
	reply.Stderr = stderr.String()

	return reply, nil
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", ln.Addr())

	s := grpc.NewServer()
	pb.RegisterShellExecutorServer(s, &server{})

	go func() {
		if err := s.Serve(ln); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	signal.Stop(c)
	s.Stop()
	log.Println("program exit")
}
