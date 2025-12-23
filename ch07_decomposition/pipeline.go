//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

type Washload string

// Washer is a thread representing a washing machine.
type Washer struct {
	In  <-chan Washload
	Out chan<- Washload
}

func (w *Washer) Run() {
	for load := range w.In {
		fmt.Printf("Washer: washing %s...\n", load)
		time.Sleep(4 * time.Second)
		w.Out <- load
	}
	close(w.Out) // no more inputs so close downstream to signal completion
}

// Dryer is a thread representing a dryer.
type Dryer struct {
	In  <-chan Washload
	Out chan<- Washload
}

func (d *Dryer) Run() {
	for load := range d.In {
		fmt.Printf("Dryer: drying %s...\n", load)
		time.Sleep(2 * time.Second)
		d.Out <- load
	}
	close(d.Out)
}

// Folder is a thread representing the folding action.
type Folder struct {
	In <-chan Washload
}

func (f *Folder) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for load := range f.In {
		fmt.Printf("Folder: folding %s...\n", load)
		time.Sleep(1 * time.Second)
		fmt.Printf("Folder: %s done!\n", load)
	}
}

// Pipeline represents a washer, dryer, and folder, linked by queues (channels).
type Pipeline struct{}

func (*Pipeline) assembleLaundryForWashing() <-chan Washload {
	washloadCount := 4

	toBeWashed := make(chan Washload, washloadCount)
	for i := range washloadCount {
		toBeWashed <- Washload(fmt.Sprintf("Washload #%d", i))
	}
	close(toBeWashed) // input stage is fully populated
	return toBeWashed
}

func (p *Pipeline) RunConcurrently() {
	// set up the queues in the pipeline
	toBeWashed := p.assembleLaundryForWashing()
	toBeDried := make(chan Washload)
	toBeFolded := make(chan Washload)

	// start the threas linked by the queues
	go (&Washer{In: toBeWashed, Out: toBeDried}).Run()
	go (&Dryer{In: toBeDried, Out: toBeFolded}).Run()

	// wait for folder to finish (as it is the last stage)
	var wg sync.WaitGroup
	wg.Add(1)
	go (&Folder{In: toBeFolded}).Run(&wg)
	wg.Wait()

	fmt.Println("All done!")
}

func main() {
	pipeline := &Pipeline{}
	pipeline.RunConcurrently()
}
