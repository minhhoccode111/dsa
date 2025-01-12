package main

import "testing"

func TestFibonacciSumSquaresLastDigit(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want int64
	}{
		{"course sample 7", 7, 3},
		{"course sample 73", 73, 1},
		{"course sample 1234567890", 1234567890, 0},
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 2},
		{"n=3", 3, 6},
		{"large n fits int64", 1000000000000000000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacciSumSquaresLastDigit(tt.n); got != tt.want {
				t.Fatalf("fibonacciSumSquaresLastDigit(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
