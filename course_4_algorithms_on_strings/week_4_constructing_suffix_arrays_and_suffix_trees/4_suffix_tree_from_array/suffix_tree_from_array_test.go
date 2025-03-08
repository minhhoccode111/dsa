package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSuffixArrayToSuffixTree(t *testing.T) {
	tests := []struct {
		name             string
		text             string
		sa, lcp          []int
		wantLabelLines   []string
		wantExactFixture bool
		wantLines        []string
	}{
		{
			name:             "one suffix",
			text:             "A$",
			sa:               []int{1, 0},
			lcp:              []int{0},
			wantExactFixture: true,
			wantLines:        []string{"1 2", "0 2"},
			wantLabelLines:   []string{"$", "A$"},
		},
		{
			name:             "repeated characters",
			text:             "AAA$",
			sa:               []int{3, 2, 1, 0},
			lcp:              []int{0, 1, 2},
			wantExactFixture: true,
			wantLines:        []string{"3 4", "0 1", "3 4", "1 2", "3 4", "2 4"},
			wantLabelLines:   []string{"$", "A", "$", "A", "$", "A$"},
		},
		{
			name:           "mixed string",
			text:           "GTAGT$",
			sa:             []int{5, 2, 3, 0, 4, 1},
			lcp:            []int{0, 0, 2, 0, 1},
			wantLabelLines: []string{"$", "AGT$", "GT", "$", "AGT$", "T", "$", "AGT$"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := suffixArrayToSuffixTree(tt.sa, tt.lcp, tt.text)
			got := outputTree(tree)
			gotLabels := make([]string, len(got))
			for i, line := range got {
				var start, end int
				if _, err := fmt.Sscanf(line, "%d %d", &start, &end); err != nil {
					t.Fatalf("bad edge line %q: %v", line, err)
				}
				if start < 0 || end > len(tt.text) || start >= end {
					t.Fatalf("edge line %q out of bounds for %q", line, tt.text)
				}
				gotLabels[i] = tt.text[start:end]
			}
			if !reflect.DeepEqual(gotLabels, tt.wantLabelLines) {
				t.Fatalf("edge labels = %v, want %v", gotLabels, tt.wantLabelLines)
			}
			if tt.wantExactFixture && !reflect.DeepEqual(got, tt.wantLines) {
				t.Fatalf("outputTree = %v, want %v", got, tt.wantLines)
			}
		})
	}
}

func TestSuffixTreeEdgesSortedByFirstCharacter(t *testing.T) {
	tree := suffixArrayToSuffixTree([]int{5, 2, 3, 0, 4, 1}, []int{0, 0, 2, 0, 1}, "GTAGT$")
	for node, edges := range tree {
		for i := 1; i < len(edges); i++ {
			previous := "GTAGT$"[edges[i-1].start]
			current := "GTAGT$"[edges[i].start]
			if previous > current {
				t.Fatalf("node %d edges not sorted: %q before %q", node, previous, current)
			}
		}
	}
}
