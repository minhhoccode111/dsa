package main

import (
	"math"
	"testing"
)

func TestMinimumDistance(t *testing.T) {
	tests := []struct {
		name string
		x, y []int
		want float64
	}{
		{"single point", []int{5}, []int{7}, 0},
		{"two points", []int{0, 3}, []int{0, 4}, 5},
		{"unit square", []int{0, 0, 1, 1}, []int{0, 1, 0, 1}, 3},
		{"fixture 1", []int{0, 0, 1, 1}, []int{0, 1, 0, 1}, 3},
		{"fixture 2", []int{0, 0, 1, 3, 3}, []int{0, 2, 1, 0, 2}, 7.06449510224598},
		{"collinear", []int{0, 1, 2, 3}, []int{0, 0, 0, 0}, 3},
		{"duplicate points", []int{1, 1, 2}, []int{1, 1, 2}, math.Sqrt2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minimumDistance(tt.x, tt.y)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("minimumDistance(%v, %v) = %.12f, want %.12f", tt.x, tt.y, got, tt.want)
			}
		})
	}
}
