// go run heavy_hitters.go

// Answer heavy-hitter queries by maintaining the exact frequency difference.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type heavyUpdate struct {
	id, value int
}

func solveHeavy(n, threshold int, updates []heavyUpdate, queries []int) string {
	frequency := make(map[int]int64)
	for _, update := range updates {
		frequency[update.id] += int64(update.value)
	}
	parts := make([]string, len(queries))
	for i, query := range queries {
		if frequency[query] >= int64(threshold) {
			parts[i] = "1"
		} else {
			parts[i] = "0"
		}
	}
	return strings.Join(parts, " ")
}

func parseHeavy(input string) (int, int, []heavyUpdate, []int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, threshold int
	if _, err := fmt.Fscan(reader, &n, &threshold); err != nil {
		return 0, 0, nil, nil, err
	}
	updates := make([]heavyUpdate, 0, 2*n)
	for i := 0; i < n; i++ {
		var update heavyUpdate
		if _, err := fmt.Fscan(reader, &update.id, &update.value); err != nil {
			return 0, 0, nil, nil, err
		}
		updates = append(updates, update)
	}
	for i := 0; i < n; i++ {
		var update heavyUpdate
		if _, err := fmt.Fscan(reader, &update.id, &update.value); err != nil {
			return 0, 0, nil, nil, err
		}
		update.value = -update.value
		updates = append(updates, update)
	}
	var queryCount int
	if _, err := fmt.Fscan(reader, &queryCount); err != nil {
		return 0, 0, nil, nil, err
	}
	queries := make([]int, queryCount)
	for i := range queries {
		if _, err := fmt.Fscan(reader, &queries[i]); err != nil {
			return 0, 0, nil, nil, err
		}
	}
	return n, threshold, updates, queries, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, threshold int
	if _, err := fmt.Fscan(reader, &n, &threshold); err != nil {
		return
	}
	updates := make([]heavyUpdate, 0, 2*n)
	for i := 0; i < n; i++ {
		var update heavyUpdate
		if _, err := fmt.Fscan(reader, &update.id, &update.value); err != nil {
			return
		}
		updates = append(updates, update)
	}
	for i := 0; i < n; i++ {
		var update heavyUpdate
		if _, err := fmt.Fscan(reader, &update.id, &update.value); err != nil {
			return
		}
		update.value = -update.value
		updates = append(updates, update)
	}
	var queryCount int
	if _, err := fmt.Fscan(reader, &queryCount); err != nil {
		return
	}
	queries := make([]int, queryCount)
	for i := range queries {
		if _, err := fmt.Fscan(reader, &queries[i]); err != nil {
			return
		}
	}
	fmt.Println(solveHeavy(n, threshold, updates, queries))
}
