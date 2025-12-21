//go:build ignore

package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
	"time"
)

// getCombinations generates all possible password combinations
func getCombinations(length, minNumber int, maxNumber *int) []string {
	var combinations []string

	// calculate maximum number based on the length if not provided
	max := 0
	if maxNumber == nil {
		max = int(math.Pow(10, float64(length)) - 1)
	} else {
		max = *maxNumber
	}

	// go through all possible combinations in a given range
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
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash) // as lowercase hexadecimal
}

// checkPassword compares the resulted cryptographic hash with the expected one
func checkPassword(expectedCryptoHash, possiblePassword string) bool {
	actualCryptoHash := getCryptoHash(possiblePassword)
	return expectedCryptoHash == actualCryptoHash
}

// crackPassword tries to find the password by checking all possible combinations (brute force)
func crackPassword(cryptoHash string, length int) {
	fmt.Println("Processing number combinations sequentially")
	startTime := time.Now()

	combinations := getCombinations(length, 0, nil)
	for _, combination := range combinations {
		if checkPassword(cryptoHash, combination) {
			fmt.Printf("PASSWORD CRACKED: %s\n", combination)
			break
		}
	}

	processTime := time.Since(startTime)
	fmt.Printf("PROCESS TIME: %s\n", processTime)
}

func main() {
	cryptoHash := "e24df920078c3dd4e7e8d2442f00e5c9ab2a231bb3918d65cc50906e49ecaef4"
	length := 8
	crackPassword(cryptoHash, length)
}
