package main

import "testing"

func validTopologicalOrder(adj [][]int, order []int) bool {
	if len(order) != len(adj) {
		return false
	}
	position := make([]int, len(adj))
	seen := make([]bool, len(adj))
	for i, v := range order {
		if v < 0 || v >= len(adj) || seen[v] {
			return false
		}
		seen[v] = true
		position[v] = i
	}
	for u, neighbors := range adj {
		for _, v := range neighbors {
			if position[u] >= position[v] {
				return false
			}
		}
	}
	return true
}

func TestToposort(t *testing.T) {
	tests := []struct {
		name      string
		adj       [][]int
		wantExact []int
	}{
		{"empty", [][]int{}, []int{}},
		{"single vertex", [][]int{{}}, []int{0}},
		{"fixture 1", [][]int{{1}, {}, {0}, {0}}, []int{3, 2, 0, 1}},
		{"fixture 3", [][]int{{}, {0}, {0, 1}, {2}, {0, 1, 2}}, []int{4, 3, 2, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toposort(tt.adj)
			if !validTopologicalOrder(tt.adj, got) {
				t.Fatalf("toposort(%v) = %v, not a valid topological order", tt.adj, got)
			}
			if tt.wantExact == nil {
				return
			}
			if len(got) != len(tt.wantExact) {
				t.Fatalf("toposort(%v) = %v, want %v", tt.adj, got, tt.wantExact)
			}
			for i := range got {
				if got[i] != tt.wantExact[i] {
					t.Fatalf("toposort(%v) = %v, want %v", tt.adj, got, tt.wantExact)
				}
			}
		})
	}
}

func TestToposortFixtureWithAlternateValidOrder(t *testing.T) {
	// Fixture 2's expected order is valid but differs from the Python
	// reference's DFS order. Accept any valid topological ordering.
	adj := [][]int{{}, {}, {}, {0}}
	got := toposort(adj)
	if !validTopologicalOrder(adj, got) {
		t.Fatalf("toposort(%v) = %v, not a valid topological order", adj, got)
	}
}
