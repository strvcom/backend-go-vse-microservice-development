package main

import (
	"fmt"
	"time"
)

func main() {
	go work("doors")
	go work("windows")
	go work("cables")

	// wait for everyone to finish
	time.Sleep(5 * time.Second)
}

func work(task string) {
	for i := 0; i < 3; i++ {
		fmt.Println("working on", task, "day number", i+1)
		time.Sleep(time.Second)
	}
}
