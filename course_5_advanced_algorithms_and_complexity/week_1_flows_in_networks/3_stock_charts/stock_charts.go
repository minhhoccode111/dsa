// go run stock_charts.go

// Cover stock price charts with the minimum number of charts using maximum
// bipartite matching.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func dominates(first, second []int) bool {
	for i := range first {
		if first[i] >= second[i] {
			return false
		}
	}
	return true
}

func maxBipartiteMatching(adj [][]bool) int {
	n := len(adj)
	match := make([]int, n)
	for i := range match {
		match[i] = -1
	}

	var dfs func(u int, visited []bool) bool
	dfs = func(u int, visited []bool) bool {
		for v := 0; v < n; v++ {
			if !adj[u][v] || visited[v] {
				continue
			}
			visited[v] = true
			if match[v] == -1 || dfs(match[v], visited) {
				match[v] = u
				return true
			}
		}
		return false
	}

	result := 0
	for u := 0; u < n; u++ {
		if dfs(u, make([]bool, n)) {
			result++
		}
	}
	return result
}

func minCharts(stocks [][]int) int {
	n := len(stocks)
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
		for j := range adj[i] {
			if i != j && dominates(stocks[i], stocks[j]) {
				adj[i][j] = true
			}
		}
	}
	return n - maxBipartiteMatching(adj)
}

func parseStocks(input string) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return nil, err
	}
	stocks := make([][]int, n)
	for i := range stocks {
		stocks[i] = make([]int, k)
		for j := range stocks[i] {
			if _, err := fmt.Fscan(reader, &stocks[i][j]); err != nil {
				return nil, err
			}
		}
	}
	return stocks, nil
}

func solveStocks(reader io.Reader, writer io.Writer) error {
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return err
	}
	stocks := make([][]int, n)
	for i := range stocks {
		stocks[i] = make([]int, k)
		for j := range stocks[i] {
			if _, err := fmt.Fscan(reader, &stocks[i][j]); err != nil {
				return err
			}
		}
	}
	fmt.Fprintln(writer, minCharts(stocks))
	return nil
}

func main() {
	if err := solveStocks(bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		return
	}
}
