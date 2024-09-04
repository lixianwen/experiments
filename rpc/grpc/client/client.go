package main

import (
	"context"
	"log"

	pb "demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// set the auth credentials
	// var opts []grpc.DialOption

	// a gRPC channel
	conn, err := grpc.NewClient("0.0.0.0:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// a client stub
	client := pb.NewShellExecutorClient(conn)

	// calling service's method
	// args := &pb.Command{Name: "w"}
	args := &pb.Command{Name: "top", Args: []string{"-bn1"}}
	result, err := client.Exec(context.Background(), args)
	if err != nil {
		log.Fatalf("failed to invoke: %v", err)
	}
	log.Printf("%+v", result)
}
