// go run bwtmatching.go

// Count occurrences of patterns using the Burrows-Wheeler transform.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func preprocessBWT(encoded string) ([256]int, [256][]int) {
	var counts [256]int
	var starts [256]int
	var occurrences [256][]int
	var present [256]bool

	for i := 0; i < len(encoded); i++ {
		counts[encoded[i]]++
		present[encoded[i]] = true
	}

	total := 0
	for char := 0; char < 256; char++ {
		starts[char] = total
		total += counts[char]
		if present[char] {
			occurrences[char] = make([]int, len(encoded)+1)
		}
	}

	for i := 0; i < len(encoded); i++ {
		char := encoded[i]
		for candidate := 0; candidate < 256; candidate++ {
			if present[candidate] {
				occurrences[candidate][i+1] = occurrences[candidate][i]
			}
		}
		occurrences[char][i+1]++
	}

	return starts, occurrences
}

func countOccurrences(pattern, encoded string, starts [256]int, occurrences [256][]int) int {
	if len(encoded) == 0 || len(pattern) == 0 {
		return 0
	}

	top := 0
	bottom := len(encoded) - 1
	for i := len(pattern) - 1; i >= 0; i-- {
		symbol := pattern[i]
		if occurrences[symbol] == nil {
			return 0
		}
		beforeTop := occurrences[symbol][top]
		beforeBottom := occurrences[symbol][bottom+1]
		if beforeBottom <= beforeTop {
			return 0
		}
		top = starts[symbol] + beforeTop
		bottom = starts[symbol] + beforeBottom - 1
	}
	return bottom - top + 1
}

func bwt(text string) string {
	n := len(text)
	if n == 0 {
		return ""
	}
	data := []byte(text)
	sa := strs.SuffixArray(data)
	result := make([]byte, n)
	for i, suffix := range sa {
		previous := suffix - 1
		if previous < 0 {
			previous = n - 1
		}
		result[i] = data[previous]
	}
	return string(result)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var encoded string
	if _, err := fmt.Fscan(reader, &encoded); err != nil {
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

	starts, occurrences := preprocessBWT(encoded)
	for i, pattern := range patterns {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(countOccurrences(pattern, encoded, starts, occurrences))
	}
	fmt.Println()
}
