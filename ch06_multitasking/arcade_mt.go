package main

import (
	"bufio"
	"os"
	"strings"
	"time"
)

// timeSlice defines how often the scheduler switches between tasks, not the game logic cadence.
const timeSlice = 500 * time.Millisecond

// readInputBlockingInto reads stdin lines into the provided string pointer.
//
// This is a dedicated blocking goroutine (I/O worker) so that the main scheduler loop never blocks on input.
func readInputBlockingInto(pending *string) {
	r := bufio.NewReader(os.Stdin)
	for {
		if isGameOver {
			return
		}
		line, err := r.ReadString('\n')
		if err != nil {
			isGameOver = true
			gameOverMsg = "input error"
			return
		}
		*pending = strings.TrimSpace(line)
	}
}

func runMT() {
	var pending string

	// blocking input runs "in the background" so the scheduler loop never blocks on stdin
	go readInputBlockingInto(&pending)

	input := &InputTask{src: NewSharedBufferSource(&pending)}
	world := &WorldTask{}
	render := &RenderTask{}
	tasks := []StepTask{input, world, render}

	// track per-task next eligible run time based on task.Period()
	nextRun := make(map[string]time.Time, len(tasks))
	now := time.Now()
	for _, t := range tasks {
		nextRun[t.Name()] = now
	}

	ticker := time.NewTicker(timeSlice)
	defer ticker.Stop()
	for now = range ticker.C {
		// "scheduler loop": interleaves tasks, granting each a slice when eligible
		for _, t := range tasks {
			if isGameOver {
				// one final render, then wait for Enter
				(&RenderTask{}).Step(now, false)
				_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
				return
			}

			nr := nextRun[t.Name()]
			if now.Before(nr) {
				continue // not yet eligible
			}

			// in multitasking mode, Step must never block
			t.Step(now, false)

			// update nextRun according to the task cadence
			if p := t.Period(); p > 0 {
				nextRun[t.Name()] = now.Add(p)
			} else {
				// period=0: eligible every scheduler tick
				nextRun[t.Name()] = now
			}
		}
	}
}
