// Package search provides exact shortest-path queries with bidirectional
// Dijkstra for non-negative directed graphs.
package search

import "container/heap"

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
