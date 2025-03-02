// go run bwtinverse.go

// Reconstruct a string from its Burrows-Wheeler transform.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func inverseBWT(encoded string) string {
	n := len(encoded)
	if n == 0 {
		return ""
	}
	data := []byte(encoded)

	var counts [256]int
	for _, char := range data {
		counts[char]++
	}
	var starts [256]int
	total := 0
	for char := 0; char < 256; char++ {
		starts[char] = total
		total += counts[char]
	}

	var seen [256]int
	next := make([]int, n)
	sentinel := -1
	for i, char := range data {
		next[i] = starts[char] + seen[char]
		seen[char]++
		if char == '$' && sentinel < 0 {
			sentinel = i
		}
	}
	if sentinel < 0 {
		return ""
	}

	result := make([]byte, n)
	row := sentinel
	for i := n - 1; i >= 0; i-- {
		result[i] = data[row]
		row = next[row]
	}
	return string(result)
}

func bwt(text string) string {
	n := len(text)
	if n == 0 {
		return ""
	}
	data := []byte(text)
	sa := strs.SuffixArray(data)
	result := make([]byte, n)
	for i, suffix := range sa {
		previous := suffix - 1
		if previous < 0 {
			previous = n - 1
		}
		result[i] = data[previous]
	}
	return string(result)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var encoded string
	if _, err := fmt.Fscan(reader, &encoded); err != nil {
		return
	}
	fmt.Println(inverseBWT(encoded))
}
