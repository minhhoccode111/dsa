// go run matching_with_mismatches.go

// Find all text positions where the pattern has at most k mismatches.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	mismatchPrime      int64 = 1000000009
	mismatchMultiplier int64 = 1234
)

func buildMismatchHashes(s string, powers []int64) []int64 {
	hashes := make([]int64, len(s)+1)
	for i := 1; i <= len(s); i++ {
		hashes[i] = (mismatchMultiplier*hashes[i-1] + int64(s[i-1])) % mismatchPrime
		powers[i] = powers[i-1] * mismatchMultiplier % mismatchPrime
	}
	return hashes
}

func mismatchSubstringHash(hashes, powers []int64, start, length int) int64 {
	value := (hashes[start+length] - powers[length]*hashes[start]) % mismatchPrime
	if value < 0 {
		value += mismatchPrime
	}
	return value
}

func solveMismatches(k int, text, pattern string) []int {
	result := make([]int, 0)
	m := len(pattern)
	if m == 0 || m > len(text) {
		return result
	}

	powers := make([]int64, len(text)+1)
	powers[0] = 1
	textHashes := buildMismatchHashes(text, powers)
	patternHashes := buildMismatchHashes(pattern, powers)

	for i := 0; i+m <= len(text); i++ {
		a := i
		baseEnd := i + m - 1
		b := baseEnd
		mismatches := 0

		for mismatches <= k {
			mismatch := -1
			for a <= b {
				mid := (a + b) / 2
				length := mid - a + 1

				textHash := mismatchSubstringHash(textHashes, powers, a, length)
				patternHash := mismatchSubstringHash(patternHashes, powers, a-i, length)

				if textHash == patternHash {
					a = mid + 1
				} else {
					mismatch = mid
					b = mid - 1
				}
			}

			if mismatch == -1 {
				result = append(result, i)
				break
			}

			mismatches++
			a = mismatch + 1
			b = baseEnd
		}
	}

	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024), 1<<20)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}

		k, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		positions := solveMismatches(k, fields[1], fields[2])

		if len(positions) == 0 {
			fmt.Println(0)
			continue
		}

		parts := make([]string, 1+len(positions))
		parts[0] = strconv.Itoa(len(positions))
		for i, position := range positions {
			parts[i+1] = strconv.Itoa(position)
		}
		fmt.Println(strings.Join(parts, " "))
	}
}
