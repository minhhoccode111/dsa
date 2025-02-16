package main

import "testing"

func TestDijkstraDistance(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		cost [][]int64
		s, t int
		want int64
	}{
		{"same vertex", [][]int{{}}, [][]int64{{}}, 0, 0, 0},
		{"direct edge", [][]int{{1}, {}}, [][]int64{{7}, {}}, 0, 1, 7},
		{"unreachable", [][]int{{}, {}}, [][]int64{{}, {}}, 0, 1, -1},
		{"prefers shorter path", [][]int{{1, 2}, {2}, {}}, [][]int64{{5, 10}, {1}, {}}, 0, 2, 6},
		{"fixture 1", [][]int{{1, 2}, {2}, {}, {0}}, [][]int64{{1, 5}, {2}, {}, {2}}, 0, 2, 3},
		{"fixture 2", [][]int{{1, 2}, {2, 3, 4}, {1, 4, 3}, {}, {3}}, [][]int64{{4, 2}, {2, 2, 3}, {1, 4, 4}, {}, {1}}, 0, 4, 6},
		{"fixture 3", [][]int{{1, 2}, {2}, {}}, [][]int64{{7, 5}, {2}, {}}, 2, 1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distance(tt.adj, tt.cost, tt.s, tt.t); got != tt.want {
				t.Fatalf("distance(%v, %v, %d, %d) = %d, want %d", tt.adj, tt.cost, tt.s, tt.t, got, tt.want)
			}
		})
	}
}
