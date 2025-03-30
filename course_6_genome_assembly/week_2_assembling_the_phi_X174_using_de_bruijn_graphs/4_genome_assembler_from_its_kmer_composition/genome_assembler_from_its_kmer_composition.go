// go run genome_assembler_from_its_kmer_composition.go

// Reconstruct a circular genome from its k-mer composition.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func assembleFromKmers(kmers []string) string {
	if len(kmers) == 0 {
		return ""
	}
	vertices := assembly.KMinusOneMers(kmers)
	edges := assembly.DeBruijnEdges(kmers, vertices)
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

func parseKmers(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) > 1 {
		if count, err := strconv.Atoi(lines[0]); err == nil && count == len(lines)-1 {
			return lines[1:]
		}
	}
	return lines
}

func main() {
	kmers := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			kmers = append(kmers, line)
		}
	}
	if len(kmers) > 1 {
		if count, err := strconv.Atoi(kmers[0]); err == nil && count == len(kmers)-1 {
			kmers = kmers[1:]
		}
	}
	fmt.Println(assembleFromKmers(kmers))
}
