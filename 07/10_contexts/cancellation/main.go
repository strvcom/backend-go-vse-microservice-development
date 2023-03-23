package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

var (
	contextTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	doors := make(chan string)
	windows := make(chan string)
	cables := make(chan string)

	go doorsFactory(ctxWithTimeout, doors)
	go windowsFactory(ctxWithTimeout, windows)
	go cablesFactory(ctxWithTimeout, cables)

	// infinite loop
	for {
		select {
		case product := <-doors:
			fmt.Println("got", product)
		case product := <-windows:
			fmt.Println("got", product)
		case product := <-cables:
			fmt.Println("got", product)
		case <-ctxWithTimeout.Done():
			time.Sleep(time.Second)
			os.Exit(0)
		case <-time.After(10 * time.Second): // break the loop after 10 seconds
			os.Exit(0)
		}
	}
}

func doorsFactory(ctx context.Context, doors chan<- string) {
	for {
		select {
		case <-time.After(time.Second):
			// make doors each second
			doors <- "doors"
		case <-ctx.Done():
			// end production if context was cancelled
			fmt.Println("ending doors production")
			return
		}
	}
}

func windowsFactory(ctx context.Context, windows chan<- string) {
	for {
		select {
		case <-time.After(time.Second):
			// make windows each second
			windows <- "windows"
		case <-ctx.Done():
			// end production if context was cancelled
			fmt.Println("ending windows production")
			return
		}
	}
}

func cablesFactory(ctx context.Context, cables chan<- string) {
	for {
		select {
		case <-time.After(time.Second):
			// make cables each second
			cables <- "cables"
		case <-ctx.Done():
			// end production if context was cancelled
			fmt.Println("ending cables production")
			return
		}
	}
}
