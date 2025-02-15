package main

import "testing"

func TestShortestDistance(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		s, t int
		want int
	}{
		{"same vertex", [][]int{{}}, 0, 0, 0},
		{"direct edge", [][]int{{1}, {0}}, 0, 1, 1},
		{"unreachable", [][]int{{}, {}}, 0, 1, -1},
		{"path", [][]int{{1}, {0, 2}, {1, 3}, {2}}, 0, 3, 3},
		{"shortcut", [][]int{{1, 2}, {0, 2}, {0, 1}}, 0, 2, 1},
		{"fixture 1", [][]int{{1, 3}, {0, 2, 3}, {1, 3}, {0, 1, 2}}, 1, 3, 1},
		{"fixture 2", [][]int{{2, 3}, {4}, {0, 3}, {0, 2}, {1}}, 2, 4, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortestDistance(tt.adj, tt.s, tt.t); got != tt.want {
				t.Fatalf("shortestDistance(%v, %d, %d) = %d, want %d", tt.adj, tt.s, tt.t, got, tt.want)
			}
		})
	}
}
