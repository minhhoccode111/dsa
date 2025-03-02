package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestBuildSuffixArray(t *testing.T) {
	tests := []struct {
		text string
		want []int
	}{
		{"GAC$", []int{3, 1, 2, 0}},
		{"GAGAGAGA$", []int{8, 7, 5, 3, 1, 6, 4, 2, 0}},
		{"AACGATAGCGGTAGA$", []int{15, 14, 0, 1, 12, 6, 4, 2, 8, 13, 3, 7, 9, 10, 11, 5}},
		{"$", []int{0}},
	}

	for _, tt := range tests {
		if got := buildSuffixArray(tt.text); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("buildSuffixArray(%q) = %v, want %v", tt.text, got, tt.want)
		}
	}
}

func TestBuildSuffixArrayMatchesBruteForce(t *testing.T) {
	texts := []string{"banana$", "mississippi$", "aaaaaa$", "ACGTACGT$", "zxyxzyx$"}
	for _, text := range texts {
		want := bruteSuffixArray(text)
		if got := buildSuffixArray(text); !reflect.DeepEqual(got, want) {
			t.Fatalf("buildSuffixArray(%q) = %v, want %v", text, got, want)
		}
	}
}

func bruteSuffixArray(text string) []int {
	sa := make([]int, len(text))
	for i := range sa {
		sa[i] = i
	}
	sort.Slice(sa, func(i, j int) bool {
		return text[sa[i]:] < text[sa[j]:]
	})
	return sa
}
