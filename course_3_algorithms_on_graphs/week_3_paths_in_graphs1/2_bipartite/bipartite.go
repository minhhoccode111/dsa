// go run bipartite.go

// Check whether an undirected graph is bipartite.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func isBipartite(adj [][]int) bool {
	color := make([]int, len(adj))
	for i := range color {
		color[i] = -1
	}

	for start := range adj {
		if color[start] != -1 {
			continue
		}

		color[start] = 0
		queue := []int{start}

		for head := 0; head < len(queue); head++ {
			u := queue[head]

			for _, v := range adj[u] {
				if color[v] == -1 {
					color[v] = 1 - color[u]
					queue = append(queue, v)
				} else if color[v] == color[u] {
					return false
				}
			}
		}
	}

	return true
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

	if isBipartite(adj) {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}
