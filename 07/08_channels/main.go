package main

import (
	"fmt"
	"os"
	"time"
)

// When you run this code with `go run main.go`, you can notice the random order of doors, windows and cables.
// When you run it with `GOMAXPROCS=1 go run main.go`, you can see that only one goroutine at a time is executed.

func main() {
	doors := make(chan string)
	windows := make(chan string)
	cables := make(chan string)

	go doorsFactory(doors)
	go windowsFactory(windows)
	go cablesFactory(cables)

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

func doorsFactory(doors chan<- string) {
	for i := 0; i < 5; i++ {
		doors <- "doors"
	}
}

func windowsFactory(windows chan<- string) {
	for i := 0; i < 5; i++ {
		windows <- "windows"
	}
}

func cablesFactory(cables chan<- string) {
	for i := 0; i < 5; i++ {
		cables <- "cables"
	}
}
