package main

import "testing"

func TestHasNegativeCycle(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		cost [][]int64
		want bool
	}{
		{"empty", [][]int{}, [][]int64{}, false},
		{"single vertex", [][]int{{}}, [][]int64{{}}, false},
		{"positive edge", [][]int{{1}, {}}, [][]int64{{5}, {}}, false},
		{"zero cycle", [][]int{{1}, {0}}, [][]int64{{0}, {0}}, false},
		{"negative two-cycle", [][]int{{1}, {0}}, [][]int64{{-1}, {-1}}, true},
		{"fixture 1", [][]int{{1}, {2}, {0}, {0}}, [][]int64{{-5}, {2}, {1}, {2}}, true},
		{"disconnected negative cycle", [][]int{{}, {3}, {1}, {2}}, [][]int64{{}, {0}, {0}, {-1}}, true},
		{"self-loop negative", [][]int{{0}}, [][]int64{{-3}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasNegativeCycle(tt.adj, tt.cost); got != tt.want {
				t.Fatalf("hasNegativeCycle(%v, %v) = %v, want %v", tt.adj, tt.cost, got, tt.want)
			}
		})
	}
}
