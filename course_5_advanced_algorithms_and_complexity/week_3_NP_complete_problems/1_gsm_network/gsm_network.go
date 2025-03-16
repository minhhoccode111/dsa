// go run gsm_network.go

// Reduce 3-coloring of a GSM network to SAT.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func gsmVar(vertex, color, n int) int {
	return vertex + (color-1)*n
}

func gsmClauses(n, m int, edges [][2]int) ([][]int, int) {
	clauses := make([][]int, 0, n*4+m*3)
	for vertex := 1; vertex <= n; vertex++ {
		literals := []int{
			gsmVar(vertex, 1, n),
			gsmVar(vertex, 2, n),
			gsmVar(vertex, 3, n),
		}
		clauses = append(clauses, append([]int(nil), literals...))
		for i := 0; i < len(literals); i++ {
			for j := i + 1; j < len(literals); j++ {
				clauses = append(clauses, []int{-literals[i], -literals[j]})
			}
		}
	}
	for _, edge := range edges {
		for color := 1; color <= 3; color++ {
			clauses = append(clauses, []int{
				-gsmVar(edge[0], color, n),
				-gsmVar(edge[1], color, n),
			})
		}
	}
	return clauses, n * 3
}

func parseGSM(input string) (int, int, [][2]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, err
	}
	edges := make([][2]int, m)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return 0, 0, nil, err
		}
	}
	return n, m, edges, nil
}

func formatCNF(clauses [][]int, variables int) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d %d\n", len(clauses), variables)
	for _, clause := range clauses {
		for _, literal := range clause {
			fmt.Fprintf(&builder, "%d ", literal)
		}
		builder.WriteString("0\n")
	}
	return builder.String()
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
	}
	clauses, variables := gsmClauses(n, m, edges)
	fmt.Print(formatCNF(clauses, variables))
}
