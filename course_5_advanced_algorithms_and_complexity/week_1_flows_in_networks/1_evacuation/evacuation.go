// go run evacuation.go

// Find the maximum flow in a network using Dinic's algorithm.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type flowEdge struct {
	to   int
	rev  int
	cap  int64
	flow int64
}

type flowGraph struct {
	adj [][]flowEdge
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{adj: make([][]flowEdge, n)}
}

func (g *flowGraph) addEdge(from, to int, capacity int64) {
	forward := flowEdge{to: to, rev: len(g.adj[to]), cap: capacity}
	backward := flowEdge{to: from, rev: len(g.adj[from]), cap: 0}
	g.adj[from] = append(g.adj[from], forward)
	g.adj[to] = append(g.adj[to], backward)
}

func (g *flowGraph) size() int {
	return len(g.adj)
}

func (g *flowGraph) maxFlow(source, sink int) int64 {
	var total int64
	level := make([]int, len(g.adj))
	for {
		for i := range level {
			level[i] = -1
		}
		level[source] = 0
		queue := []int{source}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, edge := range g.adj[current] {
				if edge.cap-edge.flow > 0 && level[edge.to] < 0 {
					level[edge.to] = level[current] + 1
					queue = append(queue, edge.to)
				}
			}
		}
		if level[sink] < 0 {
			return total
		}
		next := make([]int, len(g.adj))
		for {
			pushed := g.push(source, sink, 1<<60, level, next)
			if pushed == 0 {
				break
			}
			total += pushed
		}
	}
}

func (g *flowGraph) push(node, sink int, limit int64, level, next []int) int64 {
	if node == sink {
		return limit
	}
	for ; next[node] < len(g.adj[node]); next[node]++ {
		edge := &g.adj[node][next[node]]
		if edge.cap-edge.flow <= 0 || level[edge.to] != level[node]+1 {
			continue
		}
		available := edge.cap - edge.flow
		if available > limit {
			available = limit
		}
		pushed := g.push(edge.to, sink, available, level, next)
		if pushed > 0 {
			edge.flow += pushed
			reverse := &g.adj[edge.to][edge.rev]
			reverse.flow -= pushed
			return pushed
		}
	}
	return 0
}

func parseFlowGraph(input string) (*flowGraph, error) {
	reader := bufio.NewReaderSize(strings.NewReader(input), len(input)+1)
	var vertexCount, edgeCount int
	if _, err := fmt.Fscan(reader, &vertexCount, &edgeCount); err != nil {
		return nil, err
	}
	graph := newFlowGraph(vertexCount)
	for i := 0; i < edgeCount; i++ {
		var u, v int
		var capacity int64
		if _, err := fmt.Fscan(reader, &u, &v, &capacity); err != nil {
			return nil, err
		}
		graph.addEdge(u-1, v-1, capacity)
	}
	return graph, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var vertexCount, edgeCount int
	if _, err := fmt.Fscan(reader, &vertexCount, &edgeCount); err != nil {
		return
	}
	graph := newFlowGraph(vertexCount)
	for i := 0; i < edgeCount; i++ {
		var u, v int
		var capacity int64
		if _, err := fmt.Fscan(reader, &u, &v, &capacity); err != nil {
			return
		}
		graph.addEdge(u-1, v-1, capacity)
	}
	fmt.Println(graph.maxFlow(0, graph.size()-1))
}
