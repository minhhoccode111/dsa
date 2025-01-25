package main

import "testing"

func TestComputeHeight(t *testing.T) {
	tests := []struct {
		name    string
		parents []int
		want    int
	}{
		{"course sample", []int{4, -1, 4, 1, 1}, 3},
		{"single node", []int{-1}, 1},
		{"chain", []int{-1, 0, 1, 2}, 4},
		{"star", []int{-1, 0, 0, 0}, 2},
		{"two children", []int{-1, 0, 0}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeHeight(tt.parents); got != tt.want {
				t.Fatalf("computeHeight(%v) = %d, want %d", tt.parents, got, tt.want)
			}
		})
	}
}
