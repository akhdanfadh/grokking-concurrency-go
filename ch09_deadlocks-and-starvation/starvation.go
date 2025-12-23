//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

type LockWithName struct {
	Name string
	mu   sync.Mutex
}

func NewLockWithName(name string) *LockWithName { return &LockWithName{Name: name} }
func (l *LockWithName) Acquire()                { l.mu.Lock() }
func (l *LockWithName) Release()                { l.mu.Unlock() }
func (l *LockWithName) TryAcquire() bool        { return l.mu.TryLock() }

type Philosopher struct {
	Name        string
	Left, Right *LockWithName
	Dumplings   *int // shared resource
}

func (p *Philosopher) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	dumplingsEaten := 0
	for *p.Dumplings > 0 {
		p.Left.Acquire()
		p.Right.Acquire()

		if *p.Dumplings <= 0 {
			p.Right.Release()
			p.Left.Release()
			break
		}

		*p.Dumplings--
		dumplingsEaten++

		time.Sleep(time.Nanosecond) // should be fast enough for "yield right now" action

		p.Right.Release()
		p.Left.Release()
	}
	fmt.Printf("%s took %d pieces\n", p.Name, dumplingsEaten)
}

func main() {
	dumplings := 1000
	numPhilosophers := 10

	a := NewLockWithName("chopstick_a")
	b := NewLockWithName("chopstick_b")

	var wg sync.WaitGroup
	for i := range numPhilosophers {
		wg.Add(1)
		p := &Philosopher{fmt.Sprintf("Philosopher #%d", i), a, b, &dumplings}
		go p.Run(&wg)
	}
	wg.Wait()
}
