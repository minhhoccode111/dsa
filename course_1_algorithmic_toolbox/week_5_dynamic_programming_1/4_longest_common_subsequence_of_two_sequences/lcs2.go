// go run lcs2.go

// Compute the maximum length of a common subsequence of two sequences.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func lcs2(a, b []int) int {
	n, m := len(a), len(b)
	dp := make([]int, (n+1)*(m+1))
	index := func(i, j int) int { return i*(m+1) + j }

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[index(i, j)] = dp[index(i-1, j-1)] + 1
				continue
			}

			best := dp[index(i-1, j)]
			if candidate := dp[index(i, j-1)]; candidate > best {
				best = candidate
			}
			dp[index(i, j)] = best
		}
	}

	return dp[index(n, m)]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		if _, err := fmt.Fscan(reader, &a[i]); err != nil {
			return
		}
	}

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	b := make([]int, m)
	for i := range b {
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return
		}
	}

	fmt.Println(lcs2(a, b))
}
