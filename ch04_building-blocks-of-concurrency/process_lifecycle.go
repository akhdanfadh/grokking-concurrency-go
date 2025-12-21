package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func runWorker() {
	fmt.Println("Worker started...")
	time.Sleep(3 * time.Second)
	fmt.Println("Worker is done.")
}

func isAlive(cmd *exec.Cmd) bool {
	// Process is the underlying process, once started
	// ProcessState contains information about an exited process
	return cmd.Process != nil && cmd.ProcessState == nil
}

func main() {
	// worker mode
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		runWorker()
		return
	}

	fmt.Println("Boss requesting Worker's help.")
	cmd := exec.Command(os.Args[0], "worker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("  Worker alive?: %v \n", isAlive(cmd))

	fmt.Println("Boss tells Worker to start.")
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Printf("  Worker alive?: %v \n", isAlive(cmd))

	fmt.Println("Boss goes for coffee.")
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("  Worker alive?: %v \n", isAlive(cmd))

	fmt.Println("Boss patiently waits for Worker to finish and join...")
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Printf("  Worker alive?: %v \n", isAlive(cmd))

	fmt.Println("Boss and Worker are both done!")
}
