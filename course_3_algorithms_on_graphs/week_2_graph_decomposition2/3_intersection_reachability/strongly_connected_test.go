package main

import "testing"

func TestCountSCCs(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		want int
	}{
		{"empty", [][]int{}, 0},
		{"single vertex", [][]int{{}}, 1},
		{"two isolated", [][]int{{}, {}}, 2},
		{"one edge", [][]int{{1}, {}}, 2},
		{"two-cycle", [][]int{{1}, {0}}, 1},
		{"three-cycle", [][]int{{1}, {2}, {0}}, 1},
		{"fixture 1", [][]int{{1}, {2}, {0}, {0}}, 2},
		{"fixture 2", [][]int{{}, {0}, {1, 0}, {2, 0}, {1, 2}}, 5},
		{"chain plus cycle", [][]int{{1}, {2}, {1}, {}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countSCCs(tt.adj); got != tt.want {
				t.Fatalf("countSCCs(%v) = %d, want %d", tt.adj, got, tt.want)
			}
		})
	}
}
