package main

import (
	"math/rand"
	"testing"
)

func TestDistanceQueries(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		edges   []distanceEdge
		queries []distanceQuery
		want    []int64
	}{
		{
			name:    "fixture 1",
			n:       2,
			edges:   []distanceEdge{{0, 1, 1}},
			queries: []distanceQuery{{0, 0}, {1, 1}, {0, 1}, {1, 0}},
			want:    []int64{0, 0, 1, -1},
		},
		{
			name: "fixture 2",
			n:    4,
			edges: []distanceEdge{
				{0, 1, 1},
				{3, 0, 2},
				{1, 2, 2},
				{0, 2, 5},
			},
			queries: []distanceQuery{{0, 2}},
			want:    []int64{3},
		},
		{
			name:    "no edges",
			n:       3,
			queries: []distanceQuery{{0, 0}, {0, 2}},
			want:    []int64{0, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := distanceQueries(tt.n, tt.edges, tt.queries)
			if len(got) != len(tt.want) {
				t.Fatalf("distanceQueries(...) returned %d answers, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("answer %d = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDistanceQueriesMatchesDijkstra(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	for round := 0; round < 200; round++ {
		n := 1 + rng.Intn(25)
		m := rng.Intn(n*n + 1)
		edges := make([]distanceEdge, 0, m)
		for i := 0; i < m; i++ {
			edges = append(edges, distanceEdge{
				from:   rng.Intn(n),
				to:     rng.Intn(n),
				weight: int64(1 + rng.Intn(50)),
			})
		}
		queries := make([]distanceQuery, 0, 50)
		for i := 0; i < 50; i++ {
			queries = append(queries, distanceQuery{from: rng.Intn(n), to: rng.Intn(n)})
		}

		got := distanceQueries(n, edges, queries)
		adj, reverse := buildSearchGraph(n, edges)
		for i, query := range queries {
			want := Distance(adj, reverse, query.from, query.to)
			if got[i] != want {
				t.Fatalf("round %d: distance(%d,%d) = %d, want %d",
					round, query.from, query.to, got[i], want)
			}
		}
	}
}
