package main

import (
	"reflect"
	"testing"
)

func TestFindOccurrences(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		want    []int
	}{
		{"course sample", "aba", "abacaba", []int{0, 4}},
		{"no match", "xyz", "abacaba", []int{}},
		{"pattern longer", "abcdef", "abc", []int{}},
		{"overlapping", "aa", "aaa", []int{0, 1}},
		{"whole string", "abc", "abc", []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOccurrences(tt.pattern, tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("findOccurrences(%q, %q) = %v, want %v", tt.pattern, tt.text, got, tt.want)
			}
		})
	}
}
