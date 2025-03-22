// go run circuit_design.go

// Solve a 2-SAT instance describing wire bending directions.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/twosat"
)

func parseCircuit(input string) (int, int, []twosat.Clause, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, err
	}
	clauses := make([]twosat.Clause, m)
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
	clauses := make([]twosat.Clause, m)
	for i := range clauses {
		if _, err := fmt.Fscan(reader, &clauses[i][0], &clauses[i][1]); err != nil {
			return
		}
	}
	assignment, ok := twosat.Solve(n, clauses)
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
