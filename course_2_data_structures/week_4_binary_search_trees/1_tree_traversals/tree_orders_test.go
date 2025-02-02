package main

import (
	"reflect"
	"testing"
)

func TestTreeOrders(t *testing.T) {
	tree := []treeNode{
		{4, 1, 2},
		{2, 3, 4},
		{5, -1, -1},
		{1, -1, -1},
		{3, -1, -1},
	}

	if got := inOrder(tree); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("inOrder = %v", got)
	}
	if got := preOrder(tree); !reflect.DeepEqual(got, []int{4, 2, 1, 3, 5}) {
		t.Fatalf("preOrder = %v", got)
	}
	if got := postOrder(tree); !reflect.DeepEqual(got, []int{1, 3, 2, 5, 4}) {
		t.Fatalf("postOrder = %v", got)
	}
}

func TestTreeOrdersSingle(t *testing.T) {
	tree := []treeNode{{7, -1, -1}}
	if got := inOrder(tree); !reflect.DeepEqual(got, []int{7}) {
		t.Fatalf("inOrder = %v", got)
	}
	if got := preOrder(tree); !reflect.DeepEqual(got, []int{7}) {
		t.Fatalf("preOrder = %v", got)
	}
	if got := postOrder(tree); !reflect.DeepEqual(got, []int{7}) {
		t.Fatalf("postOrder = %v", got)
	}
}
