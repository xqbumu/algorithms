package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-time.NewTimer(time.Second).C
		cancel()
	}()
	go background(ctx, "A")
	background(ctx, "B")
}

func background(ctx context.Context, val string) {
	for {
		select {
		case <-ctx.Done():
			println(val)
			return
		}
	}
}
