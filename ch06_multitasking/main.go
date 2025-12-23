package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [nomt|mt]")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "nomt":
		runNoMT()
	case "mt":
		runMT()
	default:
		fmt.Printf("Unknown mode %q. Use 'nomt' or 'mt'.\n", mode)
		os.Exit(1)
	}
}
