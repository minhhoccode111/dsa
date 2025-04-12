// go run reschedule_exams.go

// Recolor a graph with three colors so that every vertex changes color and no
// edge is monochromatic, using 2-SAT.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func choicesFor(color byte) string {
	switch color {
	case 'R':
		return "GB"
	case 'G':
		return "RB"
	default:
		return "RG"
	}
}

func rescheduleVar(vertex, choice int) int {
	return 2*vertex + choice + 1
}

func recolorClauses(n int, colors string, edges [][2]int) []Clause {
	choices := make([]string, n)
	for i := 0; i < n; i++ {
		choices[i] = choicesFor(colors[i])
	}
	clauses := make([]Clause, 0, 2*n+4*len(edges))
	for i := 0; i < n; i++ {
		first := rescheduleVar(i, 0)
		second := rescheduleVar(i, 1)
		clauses = append(clauses, Clause{first, second})
		clauses = append(clauses, Clause{-first, -second})
	}
	for _, edge := range edges {
		u, v := edge[0]-1, edge[1]-1
		if u < 0 || u >= n || v < 0 || v >= n || u == v {
			continue
		}
		for firstChoice := 0; firstChoice < 2; firstChoice++ {
			for secondChoice := 0; secondChoice < 2; secondChoice++ {
				if choices[u][firstChoice] == choices[v][secondChoice] {
					clauses = append(clauses, Clause{
						-rescheduleVar(u, firstChoice),
						-rescheduleVar(v, secondChoice),
					})
				}
			}
		}
	}
	return clauses
}

func recolor(n, m int, colors string, edges [][2]int, loop bool) (string, bool) {
	if loop || n == 0 {
		return "", false
	}
	choices := make([]string, n)
	for i := 0; i < n; i++ {
		choices[i] = choicesFor(colors[i])
	}
	clauses := recolorClauses(n, colors, edges)
	assignment, ok := Solve(2*n, clauses)
	if !ok {
		return "", false
	}
	var builder strings.Builder
	for i := 0; i < n; i++ {
		choice := 0
		if !assignment[rescheduleVar(i, 0)] {
			choice = 1
		}
		builder.WriteByte(choices[i][choice])
	}
	return builder.String(), true
}

func parseReschedule(input string) (int, int, string, [][2]int, bool, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, "", nil, false, err
	}
	var colors string
	if _, err := fmt.Fscan(reader, &colors); err != nil {
		return 0, 0, "", nil, false, err
	}
	edges := make([][2]int, m)
	loop := false
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return 0, 0, "", nil, false, err
		}
		if edges[i][0] == edges[i][1] {
			loop = true
		}
	}
	return n, m, colors, edges, loop, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var colors string
	if _, err := fmt.Fscan(reader, &colors); err != nil {
		return
	}
	edges := make([][2]int, m)
	loop := false
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return
		}
		if edges[i][0] == edges[i][1] {
			loop = true
		}
	}
	result, ok := recolor(n, m, colors, edges, loop)
	if !ok {
		fmt.Println("Impossible")
		return
	}
	fmt.Println(result)
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
