package main

import "testing"

func TestMaxPairwiseProduct(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int64
		want    int64
	}{
		{"course sample", []int64{1, 2, 3}, 6},
		{"course sample with duplicates", []int64{7, 5, 14, 2, 8, 8, 10, 1, 2, 3}, 140},
		{"two equal numbers", []int64{5, 5}, 25},
		{"zeros", []int64{0, 0}, 0},
		{"maximum product fits int64", []int64{200000, 200000}, 40000000000},
		{"new maximum shifts previous maximum", []int64{10, 1, 2}, 20},
		{"duplicate maximum", []int64{5, 9, 9}, 81},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxPairwiseProduct(tt.numbers); got != tt.want {
				t.Fatalf("maxPairwiseProduct(%v) = %d, want %d", tt.numbers, got, tt.want)
			}
		})
	}
}
