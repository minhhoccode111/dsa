package main

import "testing"

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int64
	}{
		{"n=0", 0, 0},
		{"n=1", 1, 1},
		{"n=2", 2, 1},
		{"n=3", 3, 2},
		{"n=5", 5, 5},
		{"n=10", 10, 55},
		{"n=15", 15, 610},
		{"n=20", 20, 6765},
		{"n=45", 45, 1134903170},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacci(tt.n); got != tt.want {
				t.Fatalf("fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
