// go run tip_removal.go

// Remove graph tips from a de Bruijn graph.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func removeTips(reads []string) int {
	if len(reads) == 0 {
		return 0
	}
	k := 15
	if len(reads[0])-1 < k {
		k = len(reads[0]) - 1
	}
	if k < 1 {
		return 0
	}
	kmerSet := make(map[string]bool)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			kmerSet[read[i:i+k]] = true
		}
	}
	kmers := make([]string, 0, len(kmerSet))
	for kmer := range kmerSet {
		kmers = append(kmers, kmer)
	}
	vertexSet := make(map[string]bool)
	for _, kmer := range kmers {
		vertexSet[kmer[:k-1]] = true
		vertexSet[kmer[1:]] = true
	}
	vertices := make([]string, 0, len(vertexSet))
	for vertex := range vertexSet {
		vertices = append(vertices, vertex)
	}
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		edges = append(edges, [2]int{index[kmer[:k-1]], index[kmer[1:]]})
	}

	tipCount := 0
	for {
		forward := make([][]int, len(vertices))
		reverse := make([][]int, len(vertices))
		for _, edge := range edges {
			forward[edge[0]] = append(forward[edge[0]], edge[1])
			reverse[edge[1]] = append(reverse[edge[1]], edge[0])
		}
		remove := make([]bool, len(edges))
		removed := 0
		for i, edge := range edges {
			if len(forward[edge[1]]) == 0 || len(reverse[edge[0]]) == 0 {
				remove[i] = true
				removed++
			}
		}
		if removed == 0 {
			break
		}
		tipCount += removed
		remaining := make([][2]int, 0, len(edges)-removed)
		for i, edge := range edges {
			if !remove[i] {
				remaining = append(remaining, edge)
			}
		}
		edges = remaining
	}
	return tipCount
}

func parseTipReads(input string) []string {
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
	fmt.Println(removeTips(reads))
}
