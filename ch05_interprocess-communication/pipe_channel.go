//go:build ignore

package main

import (
	"fmt"
	"sync"
)

type writer struct {
	ch chan<- string // send-only channel
}

func (w *writer) run(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Writer: Sending rubber duck...")
	w.ch <- "Rubber duck"
}

type reader struct {
	ch <-chan string // receive-only channel
}

func (r *reader) run(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Reader: Reading...")
	msg := <-r.ch
	fmt.Printf("Reader: Received: %s\n", msg)
}

func main() {
	// a channel can be thought of as a pipe-like conduit between goroutines
	// buffer size 1 makes this a simple "one message in flight" demo
	ch := make(chan string, 1)

	writer := &writer{ch: ch}
	reader := &reader{ch: ch}

	var wg sync.WaitGroup
	wg.Add(2)
	go writer.run(&wg)
	go reader.run(&wg)
	wg.Wait() // block main until child goroutines finish
}
