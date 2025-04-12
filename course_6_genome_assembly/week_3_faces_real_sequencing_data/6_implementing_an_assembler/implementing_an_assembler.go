// go run implementing_an_assembler.go

// Assemble contigs from single reads or read-pairs using de Bruijn graphs.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
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
	contigs := UnitigContigs(kmers)
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

// The helpers below are inlined from dsa/course_6_genome_assembly/internal/assembly
// so that this file can be submitted to Coursera as a single Go file.

// FindEulerianCycle returns an Eulerian cycle as a sequence of vertices. The
// returned path repeats the starting vertex at the end.
func FindEulerianCycle(n int, edges [][2]int) ([]int, bool) {
	if !hasEulerianCycle(n, edges) {
		return nil, false
	}
	if len(edges) == 0 {
		return nil, true
	}

	adjacency := make([][]int, n)
	for id, edge := range edges {
		adjacency[edge[0]] = append(adjacency[edge[0]], id)
	}

	explored := 0
	previous := 0
	cycles := make([][]int, 0)
	for explored < len(edges) {
		edgeID := -1
		for vertex := previous; vertex < n; vertex++ {
			if len(adjacency[vertex]) > 0 {
				edgeID = adjacency[vertex][len(adjacency[vertex])-1]
				adjacency[vertex] = adjacency[vertex][:len(adjacency[vertex])-1]
				explored++
				previous = vertex
				break
			}
		}
		if edgeID < 0 {
			return nil, false
		}
		start, next := edges[edgeID][0], edges[edgeID][1]
		cycle := []int{start, next}
		for next != start {
			if len(adjacency[next]) == 0 {
				return nil, false
			}
			edgeID = adjacency[next][len(adjacency[next])-1]
			adjacency[next] = adjacency[next][:len(adjacency[next])-1]
			explored++
			next = edges[edgeID][1]
			cycle = append(cycle, next)
		}
		cycles = append(cycles, cycle)
	}

	cycleStarts := make(map[int][]int)
	for i, cycle := range cycles {
		cycleStarts[cycle[0]] = append(cycleStarts[cycle[0]], i)
	}

	path := make([]int, 0)
	stack := newDeque()
	stack.pushBack(cycles[0][0])
	for stack.len() > 0 {
		current := stack.popFront()
		if ids := cycleStarts[current]; len(ids) > 0 {
			id := ids[len(ids)-1]
			cycleStarts[current] = ids[:len(ids)-1]
			cycle := cycles[id]
			for i := len(cycle) - 1; i >= 1; i-- {
				stack.pushFront(cycle[i])
			}
		}
		path = append(path, current)
	}
	if len(edges) > 0 && len(path) != len(edges)+1 {
		return nil, false
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

type deque struct {
	front []int
	back  []int
}

func newDeque() *deque {
	return &deque{}
}

func (d *deque) len() int {
	return len(d.front) + len(d.back)
}

func (d *deque) pushFront(value int) {
	d.front = append(d.front, value)
}

func (d *deque) pushBack(value int) {
	d.back = append(d.back, value)
}

func (d *deque) popFront() int {
	if len(d.front) > 0 {
		value := d.front[len(d.front)-1]
		d.front = d.front[:len(d.front)-1]
		return value
	}
	value := d.back[0]
	d.back = d.back[1:]
	return value
}
