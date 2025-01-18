package main

import "testing"

func TestCountInversions(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		want int64
	}{
		{"course sample", []int{2, 3, 9, 2, 9}, 2},
		{"sorted", []int{1, 2, 3, 4}, 0},
		{"reverse", []int{4, 3, 2, 1}, 6},
		{"duplicates", []int{2, 2, 1}, 2},
		{"empty", nil, 0},
		{"single", []int{5}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countInversions(tt.a); got != tt.want {
				t.Fatalf("countInversions(%v) = %d, want %d", tt.a, got, tt.want)
			}
		})
	}
}
