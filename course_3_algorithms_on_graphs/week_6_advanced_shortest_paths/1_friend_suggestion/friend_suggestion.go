// go run friend_suggestion.go

// Answer shortest-path queries with bidirectional Dijkstra.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type weightedEdge struct {
	to     int
	weight int64
}

type biItem struct {
	distance int64
	vertex   int
}

type biQueue []biItem

func (q biQueue) Len() int            { return len(q) }
func (q biQueue) Less(i, j int) bool  { return q[i].distance < q[j].distance }
func (q biQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *biQueue) Push(x interface{}) { *q = append(*q, x.(biItem)) }
func (q *biQueue) Pop() interface{} {
	old := *q
	item := old[len(old)-1]
	*q = old[:len(old)-1]
	return item
}

func cleanTop(q *biQueue, dist []int64) (int64, bool) {
	for q.Len() > 0 {
		item := (*q)[0]
		if item.distance == dist[item.vertex] {
			return item.distance, true
		}
		heap.Pop(q)
	}
	return inf, false
}

func bidirectionalDistance(adj, reverse [][]weightedEdge, s, t int) int64 {
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

	pqF := &biQueue{{distance: 0, vertex: s}}
	pqB := &biQueue{{distance: 0, vertex: t}}
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
			item := heap.Pop(pqF).(biItem)
			u := item.vertex

			if distB[u] != inf && distF[u]+distB[u] < best {
				best = distF[u] + distB[u]
			}

			for _, e := range adj[u] {
				next := distF[u] + e.weight
				if next < distF[e.to] {
					distF[e.to] = next
					heap.Push(pqF, biItem{distance: next, vertex: e.to})
				}
			}
			continue
		}

		item := heap.Pop(pqB).(biItem)
		u := item.vertex

		if distF[u] != inf && distF[u]+distB[u] < best {
			best = distF[u] + distB[u]
		}

		for _, e := range reverse[u] {
			next := distB[u] + e.weight
			if next < distB[e.to] {
				distB[e.to] = next
				heap.Push(pqB, biItem{distance: next, vertex: e.to})
			}
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

	adj := make([][]weightedEdge, n)
	reverse := make([][]weightedEdge, n)
	for i := 0; i < m; i++ {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return
		}
		u--
		v--
		adj[u] = append(adj[u], weightedEdge{to: v, weight: weight})
		reverse[v] = append(reverse[v], weightedEdge{to: u, weight: weight})
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
		fmt.Println(bidirectionalDistance(adj, reverse, s-1, t-1))
	}
}
