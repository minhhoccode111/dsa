package main

import "testing"

func TestMaximumGold(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		weights  []int
		want     int
	}{
		{"course sample", 10, []int{1, 4, 8}, 9},
		{"zero capacity", 0, []int{1, 4, 8}, 0},
		{"exact fit", 10, []int{10}, 10},
		{"take two", 10, []int{5, 5, 5}, 10},
		{"no items", 10, nil, 0},
		{"all too heavy", 3, []int{4, 5}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maximumGold(tt.capacity, tt.weights); got != tt.want {
				t.Fatalf("maximumGold(%d, %v) = %d, want %d", tt.capacity, tt.weights, got, tt.want)
			}
		})
	}
}
