package main

import "testing"

func TestSumOfTwoDigits(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"course sample", 2, 3, 5},
		{"zeros", 0, 0, 0},
		{"largest digits", 9, 9, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumOfTwoDigits(tt.a, tt.b); got != tt.want {
				t.Fatalf("sumOfTwoDigits(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
