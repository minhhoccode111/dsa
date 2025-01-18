package main

import "testing"

func TestMinCoins(t *testing.T) {
	tests := []struct {
		name  string
		money int
		want  int
	}{
		{"course sample 34", 34, 9},
		{"zero", 0, 0},
		{"one", 1, 1},
		{"two", 2, 2},
		{"three", 3, 1},
		{"four", 4, 1},
		{"six", 6, 2},
		{"seven", 7, 2},
		{"ten", 10, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minCoins(tt.money); got != tt.want {
				t.Fatalf("minCoins(%d) = %d, want %d", tt.money, got, tt.want)
			}
		})
	}
}
