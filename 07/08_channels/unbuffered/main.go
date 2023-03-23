package main

import (
	"fmt"
	"os"
	"time"
)

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
		// wait before getting next product
		time.Sleep(time.Second)
	}
}

func doorsFactory(doors chan<- string) {
	for i := 0; i < 5; i++ {
		fmt.Println("doors made")
		doors <- "doors"
		time.Sleep(time.Millisecond)
	}
}

func windowsFactory(windows chan<- string) {
	for i := 0; i < 5; i++ {
		fmt.Println("windows made")
		windows <- "windows"
		time.Sleep(time.Millisecond)
	}
}

func cablesFactory(cables chan<- string) {
	for i := 0; i < 5; i++ {
		fmt.Println("cables made")
		cables <- "cables"
		time.Sleep(time.Millisecond)
	}
}
