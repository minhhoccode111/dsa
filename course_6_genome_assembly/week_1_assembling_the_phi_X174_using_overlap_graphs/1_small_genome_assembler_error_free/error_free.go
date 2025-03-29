// go run error_free.go

// Assemble a circular genome from error-free reads using overlap graphs.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func maxOverlapAbove(first, second string, current int) (int, bool) {
	for size := len(first) - 1; size > current; size-- {
		if strings.HasPrefix(second, first[len(first)-size:]) {
			return size, true
		}
	}
	return 0, false
}

func assembleErrorFree(reads []string) string {
	if len(reads) == 0 {
		return ""
	}
	if len(reads) <= 9 {
		if genome := assembly.CircularSuperstring(reads); genome != "" {
			return genome
		}
	}

	type link struct {
		overlap int
		next    int
	}
	links := make([]link, len(reads))
	for i := 0; i < len(reads)-1; i++ {
		for j := i + 1; j < len(reads); j++ {
			if overlap, ok := maxOverlapAbove(reads[i], reads[j], links[i].overlap); ok {
				links[i] = link{overlap: overlap, next: j}
			}
			if overlap, ok := maxOverlapAbove(reads[j], reads[i], links[j].overlap); ok {
				links[j] = link{overlap: overlap, next: i}
			}
		}
	}

	if len(reads) > 1 && links[0].next == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(reads[0])
	visited := make([]bool, len(reads))
	visited[0] = true
	overlap, current := links[0].overlap, links[0].next
	for current != 0 {
		if visited[current] {
			return ""
		}
		visited[current] = true
		builder.WriteString(reads[current][overlap:])
		overlap, current = links[current].overlap, links[current].next
	}
	genome := builder.String()
	if overlap > 0 && overlap < len(genome) {
		genome = genome[:len(genome)-overlap]
	}
	return genome
}

func readReads(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	reads := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			reads = append(reads, line)
		}
	}
	return reads
}

func main() {
	reads := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			reads = append(reads, line)
		}
	}
	if len(reads) > 0 {
		fmt.Println(assembleErrorFree(reads))
	}
}
