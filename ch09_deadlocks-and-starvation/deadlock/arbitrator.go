package main

import (
	"fmt"
	"sync"
	"time"
)

type Waiter struct {
	mu sync.Mutex
}

func (w *Waiter) AskForChopsticks(left, right *LockWithName) {
	w.mu.Lock()
	left.Acquire()
	fmt.Printf("%s grabbed\n", left.Name)
	right.Acquire()
	fmt.Printf("%s grabbed\n", right.Name)
	w.mu.Unlock()
}

func (w *Waiter) ReleaseChopsticks(left, right *LockWithName) {
	right.Release()
	fmt.Printf("%s released\n", right.Name)
	left.Release()
	fmt.Printf("%s released\n\n", left.Name)
}

type PhilosopherWaiter struct {
	Philosopher *Philosopher
	Waiter      *Waiter
}

func (pw *PhilosopherWaiter) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for *pw.Philosopher.Dumplings > 0 {
		fmt.Printf("%s asks waiter for chopsticks\n", pw.Philosopher.Name)
		pw.Waiter.AskForChopsticks(pw.Philosopher.Left, pw.Philosopher.Right)

		*pw.Philosopher.Dumplings--
		fmt.Printf("%s eats a dumpling. Dumplings left: %d\n", pw.Philosopher.Name, *pw.Philosopher.Dumplings)

		fmt.Printf("%s asks waiter for chopsticks\n", pw.Philosopher.Name)
		pw.Waiter.ReleaseChopsticks(pw.Philosopher.Left, pw.Philosopher.Right)

		time.Sleep(100 * time.Millisecond)
	}
}

func runArbitrator() {
	dumplings := 20

	a := NewLockWithName("Chopstick A")
	b := NewLockWithName("Chopstick B")
	w := &Waiter{}
	p1 := &Philosopher{"Philosopher #1", a, b, &dumplings}
	p2 := &Philosopher{"Philosopher #2", b, a, &dumplings}
	pw1 := &PhilosopherWaiter{p1, w}
	pw2 := &PhilosopherWaiter{p2, w}

	var wg sync.WaitGroup
	wg.Add(2)
	go pw1.Run(&wg)
	go pw2.Run(&wg)
	wg.Wait()
}
