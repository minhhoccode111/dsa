// go run bfs.go

// Find the shortest path length between two vertices of an undirected graph.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func shortestDistance(adj [][]int, s, t int) int {
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = -1
	}
	dist[s] = 0

	queue := make([]int, 1, len(adj))
	queue[0] = s

	for head := 0; head < len(queue); head++ {
		u := queue[head]

		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				if v == t {
					return dist[v]
				}
				queue = append(queue, v)
			}
		}
	}

	return dist[t]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		if _, err := fmt.Fscan(reader, &a, &b); err != nil {
			return
		}
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	var s, t int
	if _, err := fmt.Fscan(reader, &s, &t); err != nil {
		return
	}

	fmt.Println(shortestDistance(adj, s-1, t-1))
}
