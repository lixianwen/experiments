package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func slowOperation(ctx context.Context) {
	duration := time.Duration(rand.Intn(10)) * time.Second
	select {
	case <-time.After(duration):
		fmt.Println("timer fired")
	case <-ctx.Done():
		fmt.Println("timed out")
	}
}

func slowOperationWithTimeout(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	slowOperation(ctx)
}

func main() {
	slowOperationWithTimeout(context.Background())
}
