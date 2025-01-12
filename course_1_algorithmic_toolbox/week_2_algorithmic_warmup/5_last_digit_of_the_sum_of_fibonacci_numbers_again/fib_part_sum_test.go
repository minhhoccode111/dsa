package main

import "testing"

func TestFibonacciPartialSumLastDigit(t *testing.T) {
	tests := []struct {
		name string
		m, n int64
		want int64
	}{
		{"course sample 3 to 7", 3, 7, 1},
		{"course sample 10 to 10", 10, 10, 5},
		{"0 to 0", 0, 0, 0},
		{"0 to 1", 0, 1, 1},
		{"3 to 3", 3, 3, 2},
		{"cycle boundary", 59, 60, 1},
		{"large bounds fit int64", 1000000000000000000, 1000000000000000000, 5},
		{"large range from zero", 0, 1000000000000000000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacciPartialSumLastDigit(tt.m, tt.n); got != tt.want {
				t.Fatalf("fibonacciPartialSumLastDigit(%d, %d) = %d, want %d", tt.m, tt.n, got, tt.want)
			}
		})
	}
}
