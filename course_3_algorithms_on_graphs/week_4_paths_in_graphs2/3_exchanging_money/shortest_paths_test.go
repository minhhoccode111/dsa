package main

import "testing"

func TestShortestPaths(t *testing.T) {
	tests := []struct {
		name      string
		adj       [][]int
		cost      [][]int64
		s         int
		distance  []int64
		reachable []bool
		shortest  []bool
	}{
		{
			name:      "single reachable vertex",
			adj:       [][]int{{}},
			cost:      [][]int64{{}},
			s:         0,
			distance:  []int64{0},
			reachable: []bool{true},
			shortest:  []bool{true},
		},
		{
			name:      "unreachable vertex",
			adj:       [][]int{{}, {}},
			cost:      [][]int64{{}, {}},
			s:         0,
			distance:  []int64{0, inf},
			reachable: []bool{true, false},
			shortest:  []bool{true, true},
		},
		{
			name:      "negative edge without cycle",
			adj:       [][]int{{1}, {2}, {}},
			cost:      [][]int64{{4}, {-2}, {}},
			s:         0,
			distance:  []int64{0, 4, 2},
			reachable: []bool{true, true, true},
			shortest:  []bool{true, true, true},
		},
		{
			name:      "fixture 1",
			adj:       [][]int{{1, 2}, {2}, {4}, {2}, {3}, {0}},
			cost:      [][]int64{{10, 100}, {5}, {7}, {-18}, {10}, {-1}},
			s:         0,
			distance:  []int64{0, 10, 0, 0, 0, inf},
			reachable: []bool{true, true, true, true, true, false},
			shortest:  []bool{true, true, false, false, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance, reachable, shortest := shortestPaths(tt.adj, tt.cost, tt.s)
			for i := range tt.distance {
				if tt.shortest[i] && distance[i] != tt.distance[i] {
					t.Fatalf("distance[%d] = %d, want %d", i, distance[i], tt.distance[i])
				}
				if reachable[i] != tt.reachable[i] {
					t.Fatalf("reachable[%d] = %v, want %v", i, reachable[i], tt.reachable[i])
				}
				if shortest[i] != tt.shortest[i] {
					t.Fatalf("shortest[%d] = %v, want %v", i, shortest[i], tt.shortest[i])
				}
			}
		})
	}
}

func TestShortestPathsNegativeCycle(t *testing.T) {
	// 0 -> 1 (1), 1 -> 2 (-5), 2 -> 1 (1) gives a negative cycle.
	// 2 -> 3 (2) is reachable but its distance is also undefined.
	adj := [][]int{{1}, {2}, {1, 3}, {}}
	cost := [][]int64{{1}, {-5}, {1, 2}, {}}

	distance, reachable, shortest := shortestPaths(adj, cost, 0)
	if !reachable[0] || !reachable[1] || !reachable[2] || !reachable[3] {
		t.Fatalf("all vertices should be reachable, got %v", reachable)
	}
	if shortest[1] || shortest[2] || shortest[3] {
		t.Fatalf("negative-cycle vertices must be marked undefined, got %v", shortest)
	}
	if distance[0] != 0 || !shortest[0] {
		t.Fatalf("source should stay shortest, got distance=%d shortest=%v", distance[0], shortest[0])
	}
}

func TestShortestPathsUnreachableNegativeCycle(t *testing.T) {
	// Vertex 0 cannot reach the negative cycle 1 <-> 2.
	adj := [][]int{{}, {2}, {1}}
	cost := [][]int64{{}, {-1}, {-1}}

	distance, reachable, shortest := shortestPaths(adj, cost, 0)
	if distance[0] != 0 || !reachable[0] || !shortest[0] {
		t.Fatalf("source = (%d, %v, %v), want (0, true, true)", distance[0], reachable[0], shortest[0])
	}
	if reachable[1] || reachable[2] || !shortest[1] || !shortest[2] {
		t.Fatalf("unreachable cycle vertices = reachable %v shortest %v, want false, true", reachable, shortest)
	}
}
