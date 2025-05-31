package main

import (
	"fmt"
	"strconv"
	"strings"
)

type suffixNode struct {
	parent    int
	depth     int
	edgeStart int
	children  []int
	leaf      bool
}

func checkSuffixTreeFromArray(input, expected, output []byte) (bool, string) {
	lines := nonEmptyLines(input)
	if len(lines) != 3 {
		return false, fmt.Sprintf("suffix-tree: got %d input lines, want 3", len(lines))
	}
	text := strings.TrimSpace(lines[0])
	sa, err := parseIntegerLine(lines[1])
	if err != nil {
		return false, fmt.Sprintf("suffix-tree: suffix array: %v", err)
	}
	lcp, err := parseIntegerLine(lines[2])
	if err != nil {
		return false, fmt.Sprintf("suffix-tree: LCP array: %v", err)
	}
	nodes, leafFor, err := buildSuffixTree(text, sa, lcp)
	if err != nil {
		return false, fmt.Sprintf("suffix-tree: %v", err)
	}
	want := suffixTreeLabels(text, nodes, leafFor, sa)
	got, err := parseSuffixTreeOutput(text, output)
	if err != nil {
		return false, fmt.Sprintf("suffix-tree: %v", err)
	}
	if len(got) != len(want) {
		return false, fmt.Sprintf("suffix-tree: got %d edges, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			return false, fmt.Sprintf("suffix-tree: edge %d is %q, want %q", i+1, got[i], want[i])
		}
	}
	return true, ""
}

func parseIntegerLine(line string) ([]int, error) {
	fields := strings.Fields(line)
	result := make([]int, len(fields))
	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

func buildSuffixTree(text string, sa, lcp []int) ([]suffixNode, []int, error) {
	n := len(text)
	if n == 0 {
		return nil, nil, fmt.Errorf("empty text")
	}
	if len(sa) != n {
		return nil, nil, fmt.Errorf("suffix array has %d entries, want %d", len(sa), n)
	}
	if len(lcp) != n-1 {
		return nil, nil, fmt.Errorf("LCP array has %d entries, want %d", len(lcp), n-1)
	}
	used := make([]bool, n)
	for _, suffix := range sa {
		if suffix < 0 || suffix >= n || used[suffix] {
			return nil, nil, fmt.Errorf("suffix array is not a permutation")
		}
		used[suffix] = true
	}
	for i, value := range lcp {
		if value < 0 || value > n-sa[i] || value > n-sa[i+1] {
			return nil, nil, fmt.Errorf("LCP %d is out of range", i)
		}
	}

	nodes := []suffixNode{{parent: -1, depth: 0, edgeStart: 0}}
	stack := []int{0}
	lastLeaf := -1
	leafFor := make([]int, n)
	for i, suffix := range sa {
		lcpValue := 0
		if i > 0 {
			lcpValue = lcp[i-1]
		}
		child := lastLeaf
		for len(stack) > 1 && nodes[stack[len(stack)-1]].depth > lcpValue {
			child = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		top := stack[len(stack)-1]
		if nodes[top].depth < lcpValue {
			if child < 0 || nodes[child].parent != top {
				return nil, nil, fmt.Errorf("cannot insert internal node at depth %d", lcpValue)
			}
			internal := len(nodes)
			nodes = append(nodes, suffixNode{
				parent:    top,
				depth:     lcpValue,
				edgeStart: nodes[child].edgeStart,
			})
			nodes[child].parent = internal
			nodes[child].edgeStart += lcpValue - nodes[top].depth
			nodes[top].children = removeChild(nodes[top].children, child)
			nodes[top].children = append(nodes[top].children, internal)
			nodes[internal].children = append(nodes[internal].children, child)
			stack = append(stack, internal)
			top = internal
		}
		leafDepth := n - suffix
		if leafDepth <= lcpValue {
			return nil, nil, fmt.Errorf("empty leaf edge at suffix %d", suffix)
		}
		leaf := len(nodes)
		nodes = append(nodes, suffixNode{
			parent:    top,
			depth:     leafDepth,
			edgeStart: suffix + lcpValue,
			leaf:      true,
		})
		nodes[top].children = append(nodes[top].children, leaf)
		leafFor[suffix] = leaf
		lastLeaf = leaf
	}
	return nodes, leafFor, nil
}

func removeChild(children []int, target int) []int {
	for i, child := range children {
		if child == target {
			return append(children[:i], children[i+1:]...)
		}
	}
	return children
}

func suffixTreeLabels(text string, nodes []suffixNode, leafFor, sa []int) []string {
	printed := make(map[int]bool, len(nodes))
	labels := make([]string, 0, len(nodes)-1)
	for _, suffix := range sa {
		path := make([]int, 0)
		for node := leafFor[suffix]; node != 0; node = nodes[node].parent {
			path = append(path, node)
		}
		for i := len(path) - 1; i >= 0; i-- {
			node := path[i]
			if printed[node] {
				continue
			}
			printed[node] = true
			parent := nodes[node].parent
			start := nodes[node].edgeStart
			end := start + nodes[node].depth - nodes[parent].depth
			labels = append(labels, text[start:end])
		}
	}
	return labels
}

func parseSuffixTreeOutput(text string, output []byte) ([]string, error) {
	lines := nonEmptyLines(output)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	start := 0
	first := strings.Fields(lines[0])
	if len(first) == 1 {
		if first[0] != text {
			return nil, fmt.Errorf("output text %q does not match input", first[0])
		}
		start = 1
	} else if len(first) != 2 {
		return nil, fmt.Errorf("first output line is not text or an edge")
	}
	labels := make([]string, 0, len(lines)-start)
	for i := start; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != 2 {
			return nil, fmt.Errorf("edge line %d needs two integers", i+1)
		}
		from, errFrom := strconv.Atoi(fields[0])
		to, errTo := strconv.Atoi(fields[1])
		if errFrom != nil || errTo != nil {
			return nil, fmt.Errorf("edge line %d is not numeric", i+1)
		}
		if from < 0 || from >= len(text) || to <= from || to > len(text) {
			return nil, fmt.Errorf("edge %d-%d is out of range", from, to)
		}
		labels = append(labels, text[from:to])
	}
	return labels, nil
}
