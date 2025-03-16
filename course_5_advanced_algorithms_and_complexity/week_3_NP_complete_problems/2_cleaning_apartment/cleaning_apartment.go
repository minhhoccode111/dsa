// go run cleaning_apartment.go

// Reduce the Hamiltonian path problem in an apartment to SAT.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func apartmentVar(vertex, position, n int) int {
	return vertex + (position-1)*n
}

func apartmentClauses(n, m int, edges [][2]int) ([][]int, int) {
	adjacency := make([][]int, n)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		adjacency[u-1] = append(adjacency[u-1], v)
		adjacency[v-1] = append(adjacency[v-1], u)
	}

	clauses := make([][]int, 0)
	for vertex := 1; vertex <= n; vertex++ {
		literals := make([]int, n)
		for position := 1; position <= n; position++ {
			literals[position-1] = apartmentVar(vertex, position, n)
		}
		clauses = append(clauses, append([]int(nil), literals...))
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				clauses = append(clauses, []int{-literals[i], -literals[j]})
			}
		}
	}
	for position := 1; position <= n; position++ {
		literals := make([]int, n)
		for vertex := 1; vertex <= n; vertex++ {
			literals[vertex-1] = apartmentVar(vertex, position, n)
		}
		clauses = append(clauses, append([]int(nil), literals...))
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				clauses = append(clauses, []int{-literals[i], -literals[j]})
			}
		}
	}
	for vertex := 0; vertex < n; vertex++ {
		for position := 1; position < n; position++ {
			clause := []int{-apartmentVar(position, vertex+1, n)}
			for _, neighbor := range adjacency[vertex] {
				clause = append(clause, apartmentVar(position+1, neighbor, n))
			}
			clauses = append(clauses, clause)
		}
	}
	return clauses, n * n
}

func parseApartment(input string) (int, int, [][2]int, error) {
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

func formatApartmentCNF(clauses [][]int, variables int) string {
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
	clauses, variables := apartmentClauses(n, m, edges)
	fmt.Print(formatApartmentCNF(clauses, variables))
}
