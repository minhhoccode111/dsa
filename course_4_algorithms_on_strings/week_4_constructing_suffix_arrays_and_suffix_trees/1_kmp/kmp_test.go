package main

import (
	"reflect"
	"testing"
)

func TestComputePrefixFunction(t *testing.T) {
	tests := []struct {
		pattern string
		want    []int
	}{
		{"", nil},
		{"A", []int{0}},
		{"AAAA", []int{0, 1, 2, 3}},
		{"ABABACA", []int{0, 0, 1, 2, 3, 0, 1}},
		{"ABCDE", []int{0, 0, 0, 0, 0}},
	}

	for _, tt := range tests {
		if got := computePrefixFunction(tt.pattern); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("computePrefixFunction(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}

func TestFindPattern(t *testing.T) {
	tests := []struct {
		pattern, text string
		want          []int
	}{
		{"GT", "TACG", nil},
		{"ATA", "ATATA", []int{0, 2}},
		{"ATAT", "GATATATGCATATACTT", []int{1, 3, 9}},
		{"A", "AAAA", []int{0, 1, 2, 3}},
		{"ACG", "AC", nil},
	}

	for _, tt := range tests {
		got := findPattern(tt.pattern, tt.text)
		if len(got) == 0 && len(tt.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("findPattern(%q, %q) = %v, want %v", tt.pattern, tt.text, got, tt.want)
		}
	}
}
