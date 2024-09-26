package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	// As of Go 1.23, the garbage collector can recover unreferenced tickers,
	// even if they haven't been stopped.
	c := time.Tick(time.Second)
loop:
	for {
		select {
		case t := <-c:
			fmt.Println(t)
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			break loop
		}
	}
}
