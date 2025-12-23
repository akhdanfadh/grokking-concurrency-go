//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

const SIZE = 5

var (
	BUFFER      = make([]string, SIZE) // shared memory
	producerIdx = 0
	mutex       = sync.Mutex{} // protects BUFFER access

	// track empty and full slots in the BUFFER
	empty = NewSemaphore(SIZE, SIZE)
	full  = NewSemaphore(0, SIZE)
)

type Semaphore chan struct{}

func (s Semaphore) Acquire() { <-s }
func (s Semaphore) Release() { s <- struct{}{} }
func NewSemaphore(initial, max int) Semaphore {
	s := make(Semaphore, max)
	for range initial {
		s.Release()
	}
	return s
}

// Producer thread will produce an item and put it into the buffer
type Producer struct {
	name         string
	counter      int
	maximumItems int
}

func NewProducer(name string, maximumItems int) *Producer {
	return &Producer{name: name, counter: 0, maximumItems: maximumItems}
}

// nextIndex get the next empty buffer index
func (p *Producer) nextIndex(index int) int {
	return (index + 1) % SIZE
}

func (p *Producer) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for p.counter < p.maximumItems {
		empty.Acquire() // wait until the buffer have some empty slots

		// critical section for changing the buffer
		mutex.Lock()
		p.counter++
		BUFFER[producerIdx] = fmt.Sprintf("%s-%d", p.name, p.counter)
		fmt.Printf("%s produced: '%s' into slot %d\n", p.name, BUFFER[producerIdx], producerIdx)
		producerIdx = p.nextIndex(producerIdx)
		mutex.Unlock()

		full.Release()          // buffer have one more item to consume
		time.Sleep(time.Second) // simulate some real action here
	}
}

type Consumer struct {
	name         string
	idx, counter int
	maximumItems int
}

func NewConsumer(name string, maximumItems int) *Consumer {
	return &Consumer{name: name, idx: 0, counter: 0, maximumItems: maximumItems}
}

// nextIndex get the next buffer index to consume
func (c *Consumer) nextIndex() int {
	return (c.idx + 1) % SIZE
}

func (c *Consumer) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for c.counter < c.maximumItems {
		full.Acquire() // wait until the buffer have some new items to consume

		// critical section for changing the buffer
		mutex.Lock()
		item := BUFFER[c.idx]
		BUFFER[c.idx] = ""
		fmt.Printf("%s consumed item: '%s' from slot %d\n", c.name, item, c.idx)
		c.idx = c.nextIndex()
		c.counter++
		mutex.Unlock()

		empty.Release()             // one more empty slot is available in buffer
		time.Sleep(2 * time.Second) // simulate some real action here
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go NewProducer("SpongeBob", 5).run(&wg)
	go NewProducer("Patrick", 5).run(&wg)
	go NewConsumer("Squidward", 10).run(&wg)
	wg.Wait()
}
