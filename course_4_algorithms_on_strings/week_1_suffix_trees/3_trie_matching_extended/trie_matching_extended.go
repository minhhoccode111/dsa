// go run trie_matching_extended.go

// Find all starting positions where any pattern occurs in the text, including
// patterns that are prefixes of other patterns.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type extendedNode struct {
	next       [4]int
	patternEnd bool
}

func newExtendedNode() extendedNode {
	return extendedNode{next: [4]int{-1, -1, -1, -1}}
}

func extendedCharIndex(char byte) int {
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

func buildExtendedTrie(patterns []string) []extendedNode {
	nodes := []extendedNode{newExtendedNode()}
	for _, pattern := range patterns {
		current := 0
		for i := 0; i < len(pattern); i++ {
			index := extendedCharIndex(pattern[i])
			if index < 0 {
				continue
			}
			if nodes[current].next[index] == -1 {
				nodes[current].next[index] = len(nodes)
				nodes = append(nodes, newExtendedNode())
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

	trie := buildExtendedTrie(patterns)
	var result []int
	for start := 0; start < len(text); start++ {
		node := 0
		for position := start; position < len(text); position++ {
			index := extendedCharIndex(text[position])
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

	result := solve(text, patterns)
	for i, position := range result {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(position)
	}
	fmt.Println()
}
