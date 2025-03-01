package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestBuildSuffixTree(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{"one suffix", "A$", []string{"A$", "$"}},
		{"repeated prefix", "ACA$", []string{"$", "A", "$", "CA$", "CA$"}},
		{"larger tree", "ATAAATG$", []string{
			"AAATG$", "G$", "T", "ATG$", "TG$", "A", "A", "AAATG$", "G$", "T", "G$", "$",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSuffixTree(tt.text)
			sort.Strings(got)
			want := append([]string(nil), tt.want...)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("buildSuffixTree(%q) = %v, want %v", tt.text, got, want)
			}
		})
	}
}

func TestBuildSuffixTreeEmpty(t *testing.T) {
	if got := buildSuffixTree(""); len(got) != 0 {
		t.Fatalf("buildSuffixTree(\"\") = %v, want no edges", got)
	}
}
