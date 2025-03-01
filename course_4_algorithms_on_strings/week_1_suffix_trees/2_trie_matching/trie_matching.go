// go run trie_matching.go

// Find all starting positions where any pattern occurs in the text.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type matchNode struct {
	next       [4]int
	patternEnd bool
}

func newMatchNode() matchNode {
	return matchNode{next: [4]int{-1, -1, -1, -1}}
}

func charIndex(char byte) int {
	switch char {
	case 'A':
		return 0
	case 'C':
		return 1
	case 'G':
		return 2
	case 'T':
		return 3
	default:
		return -1
	}
}

func buildMatchTrie(patterns []string) []matchNode {
	nodes := []matchNode{newMatchNode()}
	for _, pattern := range patterns {
		current := 0
		for i := 0; i < len(pattern); i++ {
			index := charIndex(pattern[i])
			if index < 0 {
				continue
			}
			if nodes[current].next[index] == -1 {
				nodes[current].next[index] = len(nodes)
				nodes = append(nodes, newMatchNode())
			}
			current = nodes[current].next[index]
		}
		nodes[current].patternEnd = true
	}
	return nodes
}

func solve(text string, patterns []string) []int {
	if len(patterns) == 0 || len(text) == 0 {
		return nil
	}

	trie := buildMatchTrie(patterns)
	maxLength := 0
	for _, pattern := range patterns {
		if len(pattern) > maxLength {
			maxLength = len(pattern)
		}
	}

	var result []int
	for start := 0; start < len(text); start++ {
		node := 0
		limit := start + maxLength
		if limit > len(text) {
			limit = len(text)
		}
		for position := start; position < limit; position++ {
			index := charIndex(text[position])
			if index < 0 || trie[node].next[index] == -1 {
				break
			}
			node = trie[node].next[index]
			if trie[node].patternEnd {
				result = append(result, start)
				break
			}
		}
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	var count int
	if _, err := fmt.Fscan(reader, &count); err != nil {
		return
	}
	patterns := make([]string, count)
	for i := range patterns {
		if _, err := fmt.Fscan(reader, &patterns[i]); err != nil {
			return
		}
	}

	for i, position := range solve(text, patterns) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(position)
	}
	fmt.Println()
}
