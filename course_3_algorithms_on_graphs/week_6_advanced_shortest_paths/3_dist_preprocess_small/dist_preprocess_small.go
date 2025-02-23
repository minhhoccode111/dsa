// go run dist_preprocess_small.go

// Answer shortest-path queries on a directed graph with non-negative weights.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/internal/search"
)

type distanceEdge struct {
	from, to int
	weight   int64
}

type distanceQuery struct {
	from, to int
}

func buildSearchGraph(n int, edges []distanceEdge) ([][]search.Edge, [][]search.Edge) {
	adj := make([][]search.Edge, n)
	reverse := make([][]search.Edge, n)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], search.Edge{To: e.to, Weight: e.weight})
		reverse[e.to] = append(reverse[e.to], search.Edge{To: e.from, Weight: e.weight})
	}
	return adj, reverse
}

func distanceQueries(n int, edges []distanceEdge, queries []distanceQuery) []int64 {
	adj, reverse := buildSearchGraph(n, edges)
	answers := make([]int64, len(queries))
	for i, query := range queries {
		answers[i] = search.Distance(adj, reverse, query.from, query.to)
	}
	return answers
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	edges := make([]distanceEdge, m)
	for i := range edges {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return
		}
		edges[i] = distanceEdge{from: u - 1, to: v - 1, weight: weight}
	}

	adj, reverse := buildSearchGraph(n, edges)
	fmt.Println("Ready")

	var queries int
	if _, err := fmt.Fscan(reader, &queries); err != nil {
		return
	}

	for i := 0; i < queries; i++ {
		var s, t int
		if _, err := fmt.Fscan(reader, &s, &t); err != nil {
			return
		}
		fmt.Println(search.Distance(adj, reverse, s-1, t-1))
	}
}
