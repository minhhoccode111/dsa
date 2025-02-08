package main

import "testing"

func TestCountComponents(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		want int
	}{
		{"empty graph", [][]int{}, 0},
		{"single vertex", [][]int{{}}, 1},
		{"two isolated", [][]int{{}, {}}, 2},
		{"one edge", [][]int{{1}, {0}}, 1},
		{"fixture 1", [][]int{{1}, {0, 2}, {1}, {}}, 2},
		{"path", [][]int{{1}, {0, 2}, {1, 3}, {2}}, 1},
		{"three components", [][]int{{1}, {0}, {3}, {2}, {}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countComponents(tt.adj); got != tt.want {
				t.Fatalf("countComponents(%v) = %d, want %d", tt.adj, got, tt.want)
			}
		})
	}
}
