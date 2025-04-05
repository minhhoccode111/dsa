// go run bubble_detection.go

// Count bubble pairs in a de Bruijn graph.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type bubbleDetector struct {
	k            int
	threshold    int
	reads        []string
	visitedStack []bool
	vertexCounts []int
}

func (detector *bubbleDetector) findPaths(start int, adjacency [][]int) [][]int {
	if detector.vertexCounts == nil {
		detector.vertexCounts = make([]int, len(adjacency))
		detector.visitedStack = make([]bool, len(adjacency))
	}
	paths := make(map[string][]int)
	stack := []int{start}
	detector.visitedStack[start] = true
	add := func(path []int) {
		key := pathKey(path)
		if _, exists := paths[key]; !exists {
			copyPath := append([]int(nil), path...)
			paths[key] = copyPath
		}
	}
	for len(stack) > 0 {
		if len(stack) > 1 {
			add(stack)
		}
		current := stack[len(stack)-1]
		count := detector.vertexCounts[current]
		if count < len(adjacency[current]) {
			next := adjacency[current][count]
			detector.vertexCounts[current]++
			if !detector.visitedStack[next] {
				stack = append(stack, next)
				detector.visitedStack[next] = true
			}
		} else {
			detector.vertexCounts[current] = 0
			detector.visitedStack[current] = false
			stack = stack[:len(stack)-1]
		}
		if len(stack) == detector.threshold+1 {
			add(stack)
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			detector.visitedStack[last] = false
			detector.vertexCounts[last] = 0
		}
	}
	result := make([][]int, 0, len(paths))
	for _, path := range paths {
		result = append(result, path)
	}
	return result
}

func pathKey(path []int) string {
	var builder strings.Builder
	for _, vertex := range path {
		builder.WriteString(strconv.Itoa(vertex))
		builder.WriteByte(',')
	}
	return builder.String()
}

func countBubbles(k, threshold int, reads []string) int {
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
	adjacency := make([][]int, len(vertices))
	for _, kmer := range kmers {
		adjacency[index[kmer[:k-1]]] = append(adjacency[index[kmer[:k-1]]], index[kmer[1:]])
	}

	detector := &bubbleDetector{k: k, threshold: threshold, reads: reads}
	total := 0
	for start := 0; start < len(vertices); start++ {
		rawPaths := detector.findPaths(start, adjacency)
		endpointPaths := make(map[int][][]int)
		for _, path := range rawPaths {
			end := path[len(path)-1]
			endpointPaths[end] = append(endpointPaths[end], path)
		}
		for _, paths := range endpointPaths {
			if len(paths) <= 1 {
				continue
			}
			total += countBubblePairs(paths)
		}
	}
	return total
}

func countBubblePairs(paths [][]int) int {
	count := 0
	for i := 0; i < len(paths)-1; i++ {
		first := internalSet(paths[i])
		for j := i + 1; j < len(paths); j++ {
			second := internalSet(paths[j])
			disjoint := true
			for vertex := range first {
				if second[vertex] {
					disjoint = false
					break
				}
			}
			if disjoint {
				count++
			}
		}
	}
	return count
}

func internalSet(path []int) map[int]bool {
	set := make(map[int]bool)
	for i := 1; i+1 < len(path); i++ {
		set[path[i]] = true
	}
	return set
}

func parseBubbleInput(input string) (int, int, []string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return 0, 0, nil, fmt.Errorf("empty input")
	}
	var k, threshold int
	if _, err := fmt.Sscanf(lines[0], "%d %d", &k, &threshold); err != nil {
		return 0, 0, nil, err
	}
	return k, threshold, lines[1:], nil
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
	if len(reads) == 0 {
		return
	}
	var k, threshold int
	if _, err := fmt.Sscanf(reads[0], "%d %d", &k, &threshold); err != nil {
		return
	}
	fmt.Println(countBubbles(k, threshold, reads[1:]))
}
