// go run inversions.go

// Compute the number of inversions in a sequence of integers.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func countInversions(a []int) int64 {
	_, inversions := sortAndCount(append([]int(nil), a...))
	return inversions
}

func sortAndCount(a []int) ([]int, int64) {
	if len(a) <= 1 {
		return a, 0
	}

	mid := len(a) / 2
	left, leftInversions := sortAndCount(a[:mid])
	right, rightInversions := sortAndCount(a[mid:])

	merged := make([]int, 0, len(a))
	var splitInversions int64
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			splitInversions += int64(len(left) - i)
			j++
		}
	}
	merged = append(merged, left[i:]...)
	merged = append(merged, right[j:]...)

	return merged, leftInversions + rightInversions + splitInversions
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	elements := make([]int, n)
	for i := range elements {
		if _, err := fmt.Fscan(reader, &elements[i]); err != nil {
			return
		}
	}

	fmt.Println(countInversions(elements))
}
