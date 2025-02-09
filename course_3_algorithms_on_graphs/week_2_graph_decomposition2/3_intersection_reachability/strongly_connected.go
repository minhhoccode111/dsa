// go run strongly_connected.go

// Count strongly connected components using Kosaraju's algorithm.

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

func finishOrder(adj [][]int) []int {
	used := make([]bool, len(adj))
	order := make([]int, 0, len(adj))

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

			order = append(order, frame.vertex)
			stack = stack[:len(stack)-1]
		}
	}

	return order
}

func reverseGraph(adj [][]int) [][]int {
	reversed := make([][]int, len(adj))
	for u, neighbors := range adj {
		for _, v := range neighbors {
			reversed[v] = append(reversed[v], u)
		}
	}
	return reversed
}

func countSCCs(adj [][]int) int {
	order := finishOrder(adj)
	reversed := reverseGraph(adj)
	visited := make([]bool, len(adj))
	count := 0

	for i := len(order) - 1; i >= 0; i-- {
		start := order[i]
		if visited[start] {
			continue
		}

		count++
		visited[start] = true
		stack := []int{start}

		for len(stack) > 0 {
			vertex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for _, neighbor := range reversed[vertex] {
				if !visited[neighbor] {
					visited[neighbor] = true
					stack = append(stack, neighbor)
				}
			}
		}
	}

	return count
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

	fmt.Println(countSCCs(adj))
}
