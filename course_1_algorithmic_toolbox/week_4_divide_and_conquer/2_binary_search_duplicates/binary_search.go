// go run binary_search.go

// Find the index of the first occurrence of a key in a sorted array.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func binarySearchFirst(keys []int, query int) int {
	left, right := 0, len(keys)

	for left < right {
		mid := left + (right-left)/2
		if keys[mid] < query {
			left = mid + 1
		} else {
			right = mid
		}
	}

	if left < len(keys) && keys[left] == query {
		return left
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var numKeys int
	if _, err := fmt.Fscan(reader, &numKeys); err != nil {
		return
	}
	keys := make([]int, numKeys)
	for i := range keys {
		if _, err := fmt.Fscan(reader, &keys[i]); err != nil {
			return
		}
	}

	var numQueries int
	if _, err := fmt.Fscan(reader, &numQueries); err != nil {
		return
	}
	answers := make([]string, numQueries)
	for i := range answers {
		var query int
		if _, err := fmt.Fscan(reader, &query); err != nil {
			return
		}
		answers[i] = fmt.Sprint(binarySearchFirst(keys, query))
	}

	fmt.Println(strings.Join(answers, " "))
}
