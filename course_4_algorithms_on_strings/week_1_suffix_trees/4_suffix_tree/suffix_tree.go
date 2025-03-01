// go run suffix_tree.go

// Construct a suffix tree of a string and print all edge labels.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func buildSuffixTree(text string) []string {
	data := []byte(text)
	if len(data) == 0 {
		return nil
	}
	sa := strs.SuffixArray(data)
	lcp := strs.LCP(data, sa)
	tree := strs.BuildSuffixTree(data, sa, lcp)
	return tree.EdgeLabels()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	labels := buildSuffixTree(text)
	sort.Strings(labels)
	for _, label := range labels {
		fmt.Println(label)
	}
}
