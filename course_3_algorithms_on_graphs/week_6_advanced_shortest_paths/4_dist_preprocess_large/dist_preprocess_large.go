// go run dist_preprocess_large.go

// Answer shortest-path queries on a directed graph with non-negative weights.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type distanceEdge struct {
	from, to int
	weight   int64
}

type distanceQuery struct {
	from, to int
}

func buildSearchGraph(n int, edges []distanceEdge) ([][]Edge, [][]Edge) {
	adj := make([][]Edge, n)
	reverse := make([][]Edge, n)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], Edge{To: e.to, Weight: e.weight})
		reverse[e.to] = append(reverse[e.to], Edge{To: e.from, Weight: e.weight})
	}
	return adj, reverse
}

func distanceQueries(n int, edges []distanceEdge, queries []distanceQuery) []int64 {
	adj, reverse := buildSearchGraph(n, edges)
	answers := make([]int64, len(queries))
	for i, query := range queries {
		answers[i] = Distance(adj, reverse, query.from, query.to)
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
		fmt.Println(Distance(adj, reverse, s-1, t-1))
	}
}

// The helpers below are inlined from dsa/course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/internal/search
// so that this file can be submitted to Coursera as a single Go file.

const inf int64 = 1 << 60

// Edge is a directed weighted edge.
type Edge struct {
	To     int
	Weight int64
}

type item struct {
	distance int64
	vertex   int
}

type queue []item

func (q queue) Len() int            { return len(q) }
func (q queue) Less(i, j int) bool  { return q[i].distance < q[j].distance }
func (q queue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *queue) Push(x interface{}) { *q = append(*q, x.(item)) }
func (q *queue) Pop() interface{} {
	old := *q
	result := old[len(old)-1]
	*q = old[:len(old)-1]
	return result
}

func cleanTop(q *queue, dist []int64) (int64, bool) {
	for q.Len() > 0 {
		next := (*q)[0]
		if next.distance == dist[next.vertex] {
			return next.distance, true
		}
		heap.Pop(q)
	}
	return inf, false
}

// Distance returns the shortest path weight from s to t, or -1 when t is not
// reachable from s. adj holds forward edges and reverse holds reversed edges.
func Distance(adj, reverse [][]Edge, s, t int) int64 {
	if s == t {
		return 0
	}

	n := len(adj)
	distF := make([]int64, n)
	distB := make([]int64, n)
	for i := 0; i < n; i++ {
		distF[i] = inf
		distB[i] = inf
	}
	distF[s] = 0
	distB[t] = 0

	pqF := &queue{{distance: 0, vertex: s}}
	pqB := &queue{{distance: 0, vertex: t}}
	heap.Init(pqF)
	heap.Init(pqB)

	best := inf

	for {
		topF, okF := cleanTop(pqF, distF)
		topB, okB := cleanTop(pqB, distB)

		if !okF && !okB {
			break
		}
		if okF && okB && topF+topB >= best {
			break
		}
		if okF && !okB && topF >= best {
			break
		}
		if okB && !okF && topB >= best {
			break
		}

		if okF && (!okB || topF <= topB) {
			next := heap.Pop(pqF).(item)
			u := next.vertex

			if distB[u] != inf && distF[u]+distB[u] < best {
				best = distF[u] + distB[u]
			}

			for _, e := range adj[u] {
				candidate := distF[u] + e.Weight
				if candidate < distF[e.To] {
					distF[e.To] = candidate
					heap.Push(pqF, item{distance: candidate, vertex: e.To})
				}
			}
			continue
		}

		next := heap.Pop(pqB).(item)
		u := next.vertex

		if distF[u] != inf && distF[u]+distB[u] < best {
			best = distF[u] + distB[u]
		}

		for _, e := range reverse[u] {
			candidate := distB[u] + e.Weight
			if candidate < distB[e.To] {
				distB[e.To] = candidate
				heap.Push(pqB, item{distance: candidate, vertex: e.To})
			}
		}
	}

	if best == inf {
		return -1
	}
	return best
}
