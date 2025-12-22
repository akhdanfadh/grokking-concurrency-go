//go:build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

// task represents a unit of work to be executed by a worker
type task func(workerName string)

// worker represents a single thread of execution in the thread pool
type worker struct {
	name  string          // for visibility in output
	tasks <-chan task     // receive-only channel from the pool's queue
	wg    *sync.WaitGroup // shared with the pool to track task completion
}

func NewWorker(name string, tasks <-chan task, wg *sync.WaitGroup) *worker {
	return &worker{name: name, tasks: tasks, wg: wg}
}

// start begins the worker's task processing loop in a new goroutine/thread
func (w *worker) start() {
	go func() {
		for task := range w.tasks {
			task(w.name)
		}
	}()
}

// threadPool manages a pool of worker threads to execute submitted tasks
type threadPool struct {
	tasks chan task      // message queue
	wg    sync.WaitGroup // completion tracker
	once  sync.Once      // ensure close() is safe if called multiple times
}

func newThreadPool(numWorkers, queueSize int) *threadPool {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if queueSize <= 0 {
		queueSize = numWorkers
	}

	// creates and starts several worker threads
	tp := &threadPool{tasks: make(chan task, queueSize)}
	for i := range numWorkers {
		NewWorker(fmt.Sprintf("Thread-%d", i+1), tp.tasks, &tp.wg).start()
	}
	return tp
}

// submit enqueues a task for execution
func (tp *threadPool) submit(t task) {
	tp.wg.Add(1)

	// wrap the task so Done is always called even if it panics
	tp.tasks <- func(workerName string) {
		defer tp.wg.Done()

		// protect the worker from "task failure" (panic), akin to catching exceptions
		// without recover, a panic would terminate the worker goroutine, shrinking the pool
		// with recover, the worker can continue processing further tasks
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%s recovered task panic: %v\n", workerName, r)
			}
		}()

		t(workerName)
	}
}

// waitCompletion blocks until all submitted tasks have completed
func (tp *threadPool) waitCompletion() { tp.wg.Wait() }

// close gracefully shuts down the thread pool by closing the task channel,
// thus signaling workers to exit once all tasks are done and no new tasks will arrive
func (tp *threadPool) close() { tp.once.Do(func() { close(tp.tasks) }) }

// cpuWaster simulates a CPU-bound task by sleeping for a fixed duration
func cpuWaster(i int) task {
	return func(workerName string) {
		fmt.Printf("%s doing %d work\n", workerName, i)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	// creates a thread pool with 5 workers and a queue size of 5
	pool := newThreadPool(5, 5)
	for i := range 20 { // add 20 tasks to the pool
		pool.submit(cpuWaster(i))
	}

	fmt.Println("All work requests sent")
	pool.waitCompletion()
	fmt.Println("All work complete")
	pool.close()
}
