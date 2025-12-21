package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func runWriter(w *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	defer w.Close() // signal EOF to reader after data is consumed

	fmt.Println("Writer: Sending rubber duck...")

	// send a single line as the message frame
	if _, err := io.WriteString(w, "Rubber duck\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Writer: write error: %v\n", err)
		return
	}
}

func runReader(r *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	defer r.Close()

	fmt.Println("Reader: Reading...")

	// ReadString blocks until it sees '\n' or encounters an error/EOF
	br := bufio.NewReader(r)
	msg, err := br.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reader: read error: %v\n", err)
		return
	}

	fmt.Printf("Reader: Received: %q\n", strings.TrimSpace(msg))
}

func main() {
	// os.Pipe returns a connected pair of Files; reads from r return bytes written to w
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Main: pipe error: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go runWriter(w, &wg)
	go runReader(r, &wg)
	wg.Wait() // block main until child threads finish, python is implicit on this
}
