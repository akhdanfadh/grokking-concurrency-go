package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	counter int      = 0 // shared memory
	lock    RWLocker = NewRWLock()
	// lock RWLocker = NewRWLockFair()
)

func user(idx int) {
	for {
		lock.AcquireRead()

		fmt.Printf("User %d reading '%d'\n", idx, counter)
		time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)

		lock.ReleaseRead()
		fmt.Printf("User %d done reading\n", idx)
		time.Sleep(500 * time.Millisecond) // simulate time between reads
	}
}

func librarian() {
	for {
		lock.AcquireWrite()

		fmt.Print("Librarian writing... ")
		time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second) // simulate writing time
		counter++
		fmt.Printf("New value: %d\n", counter)

		lock.ReleaseWrite()
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup
	wg.Add(3)
	go func() { defer wg.Done(); user(0) }()
	go func() { defer wg.Done(); user(1) }()
	go func() { defer wg.Done(); librarian() }()
	wg.Wait()
}
