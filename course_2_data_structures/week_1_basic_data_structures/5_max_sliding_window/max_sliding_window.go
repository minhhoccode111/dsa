// go run max_sliding_window.go

// Find the maximum in every sliding window in linear time.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxSlidingWindow(sequence []int, m int) []int {
	if len(sequence) == 0 || m <= 0 {
		return []int{}
	}

	result := make([]int, 0, len(sequence)-m+1)
	deque := make([]int, 0, len(sequence))
	head := 0

	for i, value := range sequence {
		for len(deque) > head && sequence[deque[len(deque)-1]] <= value {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)

		if deque[head] <= i-m {
			head++
		}

		if i >= m-1 {
			result = append(result, sequence[deque[head]])
		}
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	sequence := make([]int, n)
	for i := range sequence {
		if _, err := fmt.Fscan(reader, &sequence[i]); err != nil {
			return
		}
	}

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}

	result := maxSlidingWindow(sequence, m)
	for i, value := range result {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(value)
	}
	fmt.Println()
}
