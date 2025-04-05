// go run optimal_kmer_size.go

// Find the largest k whose de Bruijn graph has an Eulerian cycle.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func hasEulerianCycle(vertices []string, kmers []string) bool {
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	inDegree := make([]int, len(vertices))
	outDegree := make([]int, len(vertices))
	for _, kmer := range kmers {
		outDegree[index[kmer[:len(kmer)-1]]]++
		inDegree[index[kmer[1:]]]++
	}
	for i := range vertices {
		if inDegree[i] != outDegree[i] {
			return false
		}
	}
	return true
}

func optimalKmerSize(reads []string) int {
	if len(reads) == 0 {
		return 0
	}
	left, right := 3, len(reads[0])
	for left <= right {
		current := (left + right) / 2
		kmerSet := make(map[string]bool)
		for _, read := range reads {
			for position := 0; position+current <= len(read); position++ {
				kmerSet[read[position:position+current]] = true
			}
		}
		kmers := make([]string, 0, len(kmerSet))
		for kmer := range kmerSet {
			kmers = append(kmers, kmer)
		}
		vertexSet := make(map[string]bool)
		for _, kmer := range kmers {
			vertexSet[kmer[:len(kmer)-1]] = true
			vertexSet[kmer[1:]] = true
		}
		vertices := make([]string, 0, len(vertexSet))
		for vertex := range vertexSet {
			vertices = append(vertices, vertex)
		}
		if hasEulerianCycle(vertices, kmers) {
			left = current + 1
		} else {
			right = current - 1
		}
	}
	return right
}

func parseSequences(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	result := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			result = append(result, line)
		}
	}
	return result
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
	fmt.Println(optimalKmerSize(reads))
}
