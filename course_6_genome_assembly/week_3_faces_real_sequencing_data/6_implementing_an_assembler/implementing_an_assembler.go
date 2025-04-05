// go run implementing_an_assembler.go

// Assemble contigs from single reads or read-pairs using de Bruijn graphs.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func parseAssemblerReads(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) > 0 {
		if count, err := strconv.Atoi(lines[0]); err == nil && count == len(lines)-1 {
			lines = lines[1:]
		}
	}
	reads := make([]string, 0)
	for _, line := range lines {
		parts := strings.Split(line, "|")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if i == len(parts)-1 && len(parts) > 1 {
				if _, err := strconv.Atoi(part); err == nil {
					continue
				}
			}
			reads = append(reads, part)
		}
	}
	return reads
}

func collectKMers(reads []string, k int, preserveDuplicates bool) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			kmer := read[i : i+k]
			if !preserveDuplicates {
				if seen[kmer] {
					continue
				}
				seen[kmer] = true
			}
			result = append(result, kmer)
		}
	}
	return result
}

func assembleContigs(reads []string) string {
	if len(reads) == 0 {
		return ""
	}
	k := len(reads[0])
	for _, read := range reads {
		if len(read) < k {
			k = len(read)
		}
	}
	if k < 2 {
		return ""
	}
	kmers := collectKMers(reads, k, len(reads) <= 100)
	contigs := assembly.UnitigContigs(kmers)
	var builder strings.Builder
	for i, contig := range contigs {
		fmt.Fprintf(&builder, ">CONTIG%d\n%s\n", i+1, contig)
	}
	return strings.TrimSuffix(builder.String(), "\n")
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	fmt.Println(assembleContigs(parseAssemblerReads(string(input))))
}
