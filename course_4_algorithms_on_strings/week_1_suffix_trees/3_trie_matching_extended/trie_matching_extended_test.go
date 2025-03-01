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
		{"fixture 2", "ACATA", []string{"AT", "A", "AG"}, []int{0, 2, 4}},
		{"prefix of another pattern", "ACATA", []string{"A", "AC"}, []int{0, 2, 4}},
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
