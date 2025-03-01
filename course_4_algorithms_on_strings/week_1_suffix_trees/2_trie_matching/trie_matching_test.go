package main

import (
	"reflect"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		patterns []string
		want     []int
	}{
		{"fixture 1", "AAA", []string{"AA"}, []int{0, 1}},
		{"fixture 2", "AA", []string{"T"}, nil},
		{"fixture 3", "AATCGGGTTCAATCGGGGT", []string{"ATCG", "GGGT"}, []int{1, 4, 11, 15}},
		{"overlapping", "AAAA", []string{"AA"}, []int{0, 1, 2}},
		{"several patterns", "ACGTACGT", []string{"CG", "TA"}, []int{1, 3, 5}},
		{"no patterns", "ACGT", nil, nil},
		{"pattern longer than text", "AC", []string{"ACG"}, nil},
		{"missing pattern", "ACGT", []string{"TT"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve(tt.text, tt.patterns); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("solve(%q, %v) = %v, want %v", tt.text, tt.patterns, got, tt.want)
			}
		})
	}
}
