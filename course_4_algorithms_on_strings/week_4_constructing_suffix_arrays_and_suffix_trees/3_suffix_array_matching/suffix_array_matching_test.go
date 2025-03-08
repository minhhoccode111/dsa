package main

import (
	"reflect"
	"testing"
)

func TestFindOccurrences(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		patterns []string
		want     []int
	}{
		{"fixture 1", "AAA", []string{"A"}, []int{0, 1, 2}},
		{"fixture 2", "ATA", []string{"T", "G", "C"}, []int{1}},
		{"fixture 3", "ATATATA", []string{"ATA", "C", "TATAT"}, []int{0, 1, 2, 4}},
		{"overlapping patterns", "AAAA", []string{"AA", "AAA"}, []int{0, 1, 2}},
		{"no patterns", "ACGT", nil, nil},
		{"no matches", "ACGT", []string{"TT", "GG"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOccurrences(tt.text, tt.patterns); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("findOccurrences(%q, %v) = %v, want %v", tt.text, tt.patterns, got, tt.want)
			}
		})
	}
}
