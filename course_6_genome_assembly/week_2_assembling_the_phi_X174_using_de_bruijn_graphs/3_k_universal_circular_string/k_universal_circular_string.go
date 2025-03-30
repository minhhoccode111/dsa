// go run k_universal_circular_string.go

// Find a binary de Bruijn sequence containing every k-mer exactly once.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func binaryStrings(length int) []string {
	result := make([]string, 1<<length)
	for value := 0; value < 1<<length; value++ {
		result[value] = fmt.Sprintf("%0*b", length, value)
	}
	return result
}

func universalString(k int) string {
	if k == 1 {
		return "01"
	}
	kmers := binaryStrings(k)
	vertices := binaryStrings(k - 1)
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		edges = append(edges, [2]int{index[kmer[:k-1]], index[kmer[1:]]})
	}
	path, ok := assembly.FindEulerianCycle(len(vertices), edges)
	if !ok {
		return ""
	}
	var builder strings.Builder
	for i := 0; i+1 < len(path); i++ {
		vertex := vertices[path[i]]
		builder.WriteByte(vertex[len(vertex)-1])
	}
	return builder.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return
	}
	k, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return
	}
	fmt.Println(universalString(k))
}
