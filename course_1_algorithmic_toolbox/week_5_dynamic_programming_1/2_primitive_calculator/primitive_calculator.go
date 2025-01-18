// go run primitive_calculator.go

// Find the minimum number of operations needed to get n from 1 using only
// add 1, multiply by 2, and multiply by 3.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func computeOperations(n int) []int {
	if n == 1 {
		return []int{1}
	}

	dp := make([]int, n+1)
	parent := make([]int, n+1)

	for value := 2; value <= n; value++ {
		dp[value] = dp[value-1] + 1
		parent[value] = value - 1

		if value%2 == 0 && dp[value/2]+1 < dp[value] {
			dp[value] = dp[value/2] + 1
			parent[value] = value / 2
		}
		if value%3 == 0 && dp[value/3]+1 < dp[value] {
			dp[value] = dp[value/3] + 1
			parent[value] = value / 3
		}
	}

	sequence := make([]int, 0, dp[n]+1)
	for current := n; current > 0; current = parent[current] {
		sequence = append(sequence, current)
	}

	for left, right := 0, len(sequence)-1; left < right; left, right = left+1, right-1 {
		sequence[left], sequence[right] = sequence[right], sequence[left]
	}

	return sequence
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	sequence := computeOperations(n)

	fmt.Println(len(sequence) - 1)
	for i, value := range sequence {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(value)
	}
	fmt.Println()
}
