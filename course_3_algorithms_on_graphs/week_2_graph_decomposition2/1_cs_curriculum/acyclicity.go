// go run acyclicity.go

// Check whether a directed graph contains a cycle.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func hasCycle(adj [][]int) bool {
	inDegree := make([]int, len(adj))
	for _, neighbors := range adj {
		for _, v := range neighbors {
			inDegree[v]++
		}
	}

	queue := make([]int, 0, len(adj))
	for v, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, v)
		}
	}

	processed := 0
	for head := 0; head < len(queue); head++ {
		vertex := queue[head]
		processed++

		for _, neighbor := range adj[vertex] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return processed != len(adj)
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

	if hasCycle(adj) {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}
