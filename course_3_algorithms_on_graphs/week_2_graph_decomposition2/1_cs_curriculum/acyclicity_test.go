package main

import "testing"

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		want bool
	}{
		{"empty", [][]int{}, false},
		{"single vertex", [][]int{{}}, false},
		{"direct acyclic", [][]int{{1}, {}}, false},
		{"two-cycle", [][]int{{1}, {0}}, true},
		{"self-loop", [][]int{{0}}, true},
		{"fixture 1", [][]int{{1}, {2}, {0, 3}, {}}, true},
		{"fixture 2", [][]int{{1, 2, 3}, {2, 4}, {3, 4}, {4}, {}}, false},
		{"disconnected cycle", [][]int{{1}, {2}, {1}, {}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasCycle(tt.adj); got != tt.want {
				t.Fatalf("hasCycle(%v) = %v, want %v", tt.adj, got, tt.want)
			}
		})
	}
}
