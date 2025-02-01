// go run hash_substring.go

// Find all occurrences of a pattern in a text using Rabin-Karp.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	rabinKarpPrime      int64 = 1000000007
	rabinKarpMultiplier int64 = 263
)

func polyHash(s string, prime, multiplier int64) int64 {
	var value int64
	for i := len(s) - 1; i >= 0; i-- {
		value = (value*multiplier + int64(s[i])) % prime
	}
	return value
}

func findOccurrences(pattern, text string) []int {
	n, m := len(text), len(pattern)
	result := make([]int, 0)
	if m == 0 || m > n {
		return result
	}

	patternHash := polyHash(pattern, rabinKarpPrime, rabinKarpMultiplier)

	y := int64(1)
	for i := 0; i < m; i++ {
		y = y * rabinKarpMultiplier % rabinKarpPrime
	}

	hashes := make([]int64, n-m+1)
	hashes[n-m] = polyHash(text[n-m:], rabinKarpPrime, rabinKarpMultiplier)
	for i := n - m - 1; i >= 0; i-- {
		hashes[i] = (rabinKarpMultiplier*hashes[i+1] +
			int64(text[i]) - y*int64(text[i+m])) % rabinKarpPrime
		if hashes[i] < 0 {
			hashes[i] += rabinKarpPrime
		}
	}

	for i := 0; i < len(hashes); i++ {
		if hashes[i] == patternHash && text[i:i+m] == pattern {
			result = append(result, i)
		}
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var pattern, text string
	if _, err := fmt.Fscan(reader, &pattern, &text); err != nil {
		return
	}

	occurrences := findOccurrences(pattern, text)
	parts := make([]string, len(occurrences))
	for i, occurrence := range occurrences {
		parts[i] = fmt.Sprint(occurrence)
	}
	fmt.Println(strings.Join(parts, " "))
}
