// go run different_summands.go

// Represent a positive integer as the sum of the maximum number of pairwise
// distinct positive integers.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func optimalSummands(n int) []int {
	summands := make([]int, 0)
	current := 1

	for n > 0 {
		if n-current > current {
			summands = append(summands, current)
			n -= current
			current++
		} else {
			summands = append(summands, n)
			break
		}
	}

	return summands
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	summands := optimalSummands(n)

	fmt.Println(len(summands))
	for i, summand := range summands {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(summand)
	}
	fmt.Println()
}
