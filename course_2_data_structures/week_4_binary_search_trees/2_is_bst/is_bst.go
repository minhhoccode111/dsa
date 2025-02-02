// go run is_bst.go

// Check whether a binary tree is a strict binary search tree.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type treeNode struct {
	key, left, right int
}

type bstBound struct {
	node            int
	low, high       int
	hasLow, hasHigh bool
}

func isBinarySearchTree(tree []treeNode) bool {
	if len(tree) == 0 {
		return true
	}

	stack := []bstBound{{node: 0}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.node == -1 {
			continue
		}

		node := tree[current.node]
		if current.hasLow && node.key <= current.low {
			return false
		}
		if current.hasHigh && node.key >= current.high {
			return false
		}

		stack = append(stack,
			bstBound{node: node.left, low: current.low, high: node.key, hasLow: current.hasLow, hasHigh: true},
			bstBound{node: node.right, low: node.key, high: current.high, hasLow: true, hasHigh: current.hasHigh},
		)
	}

	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	tree := make([]treeNode, n)
	for i := range tree {
		if _, err := fmt.Fscan(reader, &tree[i].key, &tree[i].left, &tree[i].right); err != nil {
			return
		}
	}

	if isBinarySearchTree(tree) {
		fmt.Println("CORRECT")
	} else {
		fmt.Println("INCORRECT")
	}
}
