// go run circuit_design.go

// Solve a 2-SAT instance describing wire bending directions.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseCircuit(input string) (int, int, []Clause, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, err
	}
	clauses := make([]Clause, m)
	for i := range clauses {
		if _, err := fmt.Fscan(reader, &clauses[i][0], &clauses[i][1]); err != nil {
			return 0, 0, nil, err
		}
	}
	return n, m, clauses, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	clauses := make([]Clause, m)
	for i := range clauses {
		if _, err := fmt.Fscan(reader, &clauses[i][0], &clauses[i][1]); err != nil {
			return
		}
	}
	assignment, ok := Solve(n, clauses)
	if !ok {
		fmt.Println("UNSATISFIABLE")
		return
	}
	fmt.Println("SATISFIABLE")
	parts := make([]string, n)
	for i := 1; i <= n; i++ {
		if assignment[i] {
			parts[i-1] = fmt.Sprint(i)
		} else {
			parts[i-1] = fmt.Sprintf("-%d", i)
		}
	}
	fmt.Println(strings.Join(parts, " "))
}

// The helpers below are inlined from dsa/course_5_advanced_algorithms_and_complexity/internal/twosat
// so that this file can be submitted to Coursera as a single Go file.

// Clause is a pair of literals. Variables are numbered from 1 through n;
// negative literals denote negations.
type Clause [2]int

// Solve returns a satisfying assignment (indexed from 1) and true when the
// formula is satisfiable.
func Solve(n int, clauses []Clause) ([]bool, bool) {
	graph := make([][]int, 2*n)
	reverse := make([][]int, 2*n)
	for _, clause := range clauses {
		a, b := clause[0], clause[1]
		graph[node(-a)] = append(graph[node(-a)], node(b))
		graph[node(-b)] = append(graph[node(-b)], node(a))
		reverse[node(b)] = append(reverse[node(b)], node(-a))
		reverse[node(a)] = append(reverse[node(a)], node(-b))
	}

	order := make([]int, 0, 2*n)
	visited := make([]bool, 2*n)
	for start := 0; start < 2*n; start++ {
		if visited[start] {
			continue
		}
		type frame struct {
			vertex int
			index  int
		}
		visited[start] = true
		stack := []frame{{vertex: start}}
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			if top.index < len(graph[top.vertex]) {
				next := graph[top.vertex][top.index]
				top.index++
				if !visited[next] {
					visited[next] = true
					stack = append(stack, frame{vertex: next})
				}
				continue
			}
			order = append(order, top.vertex)
			stack = stack[:len(stack)-1]
		}
	}

	component := make([]int, 2*n)
	for i := range component {
		component[i] = -1
	}
	nextComponent := 0
	for i := len(order) - 1; i >= 0; i-- {
		start := order[i]
		if component[start] != -1 {
			continue
		}
		component[start] = nextComponent
		stack := []int{start}
		for len(stack) > 0 {
			vertex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, next := range reverse[vertex] {
				if component[next] == -1 {
					component[next] = nextComponent
					stack = append(stack, next)
				}
			}
		}
		nextComponent++
	}

	assignment := make([]bool, n+1)
	for variable := 1; variable <= n; variable++ {
		positive := node(variable)
		negative := node(-variable)
		if component[positive] == component[negative] {
			return nil, false
		}
		assignment[variable] = component[positive] > component[negative]
	}
	return assignment, true
}

func node(literal int) int {
	variable := literal
	if variable < 0 {
		variable = -variable
	}
	base := 2 * (variable - 1)
	if literal > 0 {
		return base
	}
	return base + 1
}
