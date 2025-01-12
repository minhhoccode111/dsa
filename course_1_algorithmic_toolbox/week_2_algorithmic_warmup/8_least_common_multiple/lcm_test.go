package main

import "testing"

func TestLCM(t *testing.T) {
	tests := []struct {
		name string
		a, b int64
		want int64
	}{
		{"course sample", 761457, 614573, 467970912861},
		{"common multiple", 6, 8, 24},
		{"zero first", 0, 5, 0},
		{"zero second", 5, 0, 0},
		{"one", 1, 1, 1},
		{"avoids product overflow", 2000000000, 2000000000, 2000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcm(tt.a, tt.b); got != tt.want {
				t.Fatalf("lcm(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
