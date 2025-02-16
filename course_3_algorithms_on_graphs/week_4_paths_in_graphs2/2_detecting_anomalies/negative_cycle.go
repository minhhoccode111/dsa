// go run negative_cycle.go

// Check whether a directed weighted graph contains a negative cycle.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func hasNegativeCycle(adj [][]int, cost [][]int64) bool {
	n := len(adj)
	dist := make([]int64, n)

	for iteration := 0; iteration < n; iteration++ {
		updated := false

		for u := 0; u < n; u++ {
			for i, v := range adj[u] {
				candidate := dist[u] + cost[u][i]
				if candidate < dist[v] {
					dist[v] = candidate
					updated = true
				}
			}
		}

		if !updated {
			return false
		}
		if iteration == n-1 {
			return true
		}
	}

	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	cost := make([][]int64, n)
	for i := 0; i < m; i++ {
		var a, b int
		var w int64
		if _, err := fmt.Fscan(reader, &a, &b, &w); err != nil {
			return
		}
		a--
		b--
		adj[a] = append(adj[a], b)
		cost[a] = append(cost[a], w)
	}

	if hasNegativeCycle(adj, cost) {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}
