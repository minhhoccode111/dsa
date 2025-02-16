// go run shortest_paths.go

// Find shortest paths from a source in a graph that may contain negative
// edges and negative cycles.

package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func shortestPaths(adj [][]int, cost [][]int64, s int) ([]int64, []bool, []bool) {
	n := len(adj)
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[s] = 0

	// Bellman-Ford relaxations.
	for iteration := 0; iteration < n-1; iteration++ {
		updated := false

		for u := 0; u < n; u++ {
			if dist[u] == inf {
				continue
			}
			for i, v := range adj[u] {
				candidate := dist[u] + cost[u][i]
				if candidate < dist[v] {
					dist[v] = candidate
					updated = true
				}
			}
		}

		if !updated {
			break
		}
	}

	reachable := make([]bool, n)
	for v := range reachable {
		reachable[v] = dist[v] != inf
	}
	shortest := make([]bool, n)
	for v := range shortest {
		shortest[v] = true
	}

	// Any vertex whose distance can still be improved lies on or after a
	// negative cycle reachable from the source.
	queue := make([]int, 0, n)
	for u := 0; u < n; u++ {
		if dist[u] == inf {
			continue
		}
		for i, v := range adj[u] {
			if dist[u]+cost[u][i] < dist[v] && shortest[v] {
				shortest[v] = false
				queue = append(queue, v)
			}
		}
	}

	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range adj[u] {
			if shortest[v] {
				shortest[v] = false
				queue = append(queue, v)
			}
		}
	}

	return dist, reachable, shortest
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

	var s int
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	dist, reachable, shortest := shortestPaths(adj, cost, s-1)
	for v := 0; v < n; v++ {
		switch {
		case !reachable[v]:
			fmt.Println("*")
		case !shortest[v]:
			fmt.Println("-")
		default:
			fmt.Println(dist[v])
		}
	}
}
