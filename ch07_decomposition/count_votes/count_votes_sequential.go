//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type summary map[int]int

func processVotes(pile []int) summary {
	summary := make(summary)
	for _, vote := range pile {
		summary[vote]++
	}
	return summary
}

func main() {
	numCandidates := 3
	numVoters := 100000

	rand.New(rand.NewSource(time.Now().UnixNano()))
	pile := make([]int, numVoters)
	for i := range pile {
		pile[i] = 1 + rand.Intn(numCandidates)
	}

	start := time.Now()
	counts := processVotes(pile)
	elapsed := time.Since(start)

	fmt.Printf("Total number of votes: %v\n", counts)
	fmt.Printf("Processing time: %s\n", elapsed)
}
