// go run toposort.go

// Compute a topological ordering of a directed acyclic graph.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type dfsFrame struct {
	vertex int
	next   int
}

func toposort(adj [][]int) []int {
	used := make([]bool, len(adj))
	postorder := make([]int, 0, len(adj))

	for start := range adj {
		if used[start] {
			continue
		}

		used[start] = true
		stack := []dfsFrame{{vertex: start}}

		for len(stack) > 0 {
			frame := &stack[len(stack)-1]

			if frame.next < len(adj[frame.vertex]) {
				neighbor := adj[frame.vertex][frame.next]
				frame.next++
				if !used[neighbor] {
					used[neighbor] = true
					stack = append(stack, dfsFrame{vertex: neighbor})
				}
				continue
			}

			postorder = append(postorder, frame.vertex)
			stack = stack[:len(stack)-1]
		}
	}

	order := make([]int, len(postorder))
	for i, v := range postorder {
		order[len(postorder)-1-i] = v
	}
	return order
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
		adj[a-1] = append(adj[a-1], b-1)
	}

	order := toposort(adj)
	for i, v := range order {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v + 1)
	}
	fmt.Println()
}
