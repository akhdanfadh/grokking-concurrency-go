package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

var (
	role = flag.String("role", "parent", "one of: parent, reader, writer")
	path = flag.String("path", "rubberduck.fifo", "path to FIFO (named pipe)")
)

func printAndExit(msg string, code int) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func main() {
	flag.Parse()

	if runtime.GOOS == "windows" {
		printAndExit("named pipes are not supported on Windows in this example", 1)
	}

	switch *role {
	case "parent":
		runParent()
	case "reader":
		runReader()
	case "writer":
		runWriter()
	default:
		printAndExit("unknown role: "+*role, 2)
	}
}

func runParent() {
	_ = os.Remove(*path) // ensure clean state

	// create named pipe (FIFO) with rw-rw-rw- masked by umask
	if err := syscall.Mkfifo(*path, 0666); err != nil {
		printAndExit("Parent: mkfifo error: "+err.Error(), 1)
	}
	defer os.Remove(*path) // clean up FIFO on exit

	fmt.Printf("Parent: Created FIFO at %s\n", *path)

	// spawn separate tasks that will use the FIFO
	readerCmd := exec.Command(os.Args[0], "-role=reader")
	writerCmd := exec.Command(os.Args[0], "-role=writer")

	readerCmd.Stdout = os.Stdout
	readerCmd.Stderr = os.Stderr
	writerCmd.Stdout = os.Stdout
	writerCmd.Stderr = os.Stderr

	if err := readerCmd.Start(); err != nil {
		printAndExit("Parent: starting reader error: "+err.Error(), 1)
	}
	if err := writerCmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Parent: starting writer error: %v\n", err)
		_ = readerCmd.Process.Kill()
		os.Exit(1)
	}

	// wait for both to finish
	if err := readerCmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Parent: reader error: %v\n", err)
	}
	if err := writerCmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Parent: writer error: %v\n", err)
	}

	fmt.Println("Parent: Done.")
}

func runReader() {
	// open FIFO for reading, this may block until a writer opens the FIFO
	f, err := os.OpenFile(*path, os.O_RDONLY, 0)
	if err != nil {
		printAndExit("Reader: open error: "+err.Error(), 1)
	}
	defer f.Close()

	fmt.Println("Reader: Reading...")

	// ReadString blocks until it sees '\n' or encounters an error/EOF
	br := bufio.NewReader(f)
	line, err := br.ReadString('\n')
	if err != nil {
		printAndExit("Reader: read error: "+err.Error(), 1)
	}

	fmt.Printf("Reader: Received: %q\n", strings.TrimSpace(line))
}

func runWriter() {
	// open FIFO for writing, this may block until a reader opens the FIFO
	f, err := os.OpenFile(*path, os.O_WRONLY, 0)
	if err != nil {
		printAndExit("Writer: open error: "+err.Error(), 1)
	}
	defer f.Close()

	fmt.Println("Writer: Sending rubber duck...")

	// FIFO is a byte stream, newline is simple message framing
	if _, err := f.WriteString("Rubber duck\n"); err != nil {
		printAndExit("Writer: write error: "+err.Error(), 1)
	}
}
