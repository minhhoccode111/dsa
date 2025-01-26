// go run build_heap.go

// Convert an array into a min-heap and report the performed swaps.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func siftDown(data []int, i, size int, swaps *[][2]int) {
	minIndex := i
	left := 2*i + 1
	right := 2*i + 2

	if left < size && data[left] < data[minIndex] {
		minIndex = left
	}
	if right < size && data[right] < data[minIndex] {
		minIndex = right
	}

	if minIndex != i {
		*swaps = append(*swaps, [2]int{i, minIndex})
		data[i], data[minIndex] = data[minIndex], data[i]
		siftDown(data, minIndex, size, swaps)
	}
}

func buildHeap(data []int) [][2]int {
	swaps := make([][2]int, 0)
	for i := len(data)/2 - 1; i >= 0; i-- {
		siftDown(data, i, len(data), &swaps)
	}
	return swaps
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	data := make([]int, n)
	for i := range data {
		if _, err := fmt.Fscan(reader, &data[i]); err != nil {
			return
		}
	}

	swaps := buildHeap(data)
	fmt.Println(len(swaps))
	for _, swap := range swaps {
		fmt.Println(swap[0], swap[1])
	}
}
