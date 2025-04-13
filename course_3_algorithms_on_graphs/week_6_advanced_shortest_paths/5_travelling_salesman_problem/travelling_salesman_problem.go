// go run travelling_salesman_problem.go

// Compute the shortest tour from a depot through a set of stores and back.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/bits"
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
	ch      *contractionHierarchy
	scratch tourScratch
}

func newTourSolver(n int, edges []tourEdge) *tourSolver {
	return &tourSolver{ch: newContractionHierarchy(n, edges)}
}

func (s *tourSolver) minTour(stores []int) int64 {
	return minTourWithScratch(s.ch.distance, stores, &s.scratch)
}

func minTour(adj, reverse [][]Edge, stores []int) int64 {
	return minTourWithDistance(func(from, to int) int64 {
		return Distance(adj, reverse, from, to)
	}, stores)
}

func minTourWithDistance(distance func(from, to int) int64, stores []int) int64 {
	return minTourWithScratch(distance, stores, &tourScratch{})
}

type tourScratch struct {
	dp   []int32
	dist [][]int32
}

func minTourWithScratch(distance func(from, to int) int64, stores []int, scratch *tourScratch) int64 {
	count := len(stores)
	if count <= 1 {
		return 0
	}

	const inf32 int32 = 1 << 30
	if cap(scratch.dist) < count {
		scratch.dist = make([][]int32, count)
	} else {
		scratch.dist = scratch.dist[:count]
	}
	dist := scratch.dist
	for i := 0; i < count; i++ {
		if cap(dist[i]) < count {
			dist[i] = make([]int32, count)
		} else {
			dist[i] = dist[i][:count]
		}
		for j := 0; j < count; j++ {
			if i == j {
				continue
			}
			d := distance(stores[i], stores[j])
			if d < 0 || d >= int64(inf32) {
				dist[i][j] = inf32
				continue
			}
			dist[i][j] = int32(d)
		}
	}

	full := 1 << count
	size := full * count
	if cap(scratch.dp) < size {
		scratch.dp = make([]int32, size)
	} else {
		scratch.dp = scratch.dp[:size]
	}
	dp := scratch.dp
	for i := range dp {
		dp[i] = inf32
	}
	dp[count] = 0 // mask with only the depot, ending there

	allBits := full - 1
	for mask := 1; mask < full; mask++ {
		if mask&1 == 0 {
			continue
		}
		for lastBits := mask; lastBits != 0; lastBits &= lastBits - 1 {
			last := bits.TrailingZeros(uint(lastBits))
			current := dp[mask*count+last]
			if current == inf32 {
				continue
			}
			remaining := allBits ^ mask
			for nextBits := remaining; nextBits != 0; nextBits &= nextBits - 1 {
				next := bits.TrailingZeros(uint(nextBits))
				weight := dist[last][next]
				if weight == inf32 || current >= inf32-1-weight {
					continue
				}
				candidate := current + weight
				nextMask := mask | 1<<next
				index := nextMask*count + next
				if candidate < dp[index] {
					dp[index] = candidate
				}
			}
		}
	}

	best := inf32
	for lastBits := allBits &^ 1; lastBits != 0; lastBits &= lastBits - 1 {
		last := bits.TrailingZeros(uint(lastBits))
		current := dp[allBits*count+last]
		weight := dist[last][0]
		if current == inf32 || weight == inf32 || current >= inf32-1-weight {
			continue
		}
		if candidate := current + weight; candidate < best {
			best = candidate
		}
	}

	if best == inf32 {
		return -1
	}
	return int64(best)
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

	contractedNeighbors []int
	level               []int
	neighborStamp       []int32
	neighborMarkStamp   int32
	rank                int
}

type chShortcut struct {
	from, to int
	weight   int64
}

type chPriorityItem struct {
	priority int
	node     int
}

type chPriorityHeap []chPriorityItem

func (h chPriorityHeap) Len() int { return len(h) }
func (h chPriorityHeap) Less(i, j int) bool {
	if h[i].priority != h[j].priority {
		return h[i].priority < h[j].priority
	}
	return h[i].node < h[j].node
}
func (h chPriorityHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *chPriorityHeap) Push(x interface{}) {
	*h = append(*h, x.(chPriorityItem))
}
func (h *chPriorityHeap) Pop() interface{} {
	old := *h
	item := old[len(old)-1]
	*h = old[:len(old)-1]
	return item
}

func (b *chBuilder) build() {
	n := b.graph.n
	b.contractedNeighbors = make([]int, n)
	b.level = make([]int, n)
	b.neighborStamp = make([]int32, n)

	queue := make(chPriorityHeap, 0, n)
	for node := 0; node < n; node++ {
		queue = append(queue, chPriorityItem{
			priority: len(b.graph.in[node]) + len(b.graph.out[node]),
			node:     node,
		})
	}
	heap.Init(&queue)

	for queue.Len() > 0 {
		item := heap.Pop(&queue).(chPriorityItem)
		node := item.node
		if b.contracted[node] {
			continue
		}

		shortcuts := b.findShortcuts(node)
		inDegree, outDegree := 0, 0
		for _, edge := range b.graph.in[node] {
			if !b.contracted[edge.to] {
				inDegree++
			}
		}
		for _, edge := range b.graph.out[node] {
			if !b.contracted[edge.to] {
				outDegree++
			}
		}
		priority := len(shortcuts) - inDegree - outDegree +
			b.contractedNeighbors[node] + b.level[node]
		if queue.Len() > 0 && priority > queue[0].priority {
			heap.Push(&queue, chPriorityItem{priority: priority, node: node})
			continue
		}

		b.applyShortcuts(shortcuts)
		b.contracted[node] = true
		b.graph.level[node] = b.rank
		b.rank++
		b.markNeighbors(node)
	}
}

func (b *chBuilder) findShortcuts(node int) []chShortcut {
	shortcuts := make([]chShortcut, 0)
	incomingEdges := b.graph.in[node]
	outgoingEdges := b.graph.out[node]

	type target struct {
		node  int
		limit int64
	}
	targets := make([]target, 0, len(outgoingEdges))

	for _, incoming := range incomingEdges {
		u := incoming.to
		if b.contracted[u] || u == node {
			continue
		}
		targets = targets[:0]
		maxLimit := int64(-1)
		for _, outgoing := range outgoingEdges {
			w := outgoing.to
			if b.contracted[w] || w == u || w == node {
				continue
			}
			limit := incoming.weight + outgoing.weight
			if b.graph.hasEdgeAtMost(u, w, limit) {
				continue
			}
			targets = append(targets, target{node: w, limit: limit})
			if limit > maxLimit {
				maxLimit = limit
			}
		}
		if len(targets) == 0 {
			continue
		}

		b.witnessSearch(u, node, maxLimit)
		for _, t := range targets {
			if b.seen[t.node] == b.stamp && b.dist[t.node] <= t.limit {
				continue
			}
			shortcuts = append(shortcuts, chShortcut{from: u, to: t.node, weight: t.limit})
		}
	}

	return shortcuts
}

func (b *chBuilder) applyShortcuts(shortcuts []chShortcut) {
	for _, shortcut := range shortcuts {
		b.graph.addArc(shortcut.from, shortcut.to, shortcut.weight)
	}
}

func (b *chBuilder) markNeighbors(node int) {
	b.neighborMarkStamp++
	if b.neighborMarkStamp == 0 {
		for i := range b.neighborStamp {
			b.neighborStamp[i] = 0
		}
		b.neighborMarkStamp = 1
	}
	stamp := b.neighborMarkStamp
	visit := func(neighbor int) {
		if b.contracted[neighbor] || b.neighborStamp[neighbor] == stamp {
			return
		}
		b.neighborStamp[neighbor] = stamp
		b.contractedNeighbors[neighbor]++
		if b.level[node]+1 > b.level[neighbor] {
			b.level[neighbor] = b.level[node] + 1
		}
	}
	for _, edge := range b.graph.in[node] {
		visit(edge.to)
	}
	for _, edge := range b.graph.out[node] {
		visit(edge.to)
	}
}

// witnessSearch explores paths starting at source while avoiding the node
// being contracted. One search serves every outgoing edge of that node, which
// removes the repeated work of a separate search per shortcut candidate.
func (b *chBuilder) witnessSearch(source, avoid int, limit int64) {
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
			return
		}
		settled++
		if settled > contractSettleLimit {
			return
		}
		if item.hops >= contractHopLimit {
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
