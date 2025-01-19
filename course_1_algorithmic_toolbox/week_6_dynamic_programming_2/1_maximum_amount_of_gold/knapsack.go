// go run knapsack.go

// Given gold bars of various weights and a backpack of capacity W, place as
// much gold as possible into the backpack.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func maximumGold(capacity int, weights []int) int {
	dp := make([]int, capacity+1)

	for _, weight := range weights {
		for current := capacity; current >= weight; current-- {
			if candidate := dp[current-weight] + weight; candidate > dp[current] {
				dp[current] = candidate
			}
		}
	}

	return dp[capacity]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var capacity, n int
	if _, err := fmt.Fscan(reader, &capacity, &n); err != nil {
		return
	}

	weights := make([]int, n)
	for i := range weights {
		if _, err := fmt.Fscan(reader, &weights[i]); err != nil {
			return
		}
	}

	fmt.Println(maximumGold(capacity, weights))
}
