// go run tree_height.go

// Compute the height of a tree given parent indices.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func computeHeight(parents []int) int {
	n := len(parents)
	if n == 0 {
		return 0
	}

	children := make([][]int, n)
	root := -1
	for node, parent := range parents {
		if parent == -1 {
			root = node
		} else {
			children[parent] = append(children[parent], node)
		}
	}
	if root == -1 {
		return 0
	}

	type item struct {
		node  int
		depth int
	}

	queue := make([]item, 1, n)
	queue[0] = item{node: root, depth: 1}
	height := 0

	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if current.depth > height {
			height = current.depth
		}
		for _, child := range children[current.node] {
			queue = append(queue, item{node: child, depth: current.depth + 1})
		}
	}

	return height
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	parents := make([]int, n)
	for i := range parents {
		if _, err := fmt.Fscan(reader, &parents[i]); err != nil {
			return
		}
	}

	fmt.Println(computeHeight(parents))
}
