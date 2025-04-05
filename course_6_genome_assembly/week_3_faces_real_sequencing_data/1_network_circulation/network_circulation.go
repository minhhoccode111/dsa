// go run network_circulation.go

// Find a feasible circulation with lower and upper edge bounds.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type circulationEdge struct {
	from, to     int
	lower, upper int
}

type residualEdge struct {
	to   int
	rev  int
	cap  int64
	flow int64
}

type residualNetwork struct {
	adj [][]residualEdge
}

func newResidualNetwork(n int) *residualNetwork {
	return &residualNetwork{adj: make([][]residualEdge, n)}
}

func (network *residualNetwork) addEdge(from, to int, capacity int64) int {
	forwardIndex := len(network.adj[from])
	backwardIndex := len(network.adj[to])
	network.adj[from] = append(network.adj[from], residualEdge{to: to, rev: backwardIndex, cap: capacity})
	network.adj[to] = append(network.adj[to], residualEdge{to: from, rev: forwardIndex, cap: 0})
	return forwardIndex
}

func (network *residualNetwork) maxFlow(source, sink int) int64 {
	var total int64
	level := make([]int, len(network.adj))
	for {
		for i := range level {
			level[i] = -1
		}
		level[source] = 0
		queue := []int{source}
		for len(queue) > 0 {
			vertex := queue[0]
			queue = queue[1:]
			for _, edge := range network.adj[vertex] {
				if edge.cap-edge.flow > 0 && level[edge.to] < 0 {
					level[edge.to] = level[vertex] + 1
					queue = append(queue, edge.to)
				}
			}
		}
		if level[sink] < 0 {
			return total
		}
		next := make([]int, len(network.adj))
		for {
			pushed := network.push(source, sink, 1<<60, level, next)
			if pushed == 0 {
				break
			}
			total += pushed
		}
	}
}

func (network *residualNetwork) push(vertex, sink int, limit int64, level, next []int) int64 {
	if vertex == sink {
		return limit
	}
	for ; next[vertex] < len(network.adj[vertex]); next[vertex]++ {
		edge := &network.adj[vertex][next[vertex]]
		if edge.cap-edge.flow <= 0 || level[edge.to] != level[vertex]+1 {
			continue
		}
		available := edge.cap - edge.flow
		if available > limit {
			available = limit
		}
		pushed := network.push(edge.to, sink, available, level, next)
		if pushed > 0 {
			edge.flow += pushed
			reverse := &network.adj[edge.to][edge.rev]
			reverse.flow -= pushed
			return pushed
		}
	}
	return 0
}

func findCirculation(n, m int, edges []circulationEdge) ([]int64, bool) {
	demand := make([]int64, n)
	for _, edge := range edges {
		demand[edge.from] -= int64(edge.lower)
		demand[edge.to] += int64(edge.lower)
	}
	superSource, superSink := n, n+1
	network := newResidualNetwork(n + 2)
	edgeIndices := make([]int, m)
	for i, edge := range edges {
		edgeIndices[i] = network.addEdge(edge.from, edge.to, int64(edge.upper-edge.lower))
	}
	var required int64
	for vertex := 0; vertex < n; vertex++ {
		if demand[vertex] > 0 {
			network.addEdge(superSource, vertex, demand[vertex])
			required += demand[vertex]
		} else if demand[vertex] < 0 {
			network.addEdge(vertex, superSink, -demand[vertex])
		}
	}
	if flow := network.maxFlow(superSource, superSink); flow != required {
		return nil, false
	}
	flows := make([]int64, m)
	for i, edge := range edges {
		index := edgeIndices[i]
		flows[i] = int64(edge.lower) + network.adj[edge.from][index].flow
	}
	return flows, true
}

func parseCirculation(input string) (int, int, []circulationEdge, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, err
	}
	edges := make([]circulationEdge, m)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i].from, &edges[i].to, &edges[i].lower, &edges[i].upper); err != nil {
			return 0, 0, nil, err
		}
		edges[i].from--
		edges[i].to--
	}
	return n, m, edges, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([]circulationEdge, m)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i].from, &edges[i].to, &edges[i].lower, &edges[i].upper); err != nil {
			return
		}
		edges[i].from--
		edges[i].to--
	}
	flows, ok := findCirculation(n, m, edges)
	if !ok {
		fmt.Println("NO")
		return
	}
	fmt.Println("YES")
	for _, flow := range flows {
		fmt.Println(flow)
	}
}
