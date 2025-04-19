// go run kmp.go

// Find all occurrences of a pattern in a text using the KMP algorithm.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func computePrefixFunction(pattern string) []int {
	n := len(pattern)
	if n == 0 {
		return nil
	}
	prefix := make([]int, n)
	border := 0
	for i := 1; i < n; i++ {
		for border > 0 && pattern[i] != pattern[border] {
			border = prefix[border-1]
		}
		if pattern[i] == pattern[border] {
			border++
		} else {
			border = 0
		}
		prefix[i] = border
	}
	return prefix
}

func findPattern(pattern, text string) []int {
	combined := pattern + "$" + text
	prefix := computePrefixFunction(combined)
	var result []int
	for i := len(pattern) + 1; i < len(combined); i++ {
		if prefix[i] == len(pattern) {
			result = append(result, i-2*len(pattern))
		}
	}
	return result
}

func writePositions(w io.Writer, positions []int) {
	buffered := bufio.NewWriterSize(w, 1<<20)
	var buf [20]byte
	for i, position := range positions {
		if i > 0 {
			buffered.WriteByte(' ')
		}
		buffered.Write(strconv.AppendInt(buf[:0], int64(position), 10))
	}
	buffered.WriteByte('\n')
	buffered.Flush()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var pattern, text string
	if _, err := fmt.Fscan(reader, &pattern); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	writePositions(os.Stdout, findPattern(pattern, text))
}
