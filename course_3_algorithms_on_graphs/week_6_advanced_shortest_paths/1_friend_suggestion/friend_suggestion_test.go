package main

import (
	"container/heap"
	"math/rand"
	"testing"
)

func TestBidirectionalDistance(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]weightedEdge
		rev  [][]weightedEdge
		s, t int
		want int64
	}{
		{"same vertex", [][]weightedEdge{{}}, [][]weightedEdge{{}}, 0, 0, 0},
		{"direct edge", [][]weightedEdge{{{1, 1}}, {}}, [][]weightedEdge{{}, {{0, 1}}}, 0, 1, 1},
		{"reverse unreachable", [][]weightedEdge{{{1, 1}}, {}}, [][]weightedEdge{{}, {{0, 1}}}, 1, 0, -1},
		{"unreachable", [][]weightedEdge{{}, {}}, [][]weightedEdge{{}, {}}, 0, 1, -1},
		{"fixture 1", [][]weightedEdge{{{1, 1}}, {}}, [][]weightedEdge{{}, {{0, 1}}}, 0, 1, 1},
		{"fixture 2", [][]weightedEdge{{{1, 1}, {2, 5}}, {{2, 2}}, {}, {{0, 2}}}, [][]weightedEdge{{}, {{0, 1}}, {{0, 5}, {1, 2}}, {}}, 0, 2, 3},
		{"prefers alternate route", [][]weightedEdge{{{1, 10}, {2, 2}}, {{3, 1}}, {{3, 10}}, {}}, [][]weightedEdge{{}, {{0, 10}}, {{0, 2}}, {{1, 1}, {2, 10}}}, 0, 3, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bidirectionalDistance(tt.adj, tt.rev, tt.s, tt.t); got != tt.want {
				t.Fatalf("bidirectionalDistance(...) = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBidirectionalDistanceMatchesDijkstra(t *testing.T) {
	rng := rand.New(rand.NewSource(424242))

	for trial := 0; trial < 100; trial++ {
		n := rng.Intn(30) + 1
		adj := make([][]weightedEdge, n)
		rev := make([][]weightedEdge, n)

		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if u == v || rng.Intn(100) >= 15 {
					continue
				}
				w := int64(rng.Intn(20) + 1)
				adj[u] = append(adj[u], weightedEdge{v, w})
				rev[v] = append(rev[v], weightedEdge{u, w})
			}
		}

		for s := 0; s < n; s++ {
			want := referenceDijkstra(adj, s)
			for target := 0; target < n; target++ {
				if got := bidirectionalDistance(adj, rev, s, target); got != want[target] {
					t.Fatalf("trial %d: distance(%d, %d) = %d, want %d", trial, s, target, got, want[target])
				}
			}
		}
	}
}

func referenceDijkstra(adj [][]weightedEdge, s int) []int64 {
	dist := make([]int64, len(adj))
	for i := range dist {
		dist[i] = inf
	}
	dist[s] = 0

	pq := &referenceQueue{{distance: 0, vertex: s}}
	heap.Init(pq)
	for pq.Len() > 0 {
		item := heap.Pop(pq).(referenceItem)
		if item.distance != dist[item.vertex] {
			continue
		}
		for _, e := range adj[item.vertex] {
			next := item.distance + e.weight
			if next < dist[e.to] {
				dist[e.to] = next
				heap.Push(pq, referenceItem{next, e.to})
			}
		}
	}
	for i := range dist {
		if dist[i] == inf {
			dist[i] = -1
		}
	}
	return dist
}

type referenceItem struct {
	distance int64
	vertex   int
}

type referenceQueue []referenceItem

func (q referenceQueue) Len() int            { return len(q) }
func (q referenceQueue) Less(i, j int) bool  { return q[i].distance < q[j].distance }
func (q referenceQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *referenceQueue) Push(x interface{}) { *q = append(*q, x.(referenceItem)) }
func (q *referenceQueue) Pop() interface{} {
	old := *q
	item := old[len(old)-1]
	*q = old[:len(old)-1]
	return item
}
