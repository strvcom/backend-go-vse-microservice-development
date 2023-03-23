package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

type (
	ctxKeyCustomer struct{}
)

var (
	contextKey = struct {
		customer ctxKeyCustomer
	}{}
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey.customer, "Gopher")

	doors := make(chan string)
	windows := make(chan string)
	cables := make(chan string)

	go doorsFactory(ctx, doors)
	go windowsFactory(ctx, windows)
	go cablesFactory(context.Background(), cables) // we do not pass a context with value

	// infinite loop
	for {
		select {
		case product := <-doors:
			fmt.Println("got", product)
		case product := <-windows:
			fmt.Println("got", product)
		case product := <-cables:
			fmt.Println("got", product)
		case <-time.After(2 * time.Second): // break the loop after 2 seconds
			os.Exit(0)
		}
	}
}

func doorsFactory(ctx context.Context, doors chan<- string) {
	printCustomer(ctx)
	for i := 0; i < 5; i++ {
		doors <- "doors"
	}
}

func windowsFactory(ctx context.Context, windows chan<- string) {
	printCustomer(ctx)
	for i := 0; i < 5; i++ {
		windows <- "windows"
	}
}

func cablesFactory(ctx context.Context, cables chan<- string) {
	printCustomer(ctx)
	for i := 0; i < 5; i++ {
		cables <- "cables"
	}
}

func printCustomer(ctx context.Context) {
	if v := ctx.Value(contextKey.customer); v != nil {
		fmt.Println("name of our dear customer :", v)
		return
	}
	fmt.Println("customer name not found")
}
