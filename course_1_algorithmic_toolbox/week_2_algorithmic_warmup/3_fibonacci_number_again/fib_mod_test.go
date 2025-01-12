package main

import "testing"

func TestFibonacciMod(t *testing.T) {
	tests := []struct {
		name string
		n, m int64
		want int64
	}{
		{"course sample 1", 1, 239, 1},
		{"course sample 2", 115, 1000, 885},
		{"course sample 3", 2816213588, 239, 151},
		{"n=0", 0, 2, 0},
		{"n=10 m=2", 10, 2, 1},
		{"n=100 m=10", 100, 10, 5},
		{"mod 1 is always 0", 123456789, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacciMod(tt.n, tt.m); got != tt.want {
				t.Fatalf("fibonacciMod(%d, %d) = %d, want %d", tt.n, tt.m, got, tt.want)
			}
		})
	}
}
