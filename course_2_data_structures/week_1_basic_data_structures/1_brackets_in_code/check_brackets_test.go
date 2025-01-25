package main

import "testing"

func TestFindMismatch(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{"empty", "", 0},
		{"plain brackets", "[]", 0},
		{"two pairs", "{}[]", 0},
		{"nested", "[()]", 0},
		{"nested parens", "(())", 0},
		{"mixed nested", "{[]}()", 0},
		{"unclosed", "{", 1},
		{"wrong nested", "{[}", 3},
		{"code with mismatch", "foo(bar[i);", 10},
		{"closing first", ")", 1},
		{"unclosed first", "(((", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMismatch(tt.text); got != tt.want {
				t.Fatalf("findMismatch(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}
