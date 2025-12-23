package main

import (
	"sync"
)

func runHierarchy() {
	dumplings := 20

	a := NewLockWithName("Chopstick A")
	b := NewLockWithName("Chopstick B")
	// changing the order a > b, so a should be acquired first
	p1 := &Philosopher{"Philosopher #1", a, b, &dumplings}
	p2 := &Philosopher{"Philosopher #2", a, b, &dumplings}

	var wg sync.WaitGroup
	wg.Add(2)
	go p1.Run(&wg)
	go p2.Run(&wg)
	wg.Wait()
}
