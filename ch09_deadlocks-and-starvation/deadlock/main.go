package main

import (
	"fmt"
	"os"
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

		p.Right.Acquire()
		fmt.Printf("%s grabbed by %s\n", p.Right.Name, p.Name)

		*p.Dumplings--
		fmt.Printf("%s eats a dumpling. Dumplings left: %d\n", p.Name, *p.Dumplings)

		p.Right.Release()
		fmt.Printf("%s released by %s\n", p.Right.Name, p.Name)

		p.Left.Release()
		fmt.Printf("%s released by %s\n", p.Left.Name, p.Name)

		fmt.Printf("%s is thinking...\n", p.Name)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [deadlock|arbitrator|hierarchy]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deadlock":
		runDeadlock()
	case "arbitrator":
		runArbitrator()
	case "hierarchy":
		runHierarchy()
	default:
		fmt.Println("unknown mode:", os.Args[1])
		os.Exit(1)
	}
}
