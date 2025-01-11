// go run max_pairwise_product.go

// find the maximum product of two distinct numbers in a sequence of non-negative integers

package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxPairwiseProduct(numbers []int64) int64 {
	var max1, max2 int64
	for _, num := range numbers {
		switch {
		case num > max1:
			max2 = max1
			max1 = num
		case num > max2:
			max2 = num
		}
	}
	return max1 * max2
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	numbers := make([]int64, n)
	for i := range numbers {
		if _, err := fmt.Fscan(reader, &numbers[i]); err != nil {
			return
		}
	}

	fmt.Println(maxPairwiseProduct(numbers))
}
