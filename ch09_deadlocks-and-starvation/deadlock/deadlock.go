package main

import (
	"sync"
)

func runDeadlock() {
	dumplings := 20

	a := NewLockWithName("Chopstick A")
	b := NewLockWithName("Chopstick B")
	p1 := &Philosopher{"Philosopher #1", a, b, &dumplings}
	p2 := &Philosopher{"Philosopher #2", b, a, &dumplings}

	var wg sync.WaitGroup
	wg.Add(2)
	go p1.Run(&wg)
	go p2.Run(&wg)
	wg.Wait()
}
