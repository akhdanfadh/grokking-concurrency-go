//go:build ignore

package main

import (
	"fmt"
	"net" // go exposes sockets via the net package
	"os"
	"sync"
	"time"
)

const (
	// Remember: A socket file is not a socket. A socket is a kernel object.
	// The file here is to facilitate the connection between threads,
	// just like TCP port number is not the TCP connection itself.
	sockFile = "./mailbox"
	// Buffer size for receiving data from the socket connection.
	bufferSize = 1024
)

func runSender(wg *sync.WaitGroup) {
	defer wg.Done()

	// Creates a client socket for this thread, and
	// connects the socket to the "channel" (the mailbox file).
	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		fmt.Printf("Sender: dial error: %v\n", err)
		return
	}
	defer conn.Close()

	// Sends a series of messages over the client socket
	messages := []string{"Hello", " ", "world!"}
	for _, msg := range messages {
		fmt.Printf("Sender: Send: %q\n", msg)
		if _, err := conn.Write([]byte(msg)); err != nil {
			fmt.Printf("Sender: write error: %v\n", err)
			return
		}
	}
}

func runReceiver(wg *sync.WaitGroup) {
	defer wg.Done()

	// Creates a listening socket for this thread,
	// and binds it to a "channel" (the mailbox file),
	// and starts listening for incoming connections.
	ln, err := net.Listen("unix", sockFile)
	if err != nil {
		fmt.Printf("Receiver: listen error: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Receiver: Listening for incoming messages...")

	// Accepts a connection on the listening socket, and returns
	// a new connection socket for the actual communication endpoint
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Receiver: accept error: %v\n", err)
		return
	}

	// Reads messages from the connection socket
	buf := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			msg := string(buf[:n])
			fmt.Printf("Receiver: Received: %q\n", msg)
		}
		if err != nil {
			break
		}
	}
}

func main() {
	_ = os.Remove(sockFile) // ensure clean state

	var wg sync.WaitGroup
	wg.Add(2)
	go runReceiver(&wg)
	time.Sleep(1 * time.Second) // ensure receiver is ready
	go runSender(&wg)
	wg.Wait() // block until both sender and receiver complete

	_ = os.Remove(sockFile) // clean up socket file on exit
}
