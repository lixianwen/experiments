package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			panic(err)
		}

		time.Sleep(3 * time.Second)

		// On a Unix-like system, pressing Ctrl+C on a keyboard sends a
		// SIGINT signal to the process of the program in execution.
		//
		// This example simulates that by sending a SIGINT signal to itself.
		if err := p.Signal(os.Interrupt); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	fmt.Println(ctx.Err())
	stop() // stop receiving signal notifications as soon as possible.
}
