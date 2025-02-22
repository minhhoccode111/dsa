package main

import (
	"math"
	"testing"
)

func TestClustering(t *testing.T) {
	tests := []struct {
		name string
		x, y []int
		k    int
		want float64
	}{
		{"two points", []int{0, 3}, []int{0, 4}, 2, 5},
		{"three points", []int{0, 1, 2}, []int{0, 0, 0}, 2, 1},
		{"fixture 1", []int{7, 4, 5, 1, 2, 5, 3, 7, 2, 4, 6, 2}, []int{6, 3, 1, 7, 7, 7, 3, 8, 8, 4, 7, 6}, 3, 2.828427124746},
		{"fixture 2", []int{3, 1, 4, 9, 9, 8, 3, 4}, []int{1, 2, 6, 8, 9, 9, 11, 12}, 4, 5},
		{"all separate clusters", []int{0, 10, 20}, []int{0, 0, 0}, 3, 10},
		{"duplicate points", []int{1, 1, 5}, []int{1, 1, 5}, 2, math.Sqrt(32)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clustering(tt.x, tt.y, tt.k)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("clustering(%v, %v, %d) = %.12f, want %.12f", tt.x, tt.y, tt.k, got, tt.want)
			}
		})
	}
}
