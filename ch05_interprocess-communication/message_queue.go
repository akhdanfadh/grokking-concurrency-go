//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

const threadNum = 4

type worker struct {
	id int
	q  <-chan int // receive-only: worker only consumes from the queue
}

func (w *worker) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range w.q {
		fmt.Printf("Thread %d: processing item %d from the queue\n", w.id, item)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	// creates a queue with values to put into it for processing in the threads
	// buffered to hold all initial messages (like python's queue)
	q := make(chan int, 10)
	for i := range 10 {
		q <- i
	}
	close(q) // no more messages will be produced

	// run threads to process data from queue
	var wg sync.WaitGroup
	wg.Add(threadNum)
	for i := range threadNum {
		w := &worker{i + 1, q}
		go w.run(&wg)
	}
	wg.Wait() // block main threads until all workers finished
}
