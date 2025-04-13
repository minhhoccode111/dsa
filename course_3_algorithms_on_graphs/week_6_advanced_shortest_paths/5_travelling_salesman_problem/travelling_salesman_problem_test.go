package main

import "testing"

func buildTestGraph(n int, edges []tourEdge) ([][]Edge, [][]Edge) {
	return buildTourGraph(n, edges)
}

func TestMinTour(t *testing.T) {
	edges := []tourEdge{
		{0, 1, 1},
		{1, 2, 1},
		{2, 3, 1},
		{3, 0, 1},
		{1, 0, 1},
	}
	adj, reverse := buildTestGraph(4, edges)

	tests := []struct {
		name   string
		stores []int
		want   int64
	}{
		{"single store", []int{0}, 0},
		{"two stores", []int{0, 1}, 2},
		{"three stores", []int{0, 2}, 4},
		{"all stores", []int{0, 1, 2, 3}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minTour(adj, reverse, tt.stores); got != tt.want {
				t.Fatalf("minTour(%v) = %d, want %d", tt.stores, got, tt.want)
			}
		})
	}
}

func TestMinTourUnreachable(t *testing.T) {
	adj, reverse := buildTestGraph(3, []tourEdge{{0, 1, 5}})
	if got := minTour(adj, reverse, []int{0, 1, 2}); got != -1 {
		t.Fatalf("minTour with unreachable store = %d, want -1", got)
	}
}

func TestTourSolver(t *testing.T) {
	edges := []tourEdge{
		{0, 1, 1},
		{1, 2, 1},
		{2, 3, 1},
		{3, 0, 1},
		{1, 0, 1},
	}
	solver := newTourSolver(4, edges)
	tests := []struct {
		stores []int
		want   int64
	}{
		{[]int{0}, 0},
		{[]int{0, 1}, 2},
		{[]int{0, 2}, 4},
		{[]int{0, 1, 2, 3}, 4},
	}
	for _, tt := range tests {
		if got := solver.minTour(tt.stores); got != tt.want {
			t.Fatalf("solver.minTour(%v) = %d, want %d", tt.stores, got, tt.want)
		}
	}
	if got := newTourSolver(3, []tourEdge{{0, 1, 5}}).minTour([]int{0, 1, 2}); got != -1 {
		t.Fatalf("solver.minTour with unreachable store = %d, want -1", got)
	}
}

func TestMinTourTwentyStores(t *testing.T) {
	const n = 20
	edges := make([]tourEdge, n)
	stores := make([]int, n)
	for i := 0; i < n; i++ {
		edges[i] = tourEdge{from: i, to: (i + 1) % n, weight: 1}
		stores[i] = i
	}
	adj, reverse := buildTourGraph(n, edges)
	if got := minTour(adj, reverse, stores); got != n {
		t.Fatalf("dijkstra minTour = %d, want %d", got, n)
	}
	if got := newTourSolver(n, edges).minTour(stores); got != n {
		t.Fatalf("CH minTour = %d, want %d", got, n)
	}
}
