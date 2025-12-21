package main

import (
	"fmt"
	"sync"
	"time"
)

const SIZE = 5

var sharedMemory [SIZE]int

func init() {
	for i := range SIZE {
		sharedMemory[i] = -1
	}
}

func runProducer(wg *sync.WaitGroup) {
	defer wg.Done()

	name := "Producer"
	for i := range SIZE {
		fmt.Printf("%s: Writing: %d\n", name, i)
		idx := (i - 1 + SIZE) % SIZE // index wrapping
		sharedMemory[idx] = i
	}
}

func runConsumer(wg *sync.WaitGroup) {
	defer wg.Done()

	name := "Consumer"
	for i := range SIZE {
		// try reading the data until succession
		for {
			line := sharedMemory[i]
			if line == -1 {
				// data hasn't change - waiting for a second
				fmt.Printf("%s: Data not available\n", name)
				fmt.Println("Sleeping for 1 second before retrying")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("%s: Read: %d\n", name, line)
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go runConsumer(&wg)
	go runProducer(&wg)

	wg.Wait() // prevent main from exiting before goroutines finish
}
