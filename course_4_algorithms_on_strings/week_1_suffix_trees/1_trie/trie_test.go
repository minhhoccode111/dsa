package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestBuildTrie(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{"single pattern", []string{"ATA"}, []string{"0->1:A", "1->2:T", "2->3:A"}},
		{"shared prefix", []string{"AT", "AG", "AC"}, []string{"0->1:A", "1->2:T", "1->3:G", "1->4:C"}},
		{"three patterns", []string{"ATAGA", "ATC", "GAT"}, []string{
			"0->1:A", "1->2:T", "2->3:A", "3->4:G", "4->5:A", "2->6:C", "0->7:G", "7->8:A", "8->9:T",
		}},
		{"no patterns", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTrie(buildTrie(tt.patterns))
			sort.Strings(got)
			want := append([]string(nil), tt.want...)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("buildTrie(%v) = %v, want %v", tt.patterns, got, want)
			}
		})
	}
}
