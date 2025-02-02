// go run rope.go

// Implement a rope with an implicit treap.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type ropeNode struct {
	char     byte
	priority uint64
	left     *ropeNode
	right    *ropeNode
	size     int
}

func nodeSize(node *ropeNode) int {
	if node == nil {
		return 0
	}
	return node.size
}

func updateRopeNode(node *ropeNode) {
	if node == nil {
		return
	}
	node.size = 1 + nodeSize(node.left) + nodeSize(node.right)
}

func splitRope(root *ropeNode, count int) (*ropeNode, *ropeNode) {
	if root == nil {
		return nil, nil
	}

	leftSize := nodeSize(root.left)
	if count <= leftSize {
		left, right := splitRope(root.left, count)
		root.left = right
		updateRopeNode(root)
		return left, root
	}

	left, right := splitRope(root.right, count-leftSize-1)
	root.right = left
	updateRopeNode(root)
	return root, right
}

func mergeRope(left, right *ropeNode) *ropeNode {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if left.priority > right.priority {
		left.right = mergeRope(left.right, right)
		updateRopeNode(left)
		return left
	}

	right.left = mergeRope(left, right.left)
	updateRopeNode(right)
	return right
}

type rope struct {
	root *ropeNode
	rng  *rand.Rand
}

func newRope(s string) *rope {
	r := &rope{rng: rand.New(rand.NewSource(1))}
	for i := 0; i < len(s); i++ {
		r.root = mergeRope(r.root, &ropeNode{
			char:     s[i],
			priority: r.rng.Uint64(),
			size:     1,
		})
	}
	return r
}

func (r *rope) process(i, j, k int) {
	left, rest := splitRope(r.root, i)
	middle, right := splitRope(rest, j-i+1)
	remaining := mergeRope(left, right)
	before, after := splitRope(remaining, k)
	r.root = mergeRope(before, mergeRope(middle, after))
}

func (r *rope) result() string {
	var builder strings.Builder
	builder.Grow(nodeSize(r.root))

	stack := make([]*ropeNode, 0)
	node := r.root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		builder.WriteByte(node.char)
		node = node.right
	}

	return builder.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	rope := newRope(s)

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	for i := 0; i < q; i++ {
		var left, right, insertAt int
		if _, err := fmt.Fscan(reader, &left, &right, &insertAt); err != nil {
			return
		}
		rope.process(left, right, insertAt)
	}

	fmt.Println(rope.result())
}
