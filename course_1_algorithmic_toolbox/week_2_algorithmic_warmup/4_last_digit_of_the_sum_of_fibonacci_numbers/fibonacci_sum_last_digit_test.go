package main

import "testing"

func TestFibonacciSumLastDigit(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want int64
	}{
		{"course sample 3", 3, 4},
		{"course sample 100", 100, 5},
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 2},
		{"n=60 resets cycle", 60, 0},
		{"large n fits int64", 1000000000000000000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacciSumLastDigit(tt.n); got != tt.want {
				t.Fatalf("fibonacciSumLastDigit(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
