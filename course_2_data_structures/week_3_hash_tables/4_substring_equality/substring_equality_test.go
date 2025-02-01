package main

import "testing"

func TestSubstringEquality(t *testing.T) {
	solver := newSubstringSolver("trololo")

	tests := []struct {
		name    string
		a, b, l int
		want    bool
	}{
		{"same prefix", 0, 0, 7, true},
		{"course sample 2", 2, 4, 3, true},
		{"course sample 3", 3, 5, 1, true},
		{"course sample 4", 1, 3, 2, false},
		{"empty substrings", 1, 3, 0, true},
		{"single equal", 2, 4, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solver.ask(tt.a, tt.b, tt.l); got != tt.want {
				t.Fatalf("ask(%d, %d, %d) = %v, want %v", tt.a, tt.b, tt.l, got, tt.want)
			}
		})
	}
}
