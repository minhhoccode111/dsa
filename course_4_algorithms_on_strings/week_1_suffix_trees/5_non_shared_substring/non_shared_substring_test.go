package main

import (
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name        string
		p, q        string
		want        string
		wantLength  int
		exactAnswer bool
	}{
		{"fixture 1", "A", "T", "A", 1, true},
		{"fixture 2", "AAAAAAAAAAAAAAAAAAAA", "TTTTTTTTTTTTTTTTTTTT", "A", 1, true},
		{"fixture 3", "CCAAGCTGCTAGAGG", "CATGCTGGGCTGGCT", "AA", 2, true},
		{"fixture 4", "ATGCGATGACCTGACTGA", "CTCAACGTATTGGCCAGA", "", 3, false},
		{"single repeated character", "AAAA", "AAAA", "", 0, true},
		{"empty second string", "ACGT", "", "A", 1, true},
		{"all substrings shared", "AC", "ACGT", "", 0, true},
		{"lexicographic tie", "AB", "A", "B", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solve(tt.p, tt.q)
			if tt.exactAnswer {
				if got != tt.want {
					t.Fatalf("solve(%q, %q) = %q, want %q", tt.p, tt.q, got, tt.want)
				}
				return
			}
			if len(got) != tt.wantLength {
				t.Fatalf("solve(%q, %q) = %q (length %d), want length %d", tt.p, tt.q, got, len(got), tt.wantLength)
			}
			if got != "" && (!strings.Contains(tt.p, got) || strings.Contains(tt.q, got)) {
				t.Fatalf("solve(%q, %q) = %q is not a valid non-shared substring", tt.p, tt.q, got)
			}
		})
	}
}
