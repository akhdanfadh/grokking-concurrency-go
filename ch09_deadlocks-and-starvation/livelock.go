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

	for *p.Dumplings > 0 {
		p.Left.Acquire()
		fmt.Printf("%s grabbed by %s now needs %s\n", p.Left.Name, p.Name, p.Right.Name)

		if !p.Right.TryAcquire() {
			fmt.Printf("%s cannot get %s chopstick, politely concedes...\n", p.Name, p.Right.Name)
			p.Left.Release()
			continue
		}
		fmt.Printf("%s chopstick grabbed by %s\n", p.Right.Name, p.Name)

		if *p.Dumplings <= 0 {
			p.Right.Release()
			p.Left.Release()
			return
		}

		*p.Dumplings--
		fmt.Printf("%s eats a dumpling. Dumplings left: %d\n", p.Name, *p.Dumplings)
		fmt.Printf("%s is thinking...\n", p.Name)
		time.Sleep(1 * time.Second)

		p.Right.Release()
		p.Left.Release()
	}
}

func main() {
	dumplings := 20

	a := NewLockWithName("chopstick_a")
	b := NewLockWithName("chopstick_b")
	p1 := &Philosopher{"Philosopher #1", a, b, &dumplings}
	p2 := &Philosopher{"Philosopher #2", b, a, &dumplings}

	var wg sync.WaitGroup
	wg.Add(2)
	go p1.Run(&wg)
	go p2.Run(&wg)
	wg.Wait()
}
