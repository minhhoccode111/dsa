package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCountOccurrences(t *testing.T) {
	tests := []struct {
		name     string
		bwt      string
		patterns []string
		want     []int
	}{
		{"fixture 1", "AGGGAA$", []string{"GA"}, []int{3}},
		{"fixture 2", "ATT$AA", []string{"ATA", "A"}, []int{2, 3}},
		{"fixture 3", "AT$TCTATG", []string{"TCT", "TATG"}, []int{0, 0}},
		{"missing character", "AC$A", []string{"G", "AC"}, []int{0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starts, occ := preprocessBWT(tt.bwt)
			got := make([]int, len(tt.patterns))
			for i, pattern := range tt.patterns {
				got[i] = countOccurrences(pattern, tt.bwt, starts, occ)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("counts = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountOccurrencesMatchesSubstringCount(t *testing.T) {
	texts := []string{"ACGTACGT$", "AAAA$", "GATTACA$", "T$"}
	for _, text := range texts {
		encoded := bwt(text)
		starts, occ := preprocessBWT(encoded)
		for length := 1; length <= len(text); length++ {
			for start := 0; start+length <= len(text); start++ {
				pattern := text[start : start+length]
				want := overlappingCount(text, pattern)
				if got := countOccurrences(pattern, encoded, starts, occ); got != want {
					t.Fatalf("text %q pattern %q: count = %d, want %d", text, pattern, got, want)
				}
			}
		}
	}
}

func overlappingCount(text, pattern string) int {
	count := 0
	for start := 0; start+len(pattern) <= len(text); start++ {
		if strings.HasPrefix(text[start:], pattern) {
			count++
		}
	}
	return count
}
