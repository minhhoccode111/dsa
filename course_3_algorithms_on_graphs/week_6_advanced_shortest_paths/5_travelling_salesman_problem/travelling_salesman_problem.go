// go run travelling_salesman_problem.go

// Compute the shortest tour from a depot through a set of stores and back.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/internal/search"
)

const inf int64 = 1 << 60

type tourEdge struct {
	from, to int
	weight   int64
}

func buildTourGraph(n int, edges []tourEdge) ([][]search.Edge, [][]search.Edge) {
	adj := make([][]search.Edge, n)
	reverse := make([][]search.Edge, n)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], search.Edge{To: e.to, Weight: e.weight})
		reverse[e.to] = append(reverse[e.to], search.Edge{To: e.from, Weight: e.weight})
	}
	return adj, reverse
}

type tourSolver struct {
	ch *contractionHierarchy
}

func newTourSolver(n int, edges []tourEdge) *tourSolver {
	return &tourSolver{ch: newContractionHierarchy(n, edges)}
}

func (s *tourSolver) minTour(stores []int) int64 {
	return minTourWithDistance(s.ch.distance, stores)
}

func minTour(adj, reverse [][]search.Edge, stores []int) int64 {
	return minTourWithDistance(func(from, to int) int64 {
		return search.Distance(adj, reverse, from, to)
	}, stores)
}

func minTourWithDistance(distance func(from, to int) int64, stores []int) int64 {
	count := len(stores)
	if count <= 1 {
		return 0
	}

	dist := make([][]int64, count)
	for i := 0; i < count; i++ {
		dist[i] = make([]int64, count)
		for j := 0; j < count; j++ {
			d := distance(stores[i], stores[j])
			if d == -1 {
				d = inf
			}
			dist[i][j] = d
		}
	}

	full := 1 << count
	dp := make([]int64, full*count)
	for i := range dp {
		dp[i] = inf
	}
	dp[1*count] = 0

	for mask := 1; mask < full; mask++ {
		if mask&1 == 0 {
			continue
		}
		for last := 0; last < count; last++ {
			if mask&(1<<last) == 0 {
				continue
			}
			current := dp[mask*count+last]
			if current == inf {
				continue
			}
			for next := 1; next < count; next++ {
				if mask&(1<<next) != 0 || dist[last][next] == inf {
					continue
				}
				nextMask := mask | 1<<next
				candidate := current + dist[last][next]
				if candidate < dp[nextMask*count+next] {
					dp[nextMask*count+next] = candidate
				}
			}
		}
	}

	best := inf
	for last := 1; last < count; last++ {
		if dp[(full-1)*count+last] == inf || dist[last][0] == inf {
			continue
		}
		candidate := dp[(full-1)*count+last] + dist[last][0]
		if candidate < best {
			best = candidate
		}
	}

	if best == inf {
		return -1
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	edges := make([]tourEdge, m)
	for i := range edges {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return
		}
		edges[i] = tourEdge{from: u - 1, to: v - 1, weight: weight}
	}

	solver := newTourSolver(n, edges)
	fmt.Println("Ready")

	var queries int
	if _, err := fmt.Fscan(reader, &queries); err != nil {
		return
	}

	for i := 0; i < queries; i++ {
		var count int
		if _, err := fmt.Fscan(reader, &count); err != nil {
			return
		}
		stores := make([]int, count)
		for j := range stores {
			var store int
			if _, err := fmt.Fscan(reader, &store); err != nil {
				return
			}
			stores[j] = store - 1
		}
		fmt.Println(solver.minTour(stores))
	}
}
