// go run dist_with_coords.go

// Answer shortest-path queries with A* using Euclidean vertex coordinates.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

const inf int64 = 1 << 60

type weightedEdge struct {
	to     int
	weight int64
}

type astarItem struct {
	priority float64
	vertex   int
}

type astarQueue []astarItem

func (q astarQueue) Len() int            { return len(q) }
func (q astarQueue) Less(i, j int) bool  { return q[i].priority < q[j].priority }
func (q astarQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *astarQueue) Push(x interface{}) { *q = append(*q, x.(astarItem)) }
func (q *astarQueue) Pop() interface{} {
	old := *q
	item := old[len(old)-1]
	*q = old[:len(old)-1]
	return item
}

func aStarDistance(x, y []int64, adj [][]weightedEdge, s, t int) int64 {
	if s == t {
		return 0
	}

	heuristic := func(v int) float64 {
		dx := float64(x[v] - x[t])
		dy := float64(y[v] - y[t])
		return math.Hypot(dx, dy)
	}

	dist := make([]int64, len(adj))
	visited := make([]bool, len(adj))
	for i := range dist {
		dist[i] = inf
	}
	dist[s] = 0

	pq := &astarQueue{{priority: heuristic(s), vertex: s}}
	heap.Init(pq)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(astarItem)
		u := item.vertex
		if visited[u] {
			continue
		}
		visited[u] = true
		if u == t {
			return dist[u]
		}

		for _, e := range adj[u] {
			if visited[e.to] {
				continue
			}
			next := dist[u] + e.weight
			if next < dist[e.to] {
				dist[e.to] = next
				heap.Push(pq, astarItem{priority: float64(next) + heuristic(e.to), vertex: e.to})
			}
		}
	}

	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &x[i], &y[i]); err != nil {
			return
		}
	}

	adj := make([][]weightedEdge, n)
	for i := 0; i < m; i++ {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return
		}
		adj[u-1] = append(adj[u-1], weightedEdge{to: v - 1, weight: weight})
	}

	var queries int
	if _, err := fmt.Fscan(reader, &queries); err != nil {
		return
	}

	for i := 0; i < queries; i++ {
		var s, t int
		if _, err := fmt.Fscan(reader, &s, &t); err != nil {
			return
		}
		fmt.Println(aStarDistance(x, y, adj, s-1, t-1))
	}
}
