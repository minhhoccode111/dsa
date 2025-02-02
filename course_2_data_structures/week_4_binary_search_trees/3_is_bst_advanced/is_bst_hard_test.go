package main

import "testing"

func TestIsBinarySearchTreeHard(t *testing.T) {
	tests := []struct {
		name string
		tree []treeNode
		want bool
	}{
		{
			"course sample",
			[]treeNode{{2, 1, 2}, {1, -1, -1}, {3, -1, -1}},
			true,
		},
		{
			"left duplicate invalid",
			[]treeNode{{2, 1, -1}, {2, -1, -1}},
			false,
		},
		{
			"right duplicate allowed",
			[]treeNode{{2, -1, 1}, {2, -1, -1}},
			true,
		},
		{
			"deep equal violation",
			[]treeNode{
				{10, 1, 2},
				{5, -1, 3},
				{20, -1, -1},
				{12, -1, -1},
			},
			false,
		},
		{
			"right duplicates",
			[]treeNode{
				{2, 1, 2},
				{1, -1, -1},
				{2, -1, -1},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBinarySearchTreeHard(tt.tree); got != tt.want {
				t.Fatalf("isBinarySearchTreeHard(%v) = %v, want %v", tt.tree, got, tt.want)
			}
		})
	}
}
