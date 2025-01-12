package main

import "testing"

func TestGCD(t *testing.T) {
	tests := []struct {
		name string
		a, b int64
		want int64
	}{
		{"course sample", 28851538, 1183019, 17657},
		{"coprime", 18, 35, 1},
		{"common factor", 48, 18, 6},
		{"zero first", 0, 5, 5},
		{"zero second", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gcd(tt.a, tt.b); got != tt.want {
				t.Fatalf("gcd(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
