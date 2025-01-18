// go run edit_distance.go

// Compute the edit distance between two strings.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func editDistance(s, t string) int {
	n, m := len(s), len(t)
	dp := make([]int, (n+1)*(m+1))
	index := func(i, j int) int { return i*(m+1) + j }

	for i := 1; i <= n; i++ {
		dp[index(i, 0)] = i
	}
	for j := 1; j <= m; j++ {
		dp[index(0, j)] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if s[i-1] == t[j-1] {
				cost = 0
			}

			best := dp[index(i-1, j)] + 1
			if candidate := dp[index(i, j-1)] + 1; candidate < best {
				best = candidate
			}
			if candidate := dp[index(i-1, j-1)] + cost; candidate < best {
				best = candidate
			}
			dp[index(i, j)] = best
		}
	}

	return dp[index(n, m)]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var s, t string
	if _, err := fmt.Fscan(reader, &s, &t); err != nil {
		return
	}

	fmt.Println(editDistance(s, t))
}
