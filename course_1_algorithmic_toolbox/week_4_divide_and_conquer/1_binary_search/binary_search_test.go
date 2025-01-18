package main

import "testing"

func TestBinarySearch(t *testing.T) {
	keys := []int{1, 5, 8, 12, 13}

	tests := []struct {
		name  string
		query int
		want  int
	}{
		{"course sample 8", 8, 2},
		{"first element", 1, 0},
		{"missing above", 23, -1},
		{"missing between", 11, -1},
		{"last element", 13, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearch(keys, tt.query); got != tt.want {
				t.Fatalf("binarySearch(%v, %d) = %d, want %d", keys, tt.query, got, tt.want)
			}
		})
	}

	if got := binarySearch(nil, 1); got != -1 {
		t.Fatalf("binarySearch(nil, 1) = %d, want -1", got)
	}
}
