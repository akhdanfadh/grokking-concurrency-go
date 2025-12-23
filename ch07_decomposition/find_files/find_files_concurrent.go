//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func searchFile(fileLocation, searchString string) (bool, error) {
	content, err := os.ReadFile(fileLocation)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(content), searchString), nil
}

type result struct {
	fileName string
	found    bool
	err      error
}

func searchFilesConcurrently(fileLocations []string, searchString string) {
	var wg sync.WaitGroup
	results := make(chan result, len(fileLocations))

	// launch goroutines for concurrent file searches
	for _, fileName := range fileLocations {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			found, err := searchFile(fileName, searchString)
			results <- result{fileName, found, err}
		}(fileName)
	}

	// close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect and process results
	for r := range results {
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "Error searching file %s: %v\n", r.fileName, r.err)
			continue
		}
		if r.found {
			fmt.Printf("Found string in file: `%s`\n", r.fileName)
		}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	pattern := filepath.Join(cwd, "books", "*.txt")
	fileLocations, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error globbing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("What word are you trying to find?: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	searchString := scanner.Text()

	startTime := time.Now()
	searchFilesConcurrently(fileLocations, searchString)
	processTime := time.Since(startTime)

	fmt.Printf("PROCESS TIME: %v\n", processTime)
}
