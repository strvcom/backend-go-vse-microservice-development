package main

import (
	"fmt"
	"sync"
)

func main() {
	mutex := &sync.Mutex{}
	var count int

	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			count++
			mutex.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
