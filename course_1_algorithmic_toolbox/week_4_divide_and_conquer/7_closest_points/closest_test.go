package main

import (
	"math"
	"testing"
)

func TestClosestPair(t *testing.T) {
	tests := []struct {
		name   string
		points []Point
		want   float64
	}{
		{"course sample 1", []Point{{0, 0}, {3, 4}}, 5},
		{
			"course sample 2",
			[]Point{{4, 4}, {-2, -2}, {-3, -4}, {-1, 3}, {2, 3}, {-4, 0}, {1, 1}, {-1, -1}, {3, -1}, {-4, 2}, {-2, 4}},
			math.Sqrt2,
		},
		{"duplicate points", []Point{{1, 1}, {1, 1}, {5, 5}}, 0},
		{"three points", []Point{{0, 0}, {10, 0}, {0, 1}}, 1},
		{"vertical pair", []Point{{0, 0}, {0, 3}, {100, 100}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := closestPair(tt.points)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("closestPair(%v) = %v, want %v", tt.points, got, tt.want)
			}
		})
	}
}
