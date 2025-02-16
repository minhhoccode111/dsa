// go run dijkstra.go

// Find the shortest path weight between two vertices with non-negative edges.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type queueItem struct {
	distance int64
	vertex   int
}

type minQueue []queueItem

func (q minQueue) Len() int            { return len(q) }
func (q minQueue) Less(i, j int) bool  { return q[i].distance < q[j].distance }
func (q minQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *minQueue) Push(x interface{}) { *q = append(*q, x.(queueItem)) }
func (q *minQueue) Pop() interface{} {
	old := *q
	item := old[len(old)-1]
	*q = old[:len(old)-1]
	return item
}

func distance(adj [][]int, cost [][]int64, s, t int) int64 {
	dist := make([]int64, len(adj))
	for i := range dist {
		dist[i] = inf
	}
	dist[s] = 0

	pq := &minQueue{{distance: 0, vertex: s}}
	heap.Init(pq)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(queueItem)
		u := item.vertex
		if item.distance != dist[u] {
			continue
		}

		for i, v := range adj[u] {
			next := dist[u] + cost[u][i]
			if next < dist[v] {
				dist[v] = next
				heap.Push(pq, queueItem{distance: next, vertex: v})
			}
		}
	}

	if dist[t] == inf {
		return -1
	}
	return dist[t]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	cost := make([][]int64, n)
	for i := 0; i < m; i++ {
		var a, b int
		var w int64
		if _, err := fmt.Fscan(reader, &a, &b, &w); err != nil {
			return
		}
		a--
		b--
		adj[a] = append(adj[a], b)
		cost[a] = append(cost[a], w)
	}

	var s, t int
	if _, err := fmt.Fscan(reader, &s, &t); err != nil {
		return
	}

	fmt.Println(distance(adj, cost, s-1, t-1))
}
