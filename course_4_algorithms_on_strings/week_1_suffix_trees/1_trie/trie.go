// go run trie.go

// Build a trie from a collection of patterns and print its edges.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type trieNode struct {
	next  map[byte]int
	order []byte
}

type trieEdge struct {
	from, to int
	char     byte
}

func buildTrie(patterns []string) []trieEdge {
	nodes := []trieNode{{next: make(map[byte]int)}}
	var edges []trieEdge

	for _, pattern := range patterns {
		current := 0
		for i := 0; i < len(pattern); i++ {
			char := pattern[i]
			if next, ok := nodes[current].next[char]; ok {
				current = next
				continue
			}
			next := len(nodes)
			nodes[current].next[char] = next
			nodes[current].order = append(nodes[current].order, char)
			nodes = append(nodes, trieNode{next: make(map[byte]int)})
			edges = append(edges, trieEdge{from: current, to: next, char: char})
			current = next
		}
	}

	return edges
}

func formatTrie(edges []trieEdge) []string {
	if len(edges) == 0 {
		return nil
	}

	children := make(map[int][]trieEdge)
	for _, edge := range edges {
		children[edge.from] = append(children[edge.from], edge)
	}
	for node := range children {
		sort.Slice(children[node], func(i, j int) bool {
			return children[node][i].char < children[node][j].char
		})
	}

	var lines []string
	var visit func(int)
	visit = func(node int) {
		for _, edge := range children[node] {
			lines = append(lines, fmt.Sprintf("%d->%d:%c", edge.from, edge.to, edge.char))
			visit(edge.to)
		}
	}
	visit(0)
	return lines
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var count int
	if _, err := fmt.Fscan(reader, &count); err != nil {
		return
	}
	patterns := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var pattern string
		if _, err := fmt.Fscan(reader, &pattern); err != nil {
			return
		}
		patterns = append(patterns, pattern)
	}

	for _, line := range formatTrie(buildTrie(patterns)) {
		fmt.Println(line)
	}
}
