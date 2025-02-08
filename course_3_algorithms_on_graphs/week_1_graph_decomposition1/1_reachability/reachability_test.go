package main

import "testing"

func TestReach(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
		x, y int
		want bool
	}{
		{"same vertex", [][]int{{}}, 0, 0, true},
		{"isolated target", [][]int{{}, {}}, 0, 1, false},
		{"direct edge", [][]int{{1}, {0}}, 0, 1, true},
		{"path forward", [][]int{{1}, {2}, {}}, 0, 2, true},
		{"path backward", [][]int{{1}, {0, 2}, {1}}, 2, 0, true},
		{"disconnected", [][]int{{1}, {0}, {3}, {2}}, 0, 2, false},
		{"fixture 1", [][]int{{1, 3}, {0, 2}, {1, 3}, {2, 0}}, 0, 3, true},
		{"fixture 2", [][]int{{1}, {0, 2}, {1}, {}}, 0, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reach(tt.adj, tt.x, tt.y); got != tt.want {
				t.Fatalf("reach(%v, %d, %d) = %v, want %v", tt.adj, tt.x, tt.y, got, tt.want)
			}
		})
	}
}
