// go run set_range_sum.go

// Maintain a set of integers and answer range-sum queries with a treap.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const rangeSumModulo int64 = 1000000001

type treapNode struct {
	key      int64
	priority uint64
	left     *treapNode
	right    *treapNode
	sum      int64
}

func nodeSum(node *treapNode) int64 {
	if node == nil {
		return 0
	}
	return node.sum
}

func updateNode(node *treapNode) {
	if node == nil {
		return
	}
	node.sum = node.key + nodeSum(node.left) + nodeSum(node.right)
}

func split(root *treapNode, key int64) (*treapNode, *treapNode) {
	if root == nil {
		return nil, nil
	}

	if root.key < key {
		left, right := split(root.right, key)
		root.right = left
		updateNode(root)
		return root, right
	}

	left, right := split(root.left, key)
	root.left = right
	updateNode(root)
	return left, root
}

func merge(left, right *treapNode) *treapNode {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if left.priority > right.priority {
		left.right = merge(left.right, right)
		updateNode(left)
		return left
	}

	right.left = merge(left, right.left)
	updateNode(right)
	return right
}

type rangeSumSet struct {
	root *treapNode
	rng  *rand.Rand
}

func newRangeSumSet() *rangeSumSet {
	return &rangeSumSet{rng: rand.New(rand.NewSource(1))}
}

func (s *rangeSumSet) insert(value int64) {
	left, rest := split(s.root, value)
	middle, right := split(rest, value+1)
	if middle == nil {
		middle = &treapNode{key: value, priority: s.rng.Uint64()}
		middle.sum = value
	}
	s.root = merge(left, merge(middle, right))
}

func (s *rangeSumSet) erase(value int64) {
	left, rest := split(s.root, value)
	_, right := split(rest, value+1)
	s.root = merge(left, right)
}

func (s *rangeSumSet) contains(value int64) bool {
	node := s.root
	for node != nil {
		switch {
		case value == node.key:
			return true
		case value < node.key:
			node = node.left
		default:
			node = node.right
		}
	}
	return false
}

func (s *rangeSumSet) rangeSum(left, right int64) int64 {
	before, rest := split(s.root, left)
	middle, after := split(rest, right+1)
	result := nodeSum(middle)
	s.root = merge(before, merge(middle, after))
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	set := newRangeSumSet()
	var lastSum int64

	for i := 0; i < n; i++ {
		var operation string
		if _, err := fmt.Fscan(reader, &operation); err != nil {
			return
		}

		switch operation {
		case "+":
			var x int64
			fmt.Fscan(reader, &x)
			set.insert((x + lastSum) % rangeSumModulo)
		case "-":
			var x int64
			fmt.Fscan(reader, &x)
			set.erase((x + lastSum) % rangeSumModulo)
		case "?":
			var x int64
			fmt.Fscan(reader, &x)
			if set.contains((x + lastSum) % rangeSumModulo) {
				fmt.Println("Found")
			} else {
				fmt.Println("Not found")
			}
		case "s":
			var left, right int64
			fmt.Fscan(reader, &left, &right)
			result := set.rangeSum((left+lastSum)%rangeSumModulo, (right+lastSum)%rangeSumModulo)
			fmt.Println(result)
			lastSum = result % rangeSumModulo
		}
	}
}
