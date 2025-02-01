package main

import (
	"reflect"
	"testing"
)

func TestSolveMismatches(t *testing.T) {
	tests := []struct {
		name    string
		k       int
		text    string
		pattern string
		want    []int
	}{
		{"exact no match", 0, "ababab", "baaa", []int{}},
		{"one mismatch", 1, "ababab", "baaa", []int{1}},
		{"one mismatch no match", 1, "xabcabc", "ccc", []int{}},
		{"two mismatches", 2, "xabcabc", "ccc", []int{1, 2, 3, 4}},
		{"all mismatches allowed", 3, "aaa", "xxx", []int{0}},
		{"every position", 2, "aaaa", "aa", []int{0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solveMismatches(tt.k, tt.text, tt.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("solveMismatches(%d, %q, %q) = %v, want %v", tt.k, tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}
