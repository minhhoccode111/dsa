// go run clustering.go

// Compute the largest minimum distance between k clusters of points.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type clusterEdge struct {
	weight float64
	u, v   int
}

type disjointSet struct {
	parent []int
	rank   []int
}

func newDisjointSet(n int) *disjointSet {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &disjointSet{parent: parent, rank: make([]int, n)}
}

func (d *disjointSet) find(v int) int {
	root := v
	for d.parent[root] != root {
		root = d.parent[root]
	}
	for d.parent[v] != v {
		next := d.parent[v]
		d.parent[v] = root
		v = next
	}
	return root
}

func (d *disjointSet) union(u, v int) bool {
	rootU := d.find(u)
	rootV := d.find(v)
	if rootU == rootV {
		return false
	}

	if d.rank[rootU] > d.rank[rootV] {
		d.parent[rootV] = rootU
	} else {
		d.parent[rootU] = rootV
		if d.rank[rootU] == d.rank[rootV] {
			d.rank[rootV]++
		}
	}
	return true
}

func clustering(x, y []int, k int) float64 {
	n := len(x)
	edges := make([]clusterEdge, 0, n*(n-1)/2)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := float64(x[i] - x[j])
			dy := float64(y[i] - y[j])
			edges = append(edges, clusterEdge{weight: math.Hypot(dx, dy), u: i, v: j})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	set := newDisjointSet(n)
	clusters := n

	for _, edge := range edges {
		if set.find(edge.u) == set.find(edge.v) {
			continue
		}
		if clusters == k {
			return edge.weight
		}
		set.union(edge.u, edge.v)
		clusters--
	}

	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &x[i], &y[i]); err != nil {
			return
		}
	}

	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}

	fmt.Printf("%.9f\n", clustering(x, y, k))
}
