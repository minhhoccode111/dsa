package main

import (
	"container/heap"
	"math"
	"math/rand"
	"testing"
)

func TestAStarDistance(t *testing.T) {
	tests := []struct {
		name string
		x, y []int64
		adj  [][]weightedEdge
		s, t int
		want int64
	}{
		{
			name: "same vertex",
			x:    []int64{0},
			y:    []int64{0},
			adj:  [][]weightedEdge{{}},
			s:    0, t: 0, want: 0,
		},
		{
			name: "fixture 01",
			x:    []int64{0, 0},
			y:    []int64{0, 1},
			adj:  [][]weightedEdge{{{1, 1}}, {}},
			s:    0, t: 1, want: 1,
		},
		{
			name: "unreachable",
			x:    []int64{0, 10},
			y:    []int64{0, 10},
			adj:  [][]weightedEdge{{}, {}},
			s:    0, t: 1, want: -1,
		},
		{
			name: "detour beats direct",
			x:    []int64{0, 5, 10, 0},
			y:    []int64{0, 0, 0, 0},
			adj:  [][]weightedEdge{{{1, 5}, {2, 100}}, {{2, 5}}, {{3, 10}}, {}},
			s:    0, t: 3, want: 20,
		},
		{
			name: "direct edge",
			x:    []int64{0, 3},
			y:    []int64{0, 4},
			adj:  [][]weightedEdge{{{1, 5}}, {}},
			s:    0, t: 1, want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aStarDistance(tt.x, tt.y, tt.adj, tt.s, tt.t); got != tt.want {
				t.Fatalf("aStarDistance(...) = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAStarDistanceMatchesDijkstra(t *testing.T) {
	rng := rand.New(rand.NewSource(848484))

	for trial := 0; trial < 100; trial++ {
		n := rng.Intn(20) + 1
		x := make([]int64, n)
		y := make([]int64, n)
		for i := 0; i < n; i++ {
			x[i] = int64(rng.Intn(101) - 50)
			y[i] = int64(rng.Intn(101) - 50)
		}

		adj := make([][]weightedEdge, n)
		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if u == v || rng.Intn(100) >= 20 {
					continue
				}
				dx := float64(x[u] - x[v])
				dy := float64(y[u] - y[v])
				// Keep the Euclidean heuristic consistent.
				w := int64(math.Ceil(math.Hypot(dx, dy))) + int64(rng.Intn(10))
				adj[u] = append(adj[u], weightedEdge{v, w})
			}
		}

		for s := 0; s < n; s++ {
			want := referenceDijkstra(adj, s)
			for target := 0; target < n; target++ {
				if got := aStarDistance(x, y, adj, s, target); got != want[target] {
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
