package main

import (
	"fmt"
	"os"
	"os/exec"
)

// runChild run some logic inside child process
func runChild() {
	fmt.Println("Child: I am the child process")
	fmt.Printf("Child: Child's PID: %d\n", os.Getpid())
	fmt.Printf("Child: Parent's PID: %d\n", os.Getppid())
}

// startParent starts multiple child processes with workaround
//
// Go doesn't have a built-in way to fork current program and run a specific function in a child process.
// As a workaround, we will run this same program again, but tell it to behave differently with a flag.
func startParent(numChildren int) {
	fmt.Println("Parent: I am the parent process")
	fmt.Printf("Parent: Parent's PID: %d\n", os.Getpid())

	cmds := make([]*exec.Cmd, 0, numChildren)
	for i := range numChildren {
		fmt.Printf("Starting Process %d\n", i)
		cmd := exec.Command(os.Args[0], "-child")

		// prints child's output to the same console as parent
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// start the child process in the background
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting child %d: %v\n", i, err)
			continue
		}
		cmds = append(cmds, cmd)
	}

	// wait for all children to complete
	// otherwise the parent may exits immediately and some children will be "re-parented" to PID 1 by the system
	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Child %d exited with error: %v\n", i, err)
		}
	}
}

func main() {
	// check if we are running as a child
	if len(os.Args) > 1 && os.Args[1] == "-child" {
		runChild()
		return
	}

	// start as parent
	numChildren := 3
	startParent(numChildren)
}
