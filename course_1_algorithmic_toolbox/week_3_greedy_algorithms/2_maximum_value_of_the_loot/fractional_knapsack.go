// go run fractional_knapsack.go

// Find the maximum value of items that fit into the backpack.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type lootItem struct {
	value  int
	weight int
}

func fractionalKnapsack(capacity int, weights, values []int) float64 {
	items := make([]lootItem, len(weights))
	for i := range weights {
		items[i] = lootItem{value: values[i], weight: weights[i]}
	}

	sort.Slice(items, func(i, j int) bool {
		return int64(items[i].value)*int64(items[j].weight) >
			int64(items[j].value)*int64(items[i].weight)
	})

	total := 0.0
	for _, item := range items {
		if capacity == 0 {
			break
		}

		amount := item.weight
		if amount > capacity {
			amount = capacity
		}

		total += float64(item.value) * float64(amount) / float64(item.weight)
		capacity -= amount
	}

	return total
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, capacity int
	if _, err := fmt.Fscan(reader, &n, &capacity); err != nil {
		return
	}

	values := make([]int, n)
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &values[i], &weights[i]); err != nil {
			return
		}
	}

	fmt.Printf("%.4f\n", fractionalKnapsack(capacity, weights, values))
}
