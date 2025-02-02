package main

import "testing"

func TestIsBinarySearchTree(t *testing.T) {
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
			"left child too large",
			[]treeNode{{2, 1, -1}, {3, -1, -1}},
			false,
		},
		{
			"right child too small",
			[]treeNode{{2, -1, 1}, {1, -1, -1}},
			false,
		},
		{
			"duplicate is invalid",
			[]treeNode{{2, 1, -1}, {2, -1, -1}},
			false,
		},
		{
			"deep bound violation",
			[]treeNode{
				{10, 1, 2},
				{5, -1, 3},
				{20, -1, -1},
				{12, -1, -1},
			},
			false,
		},
		{"empty", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBinarySearchTree(tt.tree); got != tt.want {
				t.Fatalf("isBinarySearchTree(%v) = %v, want %v", tt.tree, got, tt.want)
			}
		})
	}
}
