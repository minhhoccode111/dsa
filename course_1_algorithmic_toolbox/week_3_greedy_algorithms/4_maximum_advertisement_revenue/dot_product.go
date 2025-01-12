// go run dot_product.go

// Find the maximum dot product of two sequences of numbers.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func maxDotProduct(first, second []int) int64 {
	sort.Slice(first, func(i, j int) bool { return first[i] > first[j] })
	sort.Slice(second, func(i, j int) bool { return second[i] > second[j] })

	var total int64
	for i := range first {
		total += int64(first[i]) * int64(second[i])
	}
	return total
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	prices := make([]int, n)
	for i := range prices {
		if _, err := fmt.Fscan(reader, &prices[i]); err != nil {
			return
		}
	}

	clicks := make([]int, n)
	for i := range clicks {
		if _, err := fmt.Fscan(reader, &clicks[i]); err != nil {
			return
		}
	}

	fmt.Println(maxDotProduct(prices, clicks))
}
