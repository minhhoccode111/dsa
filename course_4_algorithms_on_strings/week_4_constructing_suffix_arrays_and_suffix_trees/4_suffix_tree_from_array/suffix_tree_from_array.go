// go run suffix_tree_from_array.go

// Construct a suffix tree from a suffix array and an LCP array, then print its
// edges in depth-first order.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

type suffixTreeEdge struct {
	node       int
	start, end int
}

func suffixArrayToSuffixTree(sa, lcp []int, text string) map[int][]suffixTreeEdge {
	tree := strs.BuildSuffixTree([]byte(text), sa, lcp)
	if len(tree.Nodes) == 0 {
		return nil
	}

	ids := make([]int, len(tree.Nodes))
	ids[tree.Root] = 0
	nextID := 1
	edges := make(map[int][]suffixTreeEdge)

	queue := []int{tree.Root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if len(tree.Nodes[current].Children) == 0 {
			continue
		}

		chars := make([]byte, 0, len(tree.Nodes[current].Children))
		for char := range tree.Nodes[current].Children {
			chars = append(chars, char)
		}
		sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

		for _, char := range chars {
			child := tree.Nodes[current].Children[char]
			ids[child] = nextID
			nextID++
			edges[ids[current]] = append(edges[ids[current]], suffixTreeEdge{
				node:  ids[child],
				start: tree.Nodes[child].EdgeStart,
				end:   tree.Nodes[child].EdgeEnd,
			})
			queue = append(queue, child)
		}
	}

	return edges
}

func outputTree(edges map[int][]suffixTreeEdge) []string {
	type frame struct {
		node, edgeIndex int
	}

	stack := []frame{{node: 0, edgeIndex: 0}}
	var lines []string
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		nodeEdges, ok := edges[current.node]
		if !ok {
			continue
		}
		if current.edgeIndex+1 < len(nodeEdges) {
			stack = append(stack, frame{node: current.node, edgeIndex: current.edgeIndex + 1})
		}
		edge := nodeEdges[current.edgeIndex]
		lines = append(lines, fmt.Sprintf("%d %d", edge.start, edge.end))
		stack = append(stack, frame{node: edge.node, edgeIndex: 0})
	}
	return lines
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	sa := make([]int, len(text))
	for i := range sa {
		if _, err := fmt.Fscan(reader, &sa[i]); err != nil {
			return
		}
	}
	lcp := make([]int, len(text)-1)
	for i := range lcp {
		if _, err := fmt.Fscan(reader, &lcp[i]); err != nil {
			return
		}
	}

	fmt.Println(text)
	tree := suffixArrayToSuffixTree(sa, lcp, text)
	for _, line := range outputTree(tree) {
		fmt.Println(line)
	}
}
