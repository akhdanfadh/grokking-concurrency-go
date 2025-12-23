//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type summary map[int]int

// processVotes counts the number of votes each candidate received in parallel.
func processVotes(pile []int, workerCount int) summary {
	voteCount := len(pile)
	votePerWorker := voteCount / workerCount

	// fork: divide the votes among workers
	votePiles := make([][]int, 0, workerCount)
	for i := range workerCount {
		start := i * votePerWorker
		end := (i + 1) * votePerWorker
		if i == workerCount-1 {
			end = voteCount
		}
		votePiles = append(votePiles, pile[start:end])
	}

	// create a fixed-size worker set: one goroutine per chunk
	workerSummaries := make([]summary, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := range workerCount {
		go func(i int) {
			defer wg.Done()
			workerSummaries[i] = processPile(votePiles[i])
		}(i)
	}
	wg.Wait()

	// join: merge worker summaries
	totalSummary := make(summary)
	for _, workerSummary := range workerSummaries {
		fmt.Printf("Votes from staff member: %v\n", workerSummary)
		for candidate, count := range workerSummary {
			totalSummary[candidate] += count
		}
	}
	return totalSummary
}

func processPile(pile []int) summary {
	summary := make(summary)
	for _, vote := range pile {
		summary[vote]++
	}
	return summary
}

func main() {
	workerCount := 4
	numCandidates := 3
	numVoters := 100000

	rand.New(rand.NewSource(time.Now().UnixNano()))
	pile := make([]int, numVoters)
	for i := range pile {
		pile[i] = 1 + rand.Intn(numCandidates)
	}

	start := time.Now()
	counts := processVotes(pile, workerCount)
	elapsed := time.Since(start)

	fmt.Printf("Total number of votes: %v\n", counts)
	fmt.Printf("Processing time: %s\n", elapsed)
}
