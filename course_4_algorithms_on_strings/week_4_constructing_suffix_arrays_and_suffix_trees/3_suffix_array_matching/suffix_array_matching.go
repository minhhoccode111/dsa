// go run suffix_array_matching.go

// Find all occurrences of a collection of patterns in a text using a suffix
// array and binary search.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func compareSuffix(text string, position int, pattern string) int {
	limit := len(pattern)
	if remaining := len(text) - position; remaining < limit {
		limit = remaining
	}
	for i := 0; i < limit; i++ {
		if text[position+i] < pattern[i] {
			return -1
		}
		if text[position+i] > pattern[i] {
			return 1
		}
	}
	if len(text)-position < len(pattern) {
		return -1
	}
	return 0
}

func findOccurrences(text string, patterns []string) []int {
	if len(text) == 0 || len(patterns) == 0 {
		return nil
	}

	sa := strs.SuffixArray([]byte(text))
	seen := make([]bool, len(text))
	var result []int

	for _, pattern := range patterns {
		if len(pattern) == 0 {
			continue
		}
		low := sort.Search(len(sa), func(i int) bool {
			return compareSuffix(text, sa[i], pattern) >= 0
		})
		high := sort.Search(len(sa), func(i int) bool {
			return compareSuffix(text, sa[i], pattern) > 0
		})
		for i := low; i < high; i++ {
			position := sa[i]
			if !seen[position] {
				seen[position] = true
				result = append(result, position)
			}
		}
	}

	sort.Ints(result)
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	var patternCount int
	if _, err := fmt.Fscan(reader, &patternCount); err != nil {
		return
	}
	patterns := make([]string, patternCount)
	for i := range patterns {
		if _, err := fmt.Fscan(reader, &patterns[i]); err != nil {
			return
		}
	}

	result := findOccurrences(text, patterns)
	for i, position := range result {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(position)
	}
	fmt.Println()
}
