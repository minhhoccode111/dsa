package main

import "testing"

func TestPartition3(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"course sample 1", []int{3, 3, 3, 3}, 0},
		{"course sample 2", []int{30}, 0},
		{
			"course sample 3",
			[]int{1, 2, 3, 4, 5, 5, 7, 7, 8, 10, 12, 19, 25},
			1,
		},
		{"three ones", []int{1, 1, 1}, 1},
		{"simple partition", []int{1, 1, 2, 2}, 1},
		{"not divisible by three", []int{1, 2, 4}, 0},
		{"empty", nil, 0},
		{"zeros", []int{0, 0, 0}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := partition3(tt.values); got != tt.want {
				t.Fatalf("partition3(%v) = %d, want %d", tt.values, got, tt.want)
			}
		})
	}
}
