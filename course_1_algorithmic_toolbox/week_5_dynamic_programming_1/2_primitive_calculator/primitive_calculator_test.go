package main

import "testing"

func validSequence(sequence []int, n int) bool {
	if len(sequence) == 0 || sequence[0] != 1 || sequence[len(sequence)-1] != n {
		return false
	}
	for i := 1; i < len(sequence); i++ {
		previous, current := sequence[i-1], sequence[i]
		if current != previous+1 && current != previous*2 && current != previous*3 {
			return false
		}
	}
	return true
}

func TestComputeOperations(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		wantOps int
	}{
		{"one", 1, 0},
		{"two", 2, 1},
		{"three", 3, 1},
		{"four", 4, 2},
		{"five", 5, 3},
		{"ten", 10, 3},
		{"course sample", 96234, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence := computeOperations(tt.n)
			if !validSequence(sequence, tt.n) {
				t.Fatalf("computeOperations(%d) = %v, not a valid sequence", tt.n, sequence)
			}
			if got := len(sequence) - 1; got != tt.wantOps {
				t.Fatalf("computeOperations(%d) uses %d operations, want %d: %v", tt.n, got, tt.wantOps, sequence)
			}
		})
	}
}
