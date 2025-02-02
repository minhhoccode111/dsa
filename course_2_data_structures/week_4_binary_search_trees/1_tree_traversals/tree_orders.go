// go run tree_orders.go

// Print in-order, pre-order, and post-order traversals of a binary tree.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type treeNode struct {
	key, left, right int
}

func inOrder(tree []treeNode) []int {
	result := make([]int, 0, len(tree))
	stack := make([]int, 0, len(tree))
	node := 0

	for node != -1 || len(stack) > 0 {
		for node != -1 {
			stack = append(stack, node)
			node = tree[node].left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, tree[node].key)
		node = tree[node].right
	}

	return result
}

func preOrder(tree []treeNode) []int {
	result := make([]int, 0, len(tree))
	if len(tree) == 0 {
		return result
	}

	stack := []int{0}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, tree[node].key)

		if tree[node].right != -1 {
			stack = append(stack, tree[node].right)
		}
		if tree[node].left != -1 {
			stack = append(stack, tree[node].left)
		}
	}

	return result
}

func postOrder(tree []treeNode) []int {
	result := make([]int, 0, len(tree))
	if len(tree) == 0 {
		return result
	}

	stack := []int{0}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, tree[node].key)

		if tree[node].left != -1 {
			stack = append(stack, tree[node].left)
		}
		if tree[node].right != -1 {
			stack = append(stack, tree[node].right)
		}
	}

	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result
}

func printValues(values []int) {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprint(value)
	}
	fmt.Println(strings.Join(parts, " "))
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

	printValues(inOrder(tree))
	printValues(preOrder(tree))
	printValues(postOrder(tree))
}
