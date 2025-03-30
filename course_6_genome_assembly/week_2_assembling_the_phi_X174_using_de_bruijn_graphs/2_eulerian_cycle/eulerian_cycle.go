// go run eulerian_cycle.go

// Find an Eulerian cycle in a directed graph.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func findEulerian(n int, edges [][2]int) ([]int, bool) {
	return assembly.FindEulerianCycle(n, edges)
}

func formatEulerian(n int, edges [][2]int) (string, bool) {
	path, ok := findEulerian(n, edges)
	if !ok {
		return "0", false
	}
	parts := make([]string, 0, len(path))
	for i := 0; i+1 < len(path); i++ {
		parts = append(parts, strconv.Itoa(path[i]+1))
	}
	return "1\n" + strings.Join(parts, " "), true
}

func parseEulerian(input string) (int, [][2]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, nil, err
	}
	edges := make([][2]int, m)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return 0, nil, err
		}
		edges[i][0]--
		edges[i][1]--
	}
	return n, edges, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([][2]int, m)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return
		}
		edges[i][0]--
		edges[i][1]--
	}
	result, _ := formatEulerian(n, edges)
	fmt.Println(result)
}
