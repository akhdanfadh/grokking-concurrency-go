package main

import (
	"fmt"
	"time"
)

// This mode demonstrates "threads without multitasking" on a single-core CPU model.
// The first thread (input) blocks forever, starving the others.

type thread struct {
	name string
	task StepTask
}

func (th thread) run(cpu <-chan struct{}) {
	<-cpu // acquire the only cpu core
	fmt.Printf("[%s] acquired CPU (single-core)\n", th.name)

	// in this no-multitasking model, all threads run indefinitely
	// there is no scheduler to preempt it and no time slicing
	for {
		th.task.Step(time.Now(), true) // blocking
		if isGameOver {
			return
		}
	}
}

func runNoMT() {
	// unbuffered channel to model single-core cpu
	cpu := make(chan struct{})

	input := &InputTask{src: NewBlockingStdinSource()}
	world := &WorldTask{}
	render := &RenderTask{}

	// start all tasks, but ensure input gets cpu first (book's intent)
	go thread{name: input.Name(), task: input}.run(cpu)
	time.Sleep(100 * time.Millisecond) // ensure input starts first
	cpu <- struct{}{}
	go thread{name: world.Name(), task: world}.run(cpu)
	go thread{name: render.Name(), task: render}.run(cpu)

	select {} // keep main alive forever
}
