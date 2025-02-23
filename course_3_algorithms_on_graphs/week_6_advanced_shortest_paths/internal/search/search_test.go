package search

import (
	"container/heap"
	"math/rand"
	"testing"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]Edge
		rev  [][]Edge
		s, t int
		want int64
	}{
		{"same vertex", [][]Edge{{}}, [][]Edge{{}}, 0, 0, 0},
		{"direct edge", [][]Edge{{{1, 1}}, {}}, [][]Edge{{}, {{0, 1}}}, 0, 1, 1},
		{"reverse unreachable", [][]Edge{{{1, 1}}, {}}, [][]Edge{{}, {{0, 1}}}, 1, 0, -1},
		{"unreachable", [][]Edge{{}, {}}, [][]Edge{{}, {}}, 0, 1, -1},
		{"prefers shorter route", [][]Edge{{{1, 10}, {2, 2}}, {{3, 1}}, {{3, 10}}, {}}, [][]Edge{{}, {{0, 10}}, {{0, 2}}, {{1, 1}, {2, 10}}}, 0, 3, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distance(tt.adj, tt.rev, tt.s, tt.t); got != tt.want {
				t.Fatalf("Distance(...) = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDistanceMatchesDijkstra(t *testing.T) {
	rng := rand.New(rand.NewSource(919191))

	for trial := 0; trial < 100; trial++ {
		n := rng.Intn(30) + 1
		adj := make([][]Edge, n)
		rev := make([][]Edge, n)

		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if u == v || rng.Intn(100) >= 15 {
					continue
				}
				w := int64(rng.Intn(20) + 1)
				adj[u] = append(adj[u], Edge{v, w})
				rev[v] = append(rev[v], Edge{u, w})
			}
		}

		for s := 0; s < n; s++ {
			want := dijkstra(adj, s)
			for target := 0; target < n; target++ {
				if got := Distance(adj, rev, s, target); got != want[target] {
					t.Fatalf("trial %d: Distance(%d, %d) = %d, want %d", trial, s, target, got, want[target])
				}
			}
		}
	}
}

func dijkstra(adj [][]Edge, s int) []int64 {
	const inf int64 = 1 << 60
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
			next := item.distance + e.Weight
			if next < dist[e.To] {
				dist[e.To] = next
				heap.Push(pq, referenceItem{next, e.To})
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
