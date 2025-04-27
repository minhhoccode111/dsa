// go run error_prone.go

// Assemble a circular genome from error-prone reads using overlap graphs.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func hammingWithin(first, second string, errors int) bool {
	if len(first) != len(second) {
		return false
	}
	mismatches := 0
	for i := range first {
		if first[i] != second[i] {
			mismatches++
			if mismatches > errors {
				return false
			}
		}
	}
	return true
}

func removeDuplicates(reads []string, readLength, errors int) []string {
	skip := make([]bool, len(reads))
	for i := 0; i < len(reads)-1; i++ {
		if skip[i] {
			continue
		}
		for j := i + 1; j < len(reads); j++ {
			if skip[j] {
				continue
			}
			if hammingWithin(reads[i], reads[j], errors) {
				skip[j] = true
			}
		}
	}
	result := make([]string, 0, len(reads))
	for i, read := range reads {
		if !skip[i] {
			result = append(result, read)
		}
	}
	return result
}

func findOverlap(first, second string, previous, errors int) (int, bool) {
	for size := len(first) - 1; size > previous; size-- {
		mismatches := 0
		start := len(first) - size
		valid := true
		for i := 0; i < size; i++ {
			if first[start+i] != second[i] {
				mismatches++
				if mismatches > errors {
					valid = false
					break
				}
			}
		}
		if valid {
			return size, true
		}
	}
	return 0, false
}

func assembleGreedy(data []string, readLength, errors int) string {
	size := len(data)
	processed := make([]bool, size)
	overlap := make([]int, size)
	nextID := make([]int, size)
	for i := range nextID {
		nextID[i] = -1
	}

	current := 0
	for step := 0; step < size-1; step++ {
		processed[current] = true
		for j := 1; j < size; j++ {
			if processed[j] {
				continue
			}
			if found, ok := findOverlap(data[current], data[j], overlap[current], errors); ok {
				nextID[current] = j
				overlap[current] = found
			}
		}
		if nextID[current] < 0 {
			return ""
		}
		current = nextID[current]
	}

	nextID[current] = 0
	if found, ok := findOverlap(data[current], data[0], 0, errors); ok {
		overlap[current] = found
	} else {
		overlap[current] = 0
	}
	lastOverlap := overlap[current]

	type placement struct {
		read  int
		start int
	}
	placements := make([]placement, 0, 2*size)
	position := 0
	placements = append(placements, placement{read: 0, start: 0})
	current = 0
	for step := 0; step < size-1; step++ {
		position += readLength - overlap[current]
		placements = append(placements, placement{read: nextID[current], start: position})
		current = nextID[current]
	}

	genomeLength := position + readLength
	for step := size - 1; step < size*2; step++ {
		position += readLength - overlap[current]
		placements = append(placements, placement{read: nextID[current], start: position})
		current = nextID[current]
	}

	var builder strings.Builder
	for i := 0; i < genomeLength; i++ {
		counts := [256]int{}
		for _, placement := range placements {
			offset := i - placement.start
			if offset >= 0 && offset < readLength {
				counts[data[placement.read][offset]]++
			}
		}
		bestChar := byte(0)
		bestCount := 0
		for _, char := range []byte{'A', 'C', 'G', 'T'} {
			if counts[char] > bestCount {
				bestCount = counts[char]
				bestChar = char
			}
		}
		if bestChar == 0 {
			for char := 0; char < 256; char++ {
				if counts[char] > bestCount {
					bestCount = counts[char]
					bestChar = byte(char)
				}
			}
		}
		builder.WriteByte(bestChar)
	}
	result := builder.String()
	if lastOverlap > 0 && lastOverlap < len(result) {
		result = result[lastOverlap:]
	}
	return result
}

func assembleErrorProne(reads []string) string {
	if len(reads) == 0 {
		return ""
	}
	readLength := len(reads[0])
	errors := readLength / 50
	data := removeDuplicates(reads, readLength, errors)
	if len(data) <= 9 {
		if genome := CircularSuperstring(data); genome != "" {
			return genome
		}
	}
	return assembleGreedy(data, readLength, errors)
}

func parseReads(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 1 {
		if count, err := strconv.Atoi(lines[0]); err == nil && count == len(lines)-1 {
			return lines[1:]
		}
	}
	return lines
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
	if len(reads) > 1 {
		if count, err := strconv.Atoi(reads[0]); err == nil && count == len(reads)-1 {
			reads = reads[1:]
		}
	}
	if len(reads) > 0 {
		fmt.Println(assembleErrorProne(reads))
	}
}

// The helpers below are inlined from dsa/course_6_genome_assembly/internal/assembly
// so that this file can be submitted to Coursera as a single Go file.

// FindEulerianCycle returns an Eulerian cycle as a sequence of vertices. The
// returned path repeats the starting vertex at the end.
func FindEulerianCycle(n int, edges [][2]int) ([]int, bool) {
	if len(edges) == 0 {
		return nil, true
	}

	adjacency := make([][]int, n)
	for id, edge := range edges {
		if edge[0] < 0 || edge[0] >= n || edge[1] < 0 || edge[1] >= n {
			return nil, false
		}
		adjacency[edge[0]] = append(adjacency[edge[0]], id)
	}
	if !hasEulerianCycle(n, edges) {
		return nil, false
	}

	start := -1
	for vertex := 0; vertex < n; vertex++ {
		if len(adjacency[vertex]) > 0 {
			start = vertex
			break
		}
	}
	if start < 0 {
		return nil, false
	}

	stack := []int{start}
	path := make([]int, 0, len(edges)+1)
	for len(stack) > 0 {
		vertex := stack[len(stack)-1]
		if len(adjacency[vertex]) > 0 {
			edgeID := adjacency[vertex][len(adjacency[vertex])-1]
			adjacency[vertex] = adjacency[vertex][:len(adjacency[vertex])-1]
			stack = append(stack, edges[edgeID][1])
		} else {
			path = append(path, vertex)
			stack = stack[:len(stack)-1]
		}
	}
	if len(path) != len(edges)+1 {
		return nil, false
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, true
}

func hasEulerianCycle(n int, edges [][2]int) bool {
	inDegree := make([]int, n)
	outDegree := make([]int, n)
	for _, edge := range edges {
		outDegree[edge[0]]++
		inDegree[edge[1]]++
	}
	for i := 0; i < n; i++ {
		if inDegree[i] != outDegree[i] {
			return false
		}
	}
	return true
}

// MaxOverlap returns the longest suffix of first that is a prefix of second.
func MaxOverlap(first, second string) int {
	limit := len(first)
	if len(second) < limit {
		limit = len(second)
	}
	for size := limit; size > 0; size-- {
		if first[len(first)-size:] == second[:size] {
			return size
		}
	}
	return 0
}

// CircularSuperstring returns the shortest circular superstring of a small set
// of reads, or an empty string when the read count is too large.
func CircularSuperstring(reads []string) string {
	n := len(reads)
	if n == 0 {
		return ""
	}
	if n > 9 {
		return ""
	}
	used := make([]bool, n)
	order := make([]int, 0, n)
	bestLength := 1 << 30
	best := ""
	var search func()
	search = func() {
		if len(order) == n {
			total := 0
			length := 0
			for i, read := range order {
				total += len(reads[read])
				if i > 0 {
					total -= MaxOverlap(reads[order[i-1]], reads[read])
				}
				length = len(reads[read])
			}
			_ = length
			closing := MaxOverlap(reads[order[n-1]], reads[order[0]])
			total -= closing
			if total > bestLength {
				return
			}
			var builder strings.Builder
			builder.WriteString(reads[order[0]])
			for i := 1; i < n; i++ {
				overlap := MaxOverlap(reads[order[i-1]], reads[order[i]])
				builder.WriteString(reads[order[i]][overlap:])
			}
			genome := builder.String()
			if closing > 0 {
				genome = genome[:len(genome)-closing]
			}
			if total < bestLength || (total == bestLength && (best == "" || genome < best)) {
				bestLength = total
				best = genome
			}
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			used[i] = true
			order = append(order, i)
			search()
			order = order[:len(order)-1]
			used[i] = false
		}
	}
	search()
	return best
}

// UniqueKMers returns all distinct k-mers from the reads, sorted.
func UniqueKMers(reads []string, k int) []string {
	seen := make(map[string]bool)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			seen[read[i:i+k]] = true
		}
	}
	result := make([]string, 0, len(seen))
	for kmer := range seen {
		result = append(result, kmer)
	}
	sort.Strings(result)
	return result
}

// KMinusOneMers returns the sorted distinct prefixes and suffixes of length
// k-1 of the given k-mers.
func KMinusOneMers(kmers []string) []string {
	seen := make(map[string]bool)
	for _, kmer := range kmers {
		if len(kmer) < 2 {
			continue
		}
		seen[kmer[:len(kmer)-1]] = true
		seen[kmer[1:]] = true
	}
	result := make([]string, 0, len(seen))
	for mer := range seen {
		result = append(result, mer)
	}
	sort.Strings(result)
	return result
}

// DeBruijnEdges builds edges (prefix, suffix) for each k-mer.
func DeBruijnEdges(kmers, vertices []string) [][2]int {
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		if len(kmer) < 2 {
			continue
		}
		edges = append(edges, [2]int{index[kmer[:len(kmer)-1]], index[kmer[1:]]})
	}
	return edges
}

// UnitigContigs returns the strings spelled by maximal non-branching paths in
// the de Bruijn graph of the given k-mers. Duplicate edges produce duplicate
// contigs.
func UnitigContigs(kmers []string) []string {
	if len(kmers) == 0 {
		return nil
	}
	vertices := KMinusOneMers(kmers)
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		edges = append(edges, [2]int{index[kmer[:len(kmer)-1]], index[kmer[1:]]})
	}
	inDegree := make([]int, len(vertices))
	outDegree := make([]int, len(vertices))
	for _, edge := range edges {
		outDegree[edge[0]]++
		inDegree[edge[1]]++
	}

	used := make([]bool, len(edges))
	adjacency := make([][]int, len(vertices))
	for id, edge := range edges {
		adjacency[edge[0]] = append(adjacency[edge[0]], id)
	}

	contigs := make([]string, 0)
	for startEdge := 0; startEdge < len(edges); startEdge++ {
		if used[startEdge] {
			continue
		}
		start := edges[startEdge][0]
		if inDegree[start] == 1 && outDegree[start] == 1 {
			continue
		}
		currentEdge := startEdge
		var builder strings.Builder
		builder.WriteString(vertices[start])
		for {
			if used[currentEdge] {
				break
			}
			used[currentEdge] = true
			edge := edges[currentEdge]
			builder.WriteString(vertices[edge[1]][len(vertices[edge[1]])-1:])
			next := edge[1]
			if inDegree[next] != 1 || outDegree[next] != 1 {
				break
			}
			found := -1
			for _, candidate := range adjacency[next] {
				if !used[candidate] {
					found = candidate
					break
				}
			}
			if found < 0 {
				break
			}
			currentEdge = found
		}
		contigs = append(contigs, builder.String())
	}
	// Pure cycles have no vertex with in-degree or out-degree different from
	// one, so they are not discovered by the branch-based pass above.
	for startEdge := 0; startEdge < len(edges); startEdge++ {
		if used[startEdge] {
			continue
		}
		start := edges[startEdge][0]
		var builder strings.Builder
		builder.WriteString(vertices[start])
		currentEdge := startEdge
		lastVertex := start
		for !used[currentEdge] {
			used[currentEdge] = true
			edge := edges[currentEdge]
			builder.WriteString(vertices[edge[1]][len(vertices[edge[1]])-1:])
			lastVertex = edge[1]
			next := -1
			for _, candidate := range adjacency[edge[1]] {
				if !used[candidate] {
					next = candidate
					break
				}
			}
			if next < 0 {
				break
			}
			currentEdge = next
		}
		contig := builder.String()
		if lastVertex == start && len(contig) > len(vertices[start]) {
			contig = contig[:len(contig)-len(vertices[start])]
		}
		contigs = append(contigs, contig)
	}
	sort.Strings(contigs)
	return contigs
}
