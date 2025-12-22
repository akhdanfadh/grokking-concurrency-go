//go:build ignore

package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type chunkRange struct {
	start, end int
}

// getCombinations generates all possible password combinations
func getCombinations(length, minNumber int, maxNumber *int) []string {
	// calculate maximum number based on the length if not provided
	max := 0
	if maxNumber == nil {
		max = int(math.Pow(10, float64(length)) - 1)
	} else {
		max = *maxNumber
	}

	// go through all possible combinations in a given range
	var combinations []string
	for i := minNumber; i <= max; i++ {
		strNum := strconv.Itoa(i)
		// fill in the missing numbers with zeros
		zeros := ""
		for range length - len(strNum) {
			zeros += "0"
		}
		combinations = append(combinations, zeros+strNum)
	}
	return combinations
}

// getCryptoHash calculates the cryptographic hash of the password
func getCryptoHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password))) // as lowercase hexadecimal
}

// checkPassword compares the resulted cryptographic hash with the expected one
func checkPassword(expectedCryptoHash, possiblePassword string) bool {
	return expectedCryptoHash == getCryptoHash(possiblePassword)
}

// getChunks split the passwords into chunks using break points
func getChunks(numRanges, length int) []chunkRange {
	maxNumber := int(math.Pow(10, float64(length)) - 1)

	chunkStarts := make([]int, 0, numRanges)
	for i := range numRanges {
		chunkStarts = append(chunkStarts, i*(maxNumber/numRanges))
	}

	chunkEnds := make([]int, 0, numRanges)
	for i := 1; i < len(chunkStarts); i++ {
		chunkEnds = append(chunkEnds, chunkStarts[i]-1)
	}
	chunkEnds = append(chunkEnds, maxNumber)

	chunks := make([]chunkRange, 0, numRanges)
	for i := range numRanges {
		chunks = append(chunks, chunkRange{start: chunkStarts[i], end: chunkEnds[i]})
	}
	return chunks
}

// crackChunk tries to find the password in a given chunk by brute force
func crackChunk(cryptoHash string, length, chunkStart, chunkEnd int) {
	fmt.Printf("Processing %d to %d\n", chunkStart, chunkEnd)
	combinations := getCombinations(length, chunkStart, &chunkEnd)
	for _, combination := range combinations {
		if checkPassword(cryptoHash, combination) {
			fmt.Fprintln(os.Stderr, combination) // log to stderr for master process (workaround)
		}
	}
}

type worker struct {
	chunk chunkRange
	cmd   *exec.Cmd
	buf   *bytes.Buffer
}

// crackPasswordParallel orchestrate cracking the password between different processes
func crackPasswordParallel(cryptoHash string, length int) {
	numCores := runtime.NumCPU()
	fmt.Println("Processing number combinations concurrently")
	startTime := time.Now()

	chunks := getChunks(numCores, length)
	workers := make([]worker, 0, len(chunks))

	// start worker processes
	for _, chunk := range chunks {
		cmd := exec.Command(
			os.Args[0],
			"-role=worker",
			"-hash="+cryptoHash,
			"-length="+strconv.Itoa(length),
			"-start="+strconv.Itoa(chunk.start),
			"-end="+strconv.Itoa(chunk.end),
		)

		var errBuf bytes.Buffer
		cmd.Stdout = os.Stdout // live progress output
		cmd.Stderr = &errBuf   // capture "return value"

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Master: failed to start worker for %d..%d: %v\n", chunk.start, chunk.end, err)
			continue
		}

		workers = append(workers, worker{chunk: chunk, cmd: cmd, buf: &errBuf})
	}

	// wait for all workers to finish
	for _, w := range workers {
		if err := w.cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "Master: worker for %d..%d exited with error: %v\n", w.chunk.start, w.chunk.end, err)
		}
	}

	// collect the first found results
	var results []string
	for _, w := range workers {
		for line := range strings.SplitSeq(w.buf.String(), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				results = append(results, line)
			}
		}
	}

	fmt.Printf("PASSWORD CRACKED: %s\n", results[0])
	processTime := time.Since(startTime)
	fmt.Printf("PROCESS TIME: %s\n", processTime)
}

func main() {
	role := flag.String("role", "master", "Role of the process: master or worker")
	hash := flag.String("hash", "", "Cryptographic hash to crack (for worker role)")
	length := flag.Int("length", 0, "Length of the password to crack (for worker role)")
	start := flag.Int("start", 0, "Start of the chunk range (for worker role)")
	end := flag.Int("end", 0, "End of the chunk range (for worker role)")
	flag.Parse()

	switch *role {
	case "master":
		cryptoHash := "e24df920078c3dd4e7e8d2442f00e5c9ab2a231bb3918d65cc50906e49ecaef4"
		length := 8
		crackPasswordParallel(cryptoHash, length)
	case "worker":
		crackChunk(*hash, *length, *start, *end)
	default:
		fmt.Fprintf(os.Stderr, "Unknown role: %s\n", *role)
		os.Exit(2)
	}
}
