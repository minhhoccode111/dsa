// go run non_shared_substring.go

// Find the shortest substring of the first string that does not occur in the
// second string. Ties are broken lexicographically.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func chooseSeparators(p, q string) (byte, byte) {
	var used [256]bool
	for i := 0; i < len(p); i++ {
		used[p[i]] = true
	}
	for i := 0; i < len(q); i++ {
		used[q[i]] = true
	}
	first := byte(0)
	for i := 0; i < 256; i++ {
		if !used[i] {
			first = byte(i)
			used[i] = true
			break
		}
	}
	second := byte(0)
	for i := 0; i < 256; i++ {
		if !used[i] {
			second = byte(i)
			break
		}
	}
	return first, second
}

func solve(p, q string) string {
	if len(p) == 0 {
		return ""
	}

	firstSeparator, secondSeparator := chooseSeparators(p, q)
	combined := make([]byte, 0, len(p)+len(q)+2)
	combined = append(combined, q...)
	combined = append(combined, firstSeparator)
	pStart := len(combined)
	combined = append(combined, p...)
	combined = append(combined, secondSeparator)

	sa := strs.SuffixArray(combined)
	lcp := strs.LCP(combined, sa)
	rmq := strs.NewRMQ(lcp)

	rank := make([]int, len(combined))
	for i, suffix := range sa {
		rank[suffix] = i
	}

	qEnd := len(q)
	previousQ := make([]int, len(sa))
	last := -1
	for i, suffix := range sa {
		if suffix < qEnd {
			last = i
		}
		previousQ[i] = last
	}

	nextQ := make([]int, len(sa))
	next := -1
	for i := len(sa) - 1; i >= 0; i-- {
		if sa[i] < qEnd {
			next = i
		}
		nextQ[i] = next
	}

	best := ""
	for offset := 0; offset < len(p); offset++ {
		position := pStart + offset
		suffixRank := rank[position]

		matched := 0
		if left := previousQ[suffixRank]; left >= 0 && left < suffixRank {
			if common := rmq.Min(left+1, suffixRank+1); common > matched {
				matched = common
			}
		}
		if right := nextQ[suffixRank]; right > suffixRank {
			if common := rmq.Min(suffixRank+1, right+1); common > matched {
				matched = common
			}
		}

		remaining := len(p) - offset
		if matched >= remaining {
			continue
		}
		candidate := p[offset : offset+matched+1]
		if best == "" || len(candidate) < len(best) || (len(candidate) == len(best) && candidate < best) {
			best = candidate
		}
	}

	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var p, q string
	if _, err := fmt.Fscan(reader, &p); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	fmt.Println(solve(p, q))
}
