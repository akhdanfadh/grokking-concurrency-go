package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// cpuWaster wastes the processor time, professionally
func cpuWaster(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker-%d doing its work\n", i)
	time.Sleep(3 * time.Second)
}

// displayThreads displays information about current process
func displayThreads(names []string) {
	fmt.Println("----------")
	fmt.Printf("Current process PID: %d\n", os.Getpid())
	fmt.Printf("Thread count: %d\n", runtime.NumGoroutine())
	fmt.Println("Active threads:")
	for _, n := range names {
		fmt.Printf("  %s\n", n)
	}
}

func main() {
	active := []string{"Main"}
	displayThreads(active)

	numThreads := 5
	fmt.Printf("Starting %d CPU wasters...\n", numThreads)

	var wg sync.WaitGroup
	for i := range numThreads {
		wg.Add(1)
		active = append(active, fmt.Sprintf("Worker-%d", i))
		go cpuWaster(i, &wg) // schedule a new goroutine (Go's threads abstraction)
	}

	time.Sleep(100 * time.Millisecond) // give time for scheduled goroutines to start
	displayThreads(active)
	wg.Wait() // prevent the main program to exit before all goroutines finish
}
