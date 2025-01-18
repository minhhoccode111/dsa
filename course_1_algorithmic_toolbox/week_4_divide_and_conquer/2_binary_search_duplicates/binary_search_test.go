package main

import "testing"

func TestBinarySearchFirst(t *testing.T) {
	keys := []int{2, 4, 4, 4, 7, 7, 9}

	tests := []struct {
		name  string
		query int
		want  int
	}{
		{"course sample 9", 9, 6},
		{"course sample 4", 4, 1},
		{"course sample 5", 5, -1},
		{"course sample 2", 2, 0},
		{"duplicate seven", 7, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearchFirst(keys, tt.query); got != tt.want {
				t.Fatalf("binarySearchFirst(%v, %d) = %d, want %d", keys, tt.query, got, tt.want)
			}
		})
	}
}
