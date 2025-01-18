// go run binary_search.go

// Search a key in a sorted array of keys.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func binarySearch(keys []int, query int) int {
	left, right := 0, len(keys)-1

	for left <= right {
		mid := left + (right-left)/2
		switch {
		case keys[mid] == query:
			return mid
		case keys[mid] > query:
			right = mid - 1
		default:
			left = mid + 1
		}
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
		answers[i] = fmt.Sprint(binarySearch(keys, query))
	}

	fmt.Println(strings.Join(answers, " "))
}
