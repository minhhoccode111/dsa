// go run reschedule_exams.go

// Recolor a graph with three colors so that every vertex changes color and no
// edge is monochromatic, using 2-SAT.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/twosat"
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

func recolorClauses(n int, colors string, edges [][2]int) []twosat.Clause {
	choices := make([]string, n)
	for i := 0; i < n; i++ {
		choices[i] = choicesFor(colors[i])
	}
	clauses := make([]twosat.Clause, 0, 2*n+4*len(edges))
	for i := 0; i < n; i++ {
		first := rescheduleVar(i, 0)
		second := rescheduleVar(i, 1)
		clauses = append(clauses, twosat.Clause{first, second})
		clauses = append(clauses, twosat.Clause{-first, -second})
	}
	for _, edge := range edges {
		u, v := edge[0]-1, edge[1]-1
		if u < 0 || u >= n || v < 0 || v >= n || u == v {
			continue
		}
		for firstChoice := 0; firstChoice < 2; firstChoice++ {
			for secondChoice := 0; secondChoice < 2; secondChoice++ {
				if choices[u][firstChoice] == choices[v][secondChoice] {
					clauses = append(clauses, twosat.Clause{
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
	assignment, ok := twosat.Solve(2*n, clauses)
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
