// go run lcs3.go

// Compute the maximum length of a common subsequence of three sequences.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func lcs3(a, b, c []int) int {
	n, m, l := len(a), len(b), len(c)
	stride := (m + 1) * (l + 1)
	dp := make([]int, (n+1)*stride)
	index := func(i, j, k int) int { return i*stride + j*(l+1) + k }

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			for k := 1; k <= l; k++ {
				if a[i-1] == b[j-1] && a[i-1] == c[k-1] {
					dp[index(i, j, k)] = dp[index(i-1, j-1, k-1)] + 1
					continue
				}

				best := dp[index(i-1, j, k)]
				if candidate := dp[index(i, j-1, k)]; candidate > best {
					best = candidate
				}
				if candidate := dp[index(i, j, k-1)]; candidate > best {
					best = candidate
				}
				dp[index(i, j, k)] = best
			}
		}
	}

	return dp[index(n, m, l)]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	readSequence := func() ([]int, bool) {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, false
		}
		sequence := make([]int, n)
		for i := range sequence {
			if _, err := fmt.Fscan(reader, &sequence[i]); err != nil {
				return nil, false
			}
		}
		return sequence, true
	}

	a, ok := readSequence()
	if !ok {
		return
	}
	b, ok := readSequence()
	if !ok {
		return
	}
	c, ok := readSequence()
	if !ok {
		return
	}

	fmt.Println(lcs3(a, b, c))
}
