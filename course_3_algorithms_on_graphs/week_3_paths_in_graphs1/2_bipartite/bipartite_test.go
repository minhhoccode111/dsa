package main

import "testing"

func TestIsBipartite(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		want bool
	}{
		{"empty", [][]int{}, true},
		{"single vertex", [][]int{{}}, true},
		{"two isolated", [][]int{{}, {}}, true},
		{"one edge", [][]int{{1}, {0}}, true},
		{"path of four", [][]int{{1}, {0, 2}, {1, 3}, {2}}, true},
		{"even cycle", [][]int{{1, 3}, {0, 2}, {1, 3}, {0, 2}}, true},
		{"odd cycle", [][]int{{1, 2}, {0, 2}, {0, 1}}, false},
		{"self-loop", [][]int{{0}}, false},
		{"disconnected odd cycle", [][]int{{1}, {0}, {3, 4}, {2, 4}, {2, 3}}, false},
		{"fixture 1", [][]int{{1, 3, 2}, {0, 2}, {1, 0}, {0}}, false},
		{"fixture 2", [][]int{{3}, {3, 4}, {3}, {0, 1, 2}, {1}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBipartite(tt.adj); got != tt.want {
				t.Fatalf("isBipartite(%v) = %v, want %v", tt.adj, got, tt.want)
			}
		})
	}
}
