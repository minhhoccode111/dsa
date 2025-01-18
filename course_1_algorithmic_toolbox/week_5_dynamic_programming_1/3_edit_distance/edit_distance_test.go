package main

import "testing"

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want int
	}{
		{"course sample short ports", "short", "ports", 3},
		{"course sample editing distance", "editing", "distance", 5},
		{"equal", "ab", "ab", 0},
		{"empty first", "", "abc", 3},
		{"empty second", "abc", "", 3},
		{"both empty", "", "", 0},
		{"substitution", "cat", "cut", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := editDistance(tt.s, tt.t); got != tt.want {
				t.Fatalf("editDistance(%q, %q) = %d, want %d", tt.s, tt.t, got, tt.want)
			}
		})
	}
}
