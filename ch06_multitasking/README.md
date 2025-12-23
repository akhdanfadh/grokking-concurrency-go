# Chapter 6 - Go Implementation Notes

## Overview

**Problem: three jobs, one worker**

The game has three essential functions that need to run continuously.
First, we need to read input from the joystick so Pac-Man can move.
Second, we need to update the game world—ghosts have to move around, dots need to be eaten, collisions need to be checked.
Third, we need to draw everything on screen so the player can see what's happening.

If we just run these three functions one after another in a loop, we'd have a problem.
The input function would block waiting for the player to press a button, and while it's waiting, nothing else happens.
The ghosts freeze, the screen doesn't update.

**Solution: time slicing**

The book introduces preemptive multitasking as the solution.
The operating system does something like this: it gives each task a tiny slice of time (say, 500 milliseconds), then forcibly pauses it and moves to the next task.

The Python implementation in the book demonstrates this with an `InterruptService` that ticks every half second, setting an event flag that lets tasks take turns.
Each task waits for the flag, grabs it exclusively, does its work, and eventually gives it back.

## Implementation

For this chapter, you can run `go run . nomt` or `go run . mt` to see two different versions of the same game.
In the multitasking version, you'll need to press Enter quickly after typing commands (w/a/s/d to move, q to quit) because stdin is line-buffered, but the game world keeps updating between your inputs.

**The broken version (`nomt`)**

In `arcade_nomt.go`, three goroutines are created for input, game logic, and rendering.
To model a single "CPU", we use an unbuffered channel where only one goroutine can receive from it at a time.
The first goroutine to grab the CPU runs forever in a blocking loop, and the others never get a chance to run.

In our case, the input goroutine wins the race and sits waiting for stdin.
Since each task is running in their own infinite loop, once one gets the CPU, it never gives it up.
The game is effectively frozen.
This demonstrates the core problem the book is explaining: without a scheduler to preempt tasks, a single blocking operation can starve everything else.

**The working version (`mt`)**

The multitasking version in `arcade_mt.go` fixes this by separating concerns.
A dedicated goroutine sits in the background, continuously reading from stdin into a shared string buffer.
This goroutine can block all it wants because it's not part of the main scheduler loop.

Meanwhile, the scheduler runs on a ticker; every 500 milliseconds, it wakes up and gives each task a chance to run.
But crucially, when it calls each task's `Step` method, it passes `block=false`, which means "don't wait for anything, just do what you can right now and return."
This ticker-based approach naturally simulates time slicing.

The input task checks the shared buffer.
If there's a command waiting, great, process it.
If not, no problem, just return immediately and we'll check again next time.
This is fundamentally different from the Python version's approach, where Python uses threading events to coordinate access, but the principle is the same: never block the scheduler.

### the abstraction layer: `InputSource` interface

One thing I did differently from the Python code is create an `InputSource` interface with two implementations.
The `BlockingStdinSource` is what the non-multitasking version uses:
it just calls `ReadString('\n')` and blocks until input arrives.
The `SharedBufferSource` is what the multitasking version uses:
it checks a string pointer, and if there's something there, it consumes it and returns.
If not, it just returns immediately with no data.

This abstraction made it really clear to me what the difference is between blocking and non-blocking I/O.
The task code doesn't change much, but the source it reads from has completely different behavior.

### task scheduling with periods

Another aspect I added was task periods.
The input task has a period of zero, meaning it can run every scheduler tick, i.e, we want to check for input as frequently as possible.
But the world computation and rendering tasks have a period of 1 second (the `Delay` constant), meaning they only run once per second even if the scheduler ticks more often.

This is tracked with a `nextRun` map that stores when each task is eligible to run next.
If the current time is before a task's next eligible time, we skip it.
This gives us throttling for free: we can have a fast scheduler tick but slower game logic updates.

The Python code doesn't really implement per-task periods in the same way.
It just sleeps for `DELAY` inside each task function, which means the tasks themselves control their pacing.
My approach moves that control to the scheduler, which feels more like how real OS schedulers work.
