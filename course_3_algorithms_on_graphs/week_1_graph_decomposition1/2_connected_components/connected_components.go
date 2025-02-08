// go run connected_components.go

// Count connected components in an undirected graph.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func countComponents(adj [][]int) int {
	visited := make([]bool, len(adj))
	components := 0

	for start := range adj {
		if visited[start] {
			continue
		}

		components++
		visited[start] = true
		stack := []int{start}

		for len(stack) > 0 {
			vertex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, neighbor := range adj[vertex] {
				if !visited[neighbor] {
					visited[neighbor] = true
					stack = append(stack, neighbor)
				}
			}
		}
	}

	return components
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

	fmt.Println(countComponents(adj))
}
