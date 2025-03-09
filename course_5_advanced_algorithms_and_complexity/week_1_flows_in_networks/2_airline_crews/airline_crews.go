// go run airline_crews.go

// Find a maximum matching between flights and crews in a bipartite graph.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func findMatching(n, m int, adj [][]bool) []int {
	match := make([]int, m)
	for i := range match {
		match[i] = -1
	}

	var dfs func(flight int, visited []bool) bool
	dfs = func(flight int, visited []bool) bool {
		for crew := 0; crew < m; crew++ {
			if !adj[flight][crew] || visited[crew] {
				continue
			}
			visited[crew] = true
			if match[crew] == -1 || dfs(match[crew], visited) {
				match[crew] = flight
				return true
			}
		}
		return false
	}

	for flight := 0; flight < n; flight++ {
		dfs(flight, make([]bool, m))
	}

	matching := make([]int, n)
	for i := range matching {
		matching[i] = -1
	}
	for crew, flight := range match {
		if flight != -1 {
			matching[flight] = crew
		}
	}
	return matching
}

func formatMatching(matching []int) string {
	parts := make([]string, len(matching))
	for i, crew := range matching {
		if crew == -1 {
			parts[i] = "-1"
		} else {
			parts[i] = fmt.Sprint(crew + 1)
		}
	}
	return strings.Join(parts, " ")
}

func parseCrewInput(input string) (int, int, [][]bool, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, err
	}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, m)
		for j := range adj[i] {
			var value int
			if _, err := fmt.Fscan(reader, &value); err != nil {
				return 0, 0, nil, err
			}
			adj[i][j] = value != 0
		}
	}
	return n, m, adj, nil
}

func solveCrews(reader io.Reader, writer io.Writer) error {
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return err
	}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, m)
		for j := range adj[i] {
			var value int
			if _, err := fmt.Fscan(reader, &value); err != nil {
				return err
			}
			adj[i][j] = value != 0
		}
	}
	fmt.Fprintln(writer, formatMatching(findMatching(n, m, adj)))
	return nil
}

func main() {
	if err := solveCrews(bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		return
	}
}
