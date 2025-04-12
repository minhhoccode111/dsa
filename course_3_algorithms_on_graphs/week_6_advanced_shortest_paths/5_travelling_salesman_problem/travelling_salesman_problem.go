// go run travelling_salesman_problem.go

// Compute the shortest tour from a depot through a set of stores and back.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type tourEdge struct {
	from, to int
	weight   int64
}

func buildTourGraph(n int, edges []tourEdge) ([][]Edge, [][]Edge) {
	adj := make([][]Edge, n)
	reverse := make([][]Edge, n)
	for _, e := range edges {
		adj[e.from] = append(adj[e.from], Edge{To: e.to, Weight: e.weight})
		reverse[e.to] = append(reverse[e.to], Edge{To: e.from, Weight: e.weight})
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

func minTour(adj, reverse [][]Edge, stores []int) int64 {
	return minTourWithDistance(func(from, to int) int64 {
		return Distance(adj, reverse, from, to)
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

// The helpers below are inlined from dsa/course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/internal/search
// so that this file can be submitted to Coursera as a single Go file.

const searchInf int64 = 1 << 60

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
	return searchInf, false
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
		distF[i] = searchInf
		distB[i] = searchInf
	}
	distF[s] = 0
	distB[t] = 0

	pqF := &queue{{distance: 0, vertex: s}}
	pqB := &queue{{distance: 0, vertex: t}}
	heap.Init(pqF)
	heap.Init(pqB)

	best := searchInf

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

			if distB[u] != searchInf && distF[u]+distB[u] < best {
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

		if distF[u] != searchInf && distF[u]+distB[u] < best {
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

	if best == searchInf {
		return -1
	}
	return best
}

// The contraction-hierarchy implementation lives in this file so that the
// TSP solution can be submitted to Coursera as a single Go file.

const (
	contractSettleLimit = 100
	contractHopLimit    = 50
)

type chArc struct {
	to     int
	weight int64
}

type contractionHierarchy struct {
	n     int
	out   [][]chArc
	in    [][]chArc
	level []int

	distForward  []int64
	distBackward []int64
	seenForward  []int32
	seenBackward []int32
	queryStamp   int32
	forwardHeap  chQueryHeap
	backwardHeap chQueryHeap
}

func newContractionHierarchy(n int, edges []tourEdge) *contractionHierarchy {
	g := &contractionHierarchy{
		n:            n,
		out:          make([][]chArc, n),
		in:           make([][]chArc, n),
		level:        make([]int, n),
		distForward:  make([]int64, n),
		distBackward: make([]int64, n),
		seenForward:  make([]int32, n),
		seenBackward: make([]int32, n),
	}
	for _, edge := range edges {
		g.addArc(edge.from, edge.to, edge.weight)
	}

	builder := &chBuilder{
		graph:      g,
		contracted: make([]bool, n),
		dist:       make([]int64, n),
		seen:       make([]int32, n),
	}
	builder.build()
	return g
}

func (g *contractionHierarchy) addArc(from, to int, weight int64) {
	for i := range g.out[from] {
		if g.out[from][i].to != to {
			continue
		}
		if weight >= g.out[from][i].weight {
			return
		}
		g.out[from][i].weight = weight
		for j := range g.in[to] {
			if g.in[to][j].to == from {
				g.in[to][j].weight = weight
				break
			}
		}
		return
	}
	g.out[from] = append(g.out[from], chArc{to: to, weight: weight})
	g.in[to] = append(g.in[to], chArc{to: from, weight: weight})
}

func (g *contractionHierarchy) hasEdgeAtMost(from, to int, limit int64) bool {
	for _, edge := range g.out[from] {
		if edge.to == to && edge.weight <= limit {
			return true
		}
	}
	return false
}

func (g *contractionHierarchy) distance(source, target int) int64 {
	if source == target {
		return 0
	}
	g.queryStamp++
	if g.queryStamp == 0 {
		for i := range g.seenForward {
			g.seenForward[i] = 0
			g.seenBackward[i] = 0
		}
		g.queryStamp = 1
	}
	stamp := g.queryStamp

	g.seenForward[source] = stamp
	g.distForward[source] = 0
	g.seenBackward[target] = stamp
	g.distBackward[target] = 0

	forward := &g.forwardHeap
	backward := &g.backwardHeap
	forward.reset()
	backward.reset()
	forward.push(chQueryItem{node: source, distance: 0})
	backward.push(chQueryItem{node: target, distance: 0})
	best := inf

	for {
		minForward, minBackward := inf, inf
		if forward.len() > 0 {
			minForward = forward.peek().distance
		}
		if backward.len() > 0 {
			minBackward = backward.peek().distance
		}
		if minForward == inf && minBackward == inf {
			break
		}
		if minForward >= best && minBackward >= best {
			break
		}

		if minForward <= minBackward {
			item := forward.pop()
			if item.distance != g.distForward[item.node] {
				continue
			}
			if g.seenBackward[item.node] == stamp {
				if candidate := item.distance + g.distBackward[item.node]; candidate < best {
					best = candidate
				}
			}
			for _, edge := range g.out[item.node] {
				if g.level[edge.to] <= g.level[item.node] {
					continue
				}
				candidate := item.distance + edge.weight
				if g.seenForward[edge.to] != stamp || candidate < g.distForward[edge.to] {
					g.seenForward[edge.to] = stamp
					g.distForward[edge.to] = candidate
					forward.push(chQueryItem{node: edge.to, distance: candidate})
				}
			}
		} else {
			item := backward.pop()
			if item.distance != g.distBackward[item.node] {
				continue
			}
			if g.seenForward[item.node] == stamp {
				if candidate := item.distance + g.distForward[item.node]; candidate < best {
					best = candidate
				}
			}
			for _, edge := range g.in[item.node] {
				if g.level[edge.to] <= g.level[item.node] {
					continue
				}
				candidate := item.distance + edge.weight
				if g.seenBackward[edge.to] != stamp || candidate < g.distBackward[edge.to] {
					g.seenBackward[edge.to] = stamp
					g.distBackward[edge.to] = candidate
					backward.push(chQueryItem{node: edge.to, distance: candidate})
				}
			}
		}
	}

	if best == inf {
		return -1
	}
	return best
}

type chBuilder struct {
	graph       *contractionHierarchy
	contracted  []bool
	dist        []int64
	seen        []int32
	stamp       int32
	witnessHeap chWitnessHeap
}

func (b *chBuilder) build() {
	order := nestedDissectionOrder(b.graph)
	for rank, node := range order {
		b.contract(node)
		b.contracted[node] = true
		b.graph.level[node] = rank
	}
}

func nestedDissectionOrder(g *contractionHierarchy) []int {
	n := g.n
	adjacency := make([][]int, n)
	for u, edges := range g.out {
		for _, edge := range edges {
			if edge.to == u {
				continue
			}
			adjacency[u] = append(adjacency[u], edge.to)
			adjacency[edge.to] = append(adjacency[edge.to], u)
		}
	}

	order := make([]int, 0, n)
	mark := make([]int32, n)
	visited := make([]int32, n)
	dist := make([]int, n)
	var markStamp, visitStamp int32

	var decompose func(nodes []int)
	decomposeConnected := func(component []int) {
		if len(component) <= 64 {
			order = append(order, component...)
			return
		}
		markStamp++
		if markStamp == 0 {
			for i := range mark {
				mark[i] = 0
			}
			markStamp = 1
		}
		stamp := markStamp
		for _, node := range component {
			mark[node] = stamp
		}

		visitStamp++
		if visitStamp == 0 {
			for i := range visited {
				visited[i] = 0
			}
			visitStamp = 1
		}
		queue := make([]int, 0, len(component))
		start := component[0]
		queue = append(queue, start)
		visited[start] = visitStamp
		dist[start] = 0
		farthest := start
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for _, next := range adjacency[v] {
				if mark[next] != stamp || visited[next] == visitStamp {
					continue
				}
				visited[next] = visitStamp
				dist[next] = dist[v] + 1
				if dist[next] > dist[farthest] {
					farthest = next
				}
				queue = append(queue, next)
			}
		}

		visitStamp++
		if visitStamp == 0 {
			for i := range visited {
				visited[i] = 0
			}
			visitStamp = 1
		}
		queue = queue[:0]
		queue = append(queue, farthest)
		visited[farthest] = visitStamp
		dist[farthest] = 0
		maxLevel := 0
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for _, next := range adjacency[v] {
				if mark[next] != stamp || visited[next] == visitStamp {
					continue
				}
				visited[next] = visitStamp
				dist[next] = dist[v] + 1
				if dist[next] > maxLevel {
					maxLevel = dist[next]
				}
				queue = append(queue, next)
			}
		}

		counts := make([]int, maxLevel+1)
		for _, v := range component {
			counts[dist[v]]++
		}
		total := len(component)
		bestLevel := -1
		bestScore := 1 << 60
		cumulative := 0
		for level := 0; level <= maxLevel; level++ {
			if level > 0 && level < maxLevel {
				score := cumulative + counts[level]/2 - total/2
				if score < 0 {
					score = -score
				}
				if score < bestScore {
					bestScore = score
					bestLevel = level
				}
			}
			cumulative += counts[level]
		}
		if bestLevel < 0 {
			if counts[0] <= counts[maxLevel] {
				bestLevel = 0
			} else {
				bestLevel = maxLevel
			}
		}
		if counts[bestLevel] > total/2 {
			order = append(order, component...)
			return
		}

		lower := make([]int, 0, total)
		upper := make([]int, 0, total)
		separator := make([]int, 0, counts[bestLevel])
		for _, v := range component {
			switch {
			case dist[v] < bestLevel:
				lower = append(lower, v)
			case dist[v] > bestLevel:
				upper = append(upper, v)
			default:
				separator = append(separator, v)
			}
		}
		decompose(lower)
		decompose(upper)
		order = append(order, separator...)
	}

	decompose = func(nodes []int) {
		if len(nodes) == 0 {
			return
		}
		if len(nodes) <= 64 {
			order = append(order, nodes...)
			return
		}
		markStamp++
		if markStamp == 0 {
			for i := range mark {
				mark[i] = 0
			}
			markStamp = 1
		}
		for _, node := range nodes {
			mark[node] = markStamp
		}
		visitStamp++
		if visitStamp == 0 {
			for i := range visited {
				visited[i] = 0
			}
			visitStamp = 1
		}
		components := make([][]int, 0)
		for _, start := range nodes {
			if visited[start] == visitStamp {
				continue
			}
			component := make([]int, 0)
			queue := []int{start}
			visited[start] = visitStamp
			for head := 0; head < len(queue); head++ {
				v := queue[head]
				component = append(component, v)
				for _, next := range adjacency[v] {
					if mark[next] == markStamp && visited[next] != visitStamp {
						visited[next] = visitStamp
						queue = append(queue, next)
					}
				}
			}
			components = append(components, component)
		}
		for _, component := range components {
			decomposeConnected(component)
		}
	}

	decompose(makeRange(n))
	return order
}

func makeRange(n int) []int {
	nodes := make([]int, n)
	for i := range nodes {
		nodes[i] = i
	}
	return nodes
}

func (b *chBuilder) contract(node int) {
	for _, incoming := range b.graph.in[node] {
		u := incoming.to
		if b.contracted[u] {
			continue
		}
		for _, outgoing := range b.graph.out[node] {
			w := outgoing.to
			if b.contracted[w] || w == u {
				continue
			}
			limit := incoming.weight + outgoing.weight
			if b.graph.hasEdgeAtMost(u, w, limit) {
				continue
			}
			if b.hasWitness(u, w, node, limit, contractSettleLimit, contractHopLimit) {
				continue
			}
			b.graph.addArc(u, w, limit)
		}
	}
}

func (b *chBuilder) hasWitness(source, target, avoid int, limit int64, settleLimit, hopLimit int) bool {
	b.stamp++
	if b.stamp == 0 {
		for i := range b.seen {
			b.seen[i] = 0
		}
		b.stamp = 1
	}
	stamp := b.stamp
	b.seen[source] = stamp
	b.dist[source] = 0

	queue := &b.witnessHeap
	queue.reset()
	queue.push(chWitnessItem{node: source, distance: 0})
	settled := 0
	for queue.len() > 0 {
		item := queue.pop()
		if item.distance != b.dist[item.node] {
			continue
		}
		if item.distance > limit {
			return false
		}
		if item.node == target {
			return true
		}
		settled++
		if settled > settleLimit {
			return false
		}
		if item.hops >= hopLimit {
			continue
		}
		for _, edge := range b.graph.out[item.node] {
			if edge.to == avoid || b.contracted[edge.to] {
				continue
			}
			candidate := item.distance + edge.weight
			if candidate > limit {
				continue
			}
			if b.seen[edge.to] != stamp || candidate < b.dist[edge.to] {
				b.seen[edge.to] = stamp
				b.dist[edge.to] = candidate
				queue.push(chWitnessItem{node: edge.to, distance: candidate, hops: item.hops + 1})
			}
		}
	}
	return false
}

type chWitnessItem struct {
	node     int
	distance int64
	hops     int
}

type chWitnessHeap struct {
	items []chWitnessItem
}

func (h *chWitnessHeap) reset()   { h.items = h.items[:0] }
func (h *chWitnessHeap) len() int { return len(h.items) }
func (h *chWitnessHeap) push(item chWitnessItem) {
	h.items = append(h.items, item)
	index := len(h.items) - 1
	for index > 0 {
		parent := (index - 1) / 2
		if !witnessLess(h.items[index], h.items[parent]) {
			break
		}
		h.items[index], h.items[parent] = h.items[parent], h.items[index]
		index = parent
	}
}
func (h *chWitnessHeap) pop() chWitnessItem {
	last := len(h.items) - 1
	result := h.items[0]
	h.items[0] = h.items[last]
	h.items = h.items[:last]
	if last > 0 {
		index := 0
		for {
			left, right := 2*index+1, 2*index+2
			smallest := index
			if left < last && witnessLess(h.items[left], h.items[smallest]) {
				smallest = left
			}
			if right < last && witnessLess(h.items[right], h.items[smallest]) {
				smallest = right
			}
			if smallest == index {
				break
			}
			h.items[index], h.items[smallest] = h.items[smallest], h.items[index]
			index = smallest
		}
	}
	return result
}
func witnessLess(a, b chWitnessItem) bool {
	if a.distance != b.distance {
		return a.distance < b.distance
	}
	return a.node < b.node
}

type chQueryItem struct {
	node     int
	distance int64
}

type chQueryHeap struct {
	items []chQueryItem
}

func (h *chQueryHeap) reset()            { h.items = h.items[:0] }
func (h *chQueryHeap) len() int          { return len(h.items) }
func (h *chQueryHeap) peek() chQueryItem { return h.items[0] }
func (h *chQueryHeap) push(item chQueryItem) {
	h.items = append(h.items, item)
	index := len(h.items) - 1
	for index > 0 {
		parent := (index - 1) / 2
		if !queryLess(h.items[index], h.items[parent]) {
			break
		}
		h.items[index], h.items[parent] = h.items[parent], h.items[index]
		index = parent
	}
}
func (h *chQueryHeap) pop() chQueryItem {
	last := len(h.items) - 1
	result := h.items[0]
	h.items[0] = h.items[last]
	h.items = h.items[:last]
	if last > 0 {
		index := 0
		for {
			left, right := 2*index+1, 2*index+2
			smallest := index
			if left < last && queryLess(h.items[left], h.items[smallest]) {
				smallest = left
			}
			if right < last && queryLess(h.items[right], h.items[smallest]) {
				smallest = right
			}
			if smallest == index {
				break
			}
			h.items[index], h.items[smallest] = h.items[smallest], h.items[index]
			index = smallest
		}
	}
	return result
}
func queryLess(a, b chQueryItem) bool {
	if a.distance != b.distance {
		return a.distance < b.distance
	}
	return a.node < b.node
}
