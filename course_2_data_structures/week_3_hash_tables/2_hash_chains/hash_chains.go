// go run hash_chains.go

// Implement a hash table with chaining.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type hashQuery struct {
	kind  string
	index int
	value string
}

const (
	hashMultiplier int64 = 263
	hashPrime      int64 = 1000000007
)

func hashString(s string, bucketCount int) int {
	var value int64
	for i := len(s) - 1; i >= 0; i-- {
		value = (value*hashMultiplier + int64(s[i])) % hashPrime
	}
	return int(value % int64(bucketCount))
}

func contains(chain []string, value string) bool {
	for _, item := range chain {
		if item == value {
			return true
		}
	}
	return false
}

func processHashQueries(bucketCount int, queries []hashQuery) []string {
	buckets := make([][]string, bucketCount)
	result := make([]string, 0)

	for _, query := range queries {
		switch query.kind {
		case "check":
			result = append(result, strings.Join(buckets[query.index], " "))
		case "add":
			bucket := hashString(query.value, bucketCount)
			if !contains(buckets[bucket], query.value) {
				buckets[bucket] = append([]string{query.value}, buckets[bucket]...)
			}
		case "del":
			bucket := hashString(query.value, bucketCount)
			for i, item := range buckets[bucket] {
				if item == query.value {
					buckets[bucket] = append(buckets[bucket][:i], buckets[bucket][i+1:]...)
					break
				}
			}
		case "find":
			bucket := hashString(query.value, bucketCount)
			if contains(buckets[bucket], query.value) {
				result = append(result, "yes")
			} else {
				result = append(result, "no")
			}
		}
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var bucketCount int
	if _, err := fmt.Fscan(reader, &bucketCount); err != nil {
		return
	}

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	queries := make([]hashQuery, 0, n)
	for i := 0; i < n; i++ {
		var query hashQuery
		if _, err := fmt.Fscan(reader, &query.kind); err != nil {
			return
		}
		if query.kind == "check" {
			if _, err := fmt.Fscan(reader, &query.index); err != nil {
				return
			}
		} else {
			if _, err := fmt.Fscan(reader, &query.value); err != nil {
				return
			}
		}
		queries = append(queries, query)
	}

	result := processHashQueries(bucketCount, queries)
	if len(result) > 0 {
		fmt.Println(strings.Join(result, "\n"))
	}
}
