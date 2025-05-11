// go run tip_removal.go

// Remove graph tips from a de Bruijn graph.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// tipGraph is a de Bruijn graph over (k-1)-mers with out-neighbor lists in
// first-appearance order and live in-degree counters.
type tipGraph struct {
	out     [][]int
	inDeg   []int
	limit   int
	removed int
}

func (g *tipGraph) removeEdge(u, v int) {
	neighbors := g.out[u]
	for i, target := range neighbors {
		if target == v {
			g.out[u] = append(neighbors[:i:i], neighbors[i+1:]...)
			break
		}
	}
	g.inDeg[v]--
	g.removed++
}

// findRemoveOutgoing walks a tip that starts at a branching vertex: it follows
// the chain from current while it stays non-branching and short, and removes
// the chain only if it dead-ends within the depth limit.
func (g *tipGraph) findRemoveOutgoing(current, depth, limit int) bool {
	if len(g.out[current]) > 1 || g.inDeg[current] > 1 {
		return false
	}
	if depth == limit {
		return false
	}
	if len(g.out[current]) == 0 {
		return true
	}
	next := g.out[current][0]
	if g.findRemoveOutgoing(next, depth+1, limit) {
		g.removeEdge(current, next)
		return true
	}
	return false
}

// findRemoveIncoming walks a tip that starts at a vertex with no incoming
// edges: it removes the chain up to the first merging or dead-end vertex.
func (g *tipGraph) findRemoveIncoming(current, depth, limit int) bool {
	if depth == limit {
		return false
	}
	if len(g.out[current]) == 0 || g.inDeg[current] > 1 {
		return true
	}
	next := g.out[current][0]
	if g.findRemoveIncoming(next, depth+1, limit) {
		g.removeEdge(current, next)
		return true
	}
	return false
}

func (g *tipGraph) removeTips() int {
	for u := range g.out {
		incoming := false
		if len(g.out[u]) == 1 && g.inDeg[u] == 0 {
			incoming = true
		} else if len(g.out[u]) <= 1 {
			continue
		}
		for {
			progress := false
			for _, w := range g.out[u] {
				var ok bool
				if incoming {
					ok = g.findRemoveIncoming(w, 0, g.limit)
				} else {
					ok = g.findRemoveOutgoing(w, 0, g.limit)
				}
				if ok {
					g.removeEdge(u, w)
					progress = true
					break
				}
			}
			if !progress {
				break
			}
		}
	}
	return g.removed
}

type tipGraphBuilder struct {
	index   map[string]int
	graph   tipGraph
	present []map[int]bool
}

func (b *tipGraphBuilder) vertex(name string) int {
	if id, ok := b.index[name]; ok {
		return id
	}
	id := len(b.index)
	b.index[name] = id
	b.graph.out = append(b.graph.out, nil)
	b.graph.inDeg = append(b.graph.inDeg, 0)
	b.present = append(b.present, nil)
	return id
}

func (b *tipGraphBuilder) addEdge(left, right int) {
	if b.present[left] == nil {
		b.present[left] = make(map[int]bool, 4)
	}
	if b.present[left][right] {
		return
	}
	b.present[left][right] = true
	b.graph.out[left] = append(b.graph.out[left], right)
	b.graph.inDeg[right]++
}

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
	builder := &tipGraphBuilder{index: make(map[string]int)}
	builder.graph.limit = k
	for _, read := range reads {
		for start := 0; start+k <= len(read); start++ {
			kmer := read[start : start+k]
			left := kmer[:k-1]
			right := kmer[1:]
			if left == right {
				continue
			}
			builder.addEdge(builder.vertex(left), builder.vertex(right))
		}
	}
	return builder.graph.removeTips()
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
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			reads = append(reads, line)
		}
	}
	fmt.Println(removeTips(reads))
}
