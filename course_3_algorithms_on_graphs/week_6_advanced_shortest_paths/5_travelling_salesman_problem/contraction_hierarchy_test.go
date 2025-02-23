package main

import (
	"math/rand"
	"testing"

	"dsa/course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/internal/search"
)

func TestContractionHierarchyMatchesBidirectionalDijkstra(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for round := 0; round < 500; round++ {
		n := 1 + rng.Intn(12)
		m := rng.Intn(n*n + 1)
		edges := make([]tourEdge, 0, m)
		adj := make([][]search.Edge, n)
		reverse := make([][]search.Edge, n)
		for i := 0; i < m; i++ {
			u := rng.Intn(n)
			v := rng.Intn(n)
			w := int64(1 + rng.Intn(20))
			edges = append(edges, tourEdge{from: u, to: v, weight: w})
			adj[u] = append(adj[u], search.Edge{To: v, Weight: w})
			reverse[v] = append(reverse[v], search.Edge{To: u, Weight: w})
		}
		ch := newContractionHierarchy(n, edges)
		for s := 0; s < n; s++ {
			for target := 0; target < n; target++ {
				want := search.Distance(adj, reverse, s, target)
				got := ch.distance(s, target)
				if want == -1 {
					if got != -1 {
						t.Fatalf("round %d: distance(%d,%d) = %d, want -1; edges=%v", round, s, target, got, edges)
					}
					continue
				}
				if got != want {
					t.Fatalf("round %d: distance(%d,%d) = %d, want %d; edges=%v", round, s, target, got, want, edges)
				}
			}
		}
	}
}

func TestContractionHierarchyDisconnected(t *testing.T) {
	ch := newContractionHierarchy(3, []tourEdge{{from: 0, to: 1, weight: 5}})
	if got := ch.distance(0, 2); got != -1 {
		t.Fatalf("distance(0,2) = %d, want -1", got)
	}
	if got := ch.distance(0, 0); got != 0 {
		t.Fatalf("distance(0,0) = %d, want 0", got)
	}
}

func TestContractionHierarchyGrid(t *testing.T) {
	const side = 8
	n := side * side
	edges := make([]tourEdge, 0, 4*n)
	adj := make([][]search.Edge, n)
	reverse := make([][]search.Edge, n)
	add := func(from, to int, weight int64) {
		edges = append(edges, tourEdge{from: from, to: to, weight: weight})
		adj[from] = append(adj[from], search.Edge{To: to, Weight: weight})
		reverse[to] = append(reverse[to], search.Edge{To: from, Weight: weight})
	}
	for row := 0; row < side; row++ {
		for col := 0; col < side; col++ {
			vertex := row*side + col
			if col+1 < side {
				add(vertex, vertex+1, int64(1+(row*side+col)%17))
				add(vertex+1, vertex, int64(2+(row*side+col)%19))
			}
			if row+1 < side {
				add(vertex, vertex+side, int64(1+(row*side+col)%13))
				add(vertex+side, vertex, int64(2+(row*side+col)%11))
			}
		}
	}

	ch := newContractionHierarchy(n, edges)
	seenLevel := make([]bool, n)
	for _, level := range ch.level {
		if level < 0 || level >= n || seenLevel[level] {
			t.Fatalf("level %d is not a permutation position", level)
		}
		seenLevel[level] = true
	}
	for source := 0; source < n; source++ {
		for target := 0; target < n; target++ {
			want := search.Distance(adj, reverse, source, target)
			got := ch.distance(source, target)
			if got != want {
				t.Fatalf("distance(%d,%d) = %d, want %d", source, target, got, want)
			}
		}
	}
}
