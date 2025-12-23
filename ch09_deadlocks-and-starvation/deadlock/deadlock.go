package main

import (
	"fmt"
	"sync"
	"time"
)

// startKeepAlive starts a goroutine that periodically print a message.
// This prevents the Go runtime from declaring a global deadlock in example.
// It does not resolve deadlock; it only keeps the program alive to observe it.
func startKeepAlive(prefix string, every time.Duration) (stop func()) {
	done := make(chan struct{})
	go func() {
		t := time.NewTicker(every)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				fmt.Printf("%s still waiting...\n", prefix)
			case <-done:
				return
			}
		}
	}()
	return func() { close(done) }
}

func runDeadlock() {
	dumplings := 20

	a := NewLockWithName("Chopstick A")
	b := NewLockWithName("Chopstick B")
	p1 := &Philosopher{"Philosopher #1", a, b, &dumplings}
	p2 := &Philosopher{"Philosopher #2", b, a, &dumplings}

	stop := startKeepAlive("(deadlock demo)", 500*time.Millisecond)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)
	go p1.Run(&wg)
	go p2.Run(&wg)
	wg.Wait()
}
