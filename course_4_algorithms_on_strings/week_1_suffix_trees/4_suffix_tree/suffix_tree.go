// go run suffix_tree.go

// Construct a suffix tree of a string and print all edge labels.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func buildSuffixTree(text string) []string {
	data := []byte(text)
	if len(data) == 0 {
		return nil
	}
	sa := SuffixArray(data)
	lcp := LCP(data, sa)
	tree := BuildSuffixTree(data, sa, lcp)
	return tree.EdgeLabels()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	labels := buildSuffixTree(text)
	sort.Strings(labels)
	for _, label := range labels {
		fmt.Println(label)
	}
}

// The helpers below are inlined from dsa/course_4_algorithms_on_strings/internal/strs
// so that this file can be submitted to Coursera as a single Go file.

// SuffixArray returns the starting positions of all suffixes of text in
// lexicographic order. It uses prefix doubling with counting sort in
// O(n log n) time.
func SuffixArray(text []byte) []int {
	n := len(text)
	if n == 0 {
		return nil
	}

	sa := make([]int, n)
	for i := range sa {
		sa[i] = i
	}
	sort.Slice(sa, func(i, j int) bool {
		return text[sa[i]] < text[sa[j]]
	})

	rank := make([]int, n)
	classes := 1
	rank[sa[0]] = 0
	for i := 1; i < n; i++ {
		if text[sa[i]] != text[sa[i-1]] {
			classes++
		}
		rank[sa[i]] = classes - 1
	}

	nextRank := make([]int, n)
	bySecond := make([]int, n)
	count := make([]int, n)

	for length := 1; length < n && classes < n; length *= 2 {
		position := 0
		for i := n - length; i < n; i++ {
			bySecond[position] = i
			position++
		}
		for _, suffix := range sa {
			if suffix >= length {
				bySecond[position] = suffix - length
				position++
			}
		}

		for i := 0; i < classes; i++ {
			count[i] = 0
		}
		for _, suffix := range bySecond {
			count[rank[suffix]]++
		}
		total := 0
		for i := 0; i < classes; i++ {
			frequency := count[i]
			count[i] = total
			total += frequency
		}
		for _, suffix := range bySecond {
			class := rank[suffix]
			sa[count[class]] = suffix
			count[class]++
		}

		nextRank[sa[0]] = 0
		newClasses := 1
		for i := 1; i < n; i++ {
			current, previous := sa[i], sa[i-1]
			same := rank[current] == rank[previous]
			if same {
				currentSecond, previousSecond := -1, -1
				if current+length < n {
					currentSecond = rank[current+length]
				}
				if previous+length < n {
					previousSecond = rank[previous+length]
				}
				same = currentSecond == previousSecond
			}
			if !same {
				newClasses++
			}
			nextRank[current] = newClasses - 1
		}

		rank, nextRank = nextRank, rank
		classes = newClasses
	}

	return sa
}

// LCP returns the longest-common-prefix array aligned with sa. LCP[0] is zero
// and LCP[i] is the LCP of suffixes sa[i-1] and sa[i].
func LCP(text []byte, sa []int) []int {
	n := len(text)
	if n == 0 {
		return nil
	}

	position := make([]int, n)
	for i, suffix := range sa {
		position[suffix] = i
	}

	lcp := make([]int, n)
	common := 0
	for suffix := 0; suffix < n; suffix++ {
		rank := position[suffix]
		if rank == 0 {
			common = 0
			continue
		}
		other := sa[rank-1]
		for suffix+common < n && other+common < n && text[suffix+common] == text[other+common] {
			common++
		}
		lcp[rank] = common
		if common > 0 {
			common--
		}
	}
	return lcp
}

// RMQ answers range-minimum queries on an integer slice.
type RMQ struct {
	size int
	tree []int
}

const rmqInfinity = int(^uint(0) >> 1)

// NewRMQ builds a range-minimum query structure over values.
func NewRMQ(values []int) *RMQ {
	size := 1
	for size < len(values) {
		size *= 2
	}
	tree := make([]int, 2*size)
	for i := range tree {
		tree[i] = rmqInfinity
	}
	for i, value := range values {
		tree[size+i] = value
	}
	for i := size - 1; i > 0; i-- {
		if tree[2*i] < tree[2*i+1] {
			tree[i] = tree[2*i]
		} else {
			tree[i] = tree[2*i+1]
		}
	}
	return &RMQ{size: size, tree: tree}
}

// Min returns the minimum of values[left:right]. The range must be nonempty.
func (r *RMQ) Min(left, right int) int {
	left += r.size
	right += r.size
	result := rmqInfinity
	for left < right {
		if left&1 == 1 {
			if r.tree[left] < result {
				result = r.tree[left]
			}
			left++
		}
		if right&1 == 1 {
			right--
			if r.tree[right] < result {
				result = r.tree[right]
			}
		}
		left /= 2
		right /= 2
	}
	return result
}

// SuffixTreeNode is one node of a suffix tree. EdgeEnd is exclusive.
type SuffixTreeNode struct {
	Parent    int
	Depth     int
	EdgeStart int
	EdgeEnd   int
	Children  map[byte]int
}

// SuffixTree stores a suffix tree over Text. Node zero is the root.
type SuffixTree struct {
	Text  []byte
	Nodes []SuffixTreeNode
	Root  int
}

// BuildSuffixTree constructs a suffix tree from text, its suffix array, and
// its LCP array. Internal edge starts are canonicalized to the earliest
// descendant suffix.
func BuildSuffixTree(text []byte, sa, lcp []int) *SuffixTree {
	tree := &SuffixTree{Text: text}
	tree.Nodes = append(tree.Nodes, SuffixTreeNode{
		Parent:    -1,
		Depth:     0,
		EdgeStart: -1,
		EdgeEnd:   -1,
		Children:  make(map[byte]int),
	})
	tree.Root = 0

	alignedLCP := alignLCP(lcp, len(sa))
	previousLCP := 0
	current := tree.Root

	for i, suffix := range sa {
		if i > 0 {
			previousLCP = alignedLCP[i]
		}
		for tree.Nodes[current].Depth > previousLCP {
			current = tree.Nodes[current].Parent
		}

		if tree.Nodes[current].Depth == previousLCP {
			current = tree.addLeaf(current, suffix)
		} else {
			edgeStart := sa[i-1] + tree.Nodes[current].Depth
			offset := previousLCP - tree.Nodes[current].Depth
			middle := tree.breakEdge(current, edgeStart, offset)
			current = tree.addLeaf(middle, suffix)
		}
	}

	tree.canonicalize()
	return tree
}

func alignLCP(lcp []int, n int) []int {
	aligned := make([]int, n)
	switch len(lcp) {
	case n:
		copy(aligned, lcp)
	case n - 1:
		for i := 1; i < n; i++ {
			aligned[i] = lcp[i-1]
		}
	}
	return aligned
}

func (t *SuffixTree) addNode(parent, depth, edgeStart, edgeEnd int) int {
	t.Nodes = append(t.Nodes, SuffixTreeNode{
		Parent:    parent,
		Depth:     depth,
		EdgeStart: edgeStart,
		EdgeEnd:   edgeEnd,
		Children:  make(map[byte]int),
	})
	return len(t.Nodes) - 1
}

func (t *SuffixTree) addLeaf(parent, suffix int) int {
	start := t.Nodes[parent].Depth + suffix
	leaf := t.addNode(parent, len(t.Text)-suffix, start, len(t.Text))
	t.Nodes[parent].Children[t.Text[start]] = leaf
	return leaf
}

func (t *SuffixTree) breakEdge(node, start, offset int) int {
	text := t.Text
	startChar := text[start]
	middleChar := text[start+offset]
	child := t.Nodes[node].Children[startChar]

	middle := t.addNode(node, t.Nodes[node].Depth+offset, start, start+offset)
	t.Nodes[middle].Children[middleChar] = child
	t.Nodes[child].Parent = middle
	t.Nodes[child].EdgeStart += offset
	t.Nodes[node].Children[startChar] = middle
	return middle
}

func (t *SuffixTree) canonicalize() {
	if len(t.Nodes) == 0 {
		return
	}

	order := make([]int, 0, len(t.Nodes))
	stack := []int{t.Root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, node)
		for _, child := range t.Nodes[node].Children {
			stack = append(stack, child)
		}
	}

	minSuffix := make([]int, len(t.Nodes))
	for i := len(order) - 1; i >= 0; i-- {
		node := order[i]
		if len(t.Nodes[node].Children) == 0 {
			minSuffix[node] = len(t.Text) - t.Nodes[node].Depth
			continue
		}
		best := len(t.Text) + 1
		for _, child := range t.Nodes[node].Children {
			if minSuffix[child] < best {
				best = minSuffix[child]
			}
		}
		minSuffix[node] = best
	}

	for node := range t.Nodes {
		if node == t.Root || len(t.Nodes[node].Children) == 0 {
			continue
		}
		parentDepth := t.Nodes[t.Nodes[node].Parent].Depth
		start := minSuffix[node] + parentDepth
		t.Nodes[node].EdgeStart = start
		t.Nodes[node].EdgeEnd = start + (t.Nodes[node].Depth - parentDepth)
	}
}

// EdgeLabels returns all edge labels of the tree in lexicographic order.
func (t *SuffixTree) EdgeLabels() []string {
	labels := make([]string, 0, len(t.Nodes)-1)
	for node := range t.Nodes {
		if node == t.Root {
			continue
		}
		labels = append(labels, string(t.Text[t.Nodes[node].EdgeStart:t.Nodes[node].EdgeEnd]))
	}
	sort.Strings(labels)
	return labels
}
