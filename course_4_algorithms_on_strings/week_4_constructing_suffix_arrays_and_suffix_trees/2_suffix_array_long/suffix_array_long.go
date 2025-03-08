// go run suffix_array_long.go

// Construct the suffix array of a string in O(n log n).

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func buildSuffixArray(text string) []int {
	return strs.SuffixArray([]byte(text))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	for i, suffix := range buildSuffixArray(text) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(suffix)
	}
	fmt.Println()
}
