package main

import "testing"

func TestFibonacciLastDigit(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want int64
	}{
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 1},
		{"n=3", 3, 2},
		{"n=10", 10, 5},
		{"n=60 resets cycle", 60, 0},
		{"n=61", 61, 1},
		{"course sample 139", 139, 1},
		{"course sample 91239", 91239, 6},
		{"course sample 239", 239, 1},
		{"large n fits int64", 1000000000000000000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacciLastDigit(tt.n); got != tt.want {
				t.Fatalf("fibonacciLastDigit(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
