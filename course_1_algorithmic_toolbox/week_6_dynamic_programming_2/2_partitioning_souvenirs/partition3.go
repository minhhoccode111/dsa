// go run partition3.go

// Check whether the values can be partitioned into three subsets with equal sums.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func partition3(values []int) int {
	if len(values) < 3 {
		return 0
	}

	total := 0
	for _, value := range values {
		total += value
	}
	if total%3 != 0 {
		return 0
	}

	target := total / 3
	n := len(values)

	dp := make([][][]bool, n+1)
	for i := range dp {
		dp[i] = make([][]bool, target+1)
		for s1 := range dp[i] {
			dp[i][s1] = make([]bool, target+1)
		}
	}
	dp[0][0][0] = true

	prefix := 0
	for i := 1; i <= n; i++ {
		value := values[i-1]
		for s1 := 0; s1 <= target; s1++ {
			for s2 := 0; s2 <= target; s2++ {
				if !dp[i-1][s1][s2] {
					continue
				}

				if s1+value <= target {
					dp[i][s1+value][s2] = true
				}
				if s2+value <= target {
					dp[i][s1][s2+value] = true
				}

				thirdSum := prefix - s1 - s2
				if thirdSum+value <= target {
					dp[i][s1][s2] = true
				}
			}
		}
		prefix += value
	}

	if dp[n][target][target] {
		return 1
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	values := make([]int, n)
	for i := range values {
		if _, err := fmt.Fscan(reader, &values[i]); err != nil {
			return
		}
	}

	fmt.Println(partition3(values))
}
