// go run kmp.go

// Find all occurrences of a pattern in a text using the KMP algorithm.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func computePrefixFunction(pattern string) []int {
	n := len(pattern)
	if n == 0 {
		return nil
	}
	prefix := make([]int, n)
	border := 0
	for i := 1; i < n; i++ {
		for border > 0 && pattern[i] != pattern[border] {
			border = prefix[border-1]
		}
		if pattern[i] == pattern[border] {
			border++
		} else {
			border = 0
		}
		prefix[i] = border
	}
	return prefix
}

func findPattern(pattern, text string) []int {
	combined := pattern + "$" + text
	prefix := computePrefixFunction(combined)
	var result []int
	for i := len(pattern) + 1; i < len(combined); i++ {
		if prefix[i] == len(pattern) {
			result = append(result, i-2*len(pattern))
		}
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var pattern, text string
	if _, err := fmt.Fscan(reader, &pattern); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	for i, position := range findPattern(pattern, text) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(position)
	}
	fmt.Println()
}
