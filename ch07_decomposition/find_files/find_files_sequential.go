//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func searchFile(fileLocation, searchString string) (bool, error) {
	content, err := os.ReadFile(fileLocation)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(content), searchString), nil
}

func searchFilesSequentially(fileLocations []string, searchString string) {
	for _, fileName := range fileLocations {
		result, err := searchFile(fileName, searchString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching file %s: %v\n", fileName, err)
			continue
		}
		if result {
			fmt.Printf("Found word in file: `%s`\n", fileName)
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
	searchFilesSequentially(fileLocations, searchString)
	processTime := time.Since(startTime)

	fmt.Printf("PROCESS TIME: %v\n", processTime)
}
