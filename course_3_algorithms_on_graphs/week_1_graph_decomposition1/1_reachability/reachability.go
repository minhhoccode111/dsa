// go run reachability.go

// Check whether two vertices are connected in an undirected graph.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func reach(adj [][]int, x, y int) bool {
	visited := make([]bool, len(adj))
	stack := []int{x}
	visited[x] = true

	for len(stack) > 0 {
		vertex := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if vertex == y {
			return true
		}

		for _, neighbor := range adj[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				stack = append(stack, neighbor)
			}
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

	var x, y int
	if _, err := fmt.Fscan(reader, &x, &y); err != nil {
		return
	}

	if reach(adj, x-1, y-1) {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}
