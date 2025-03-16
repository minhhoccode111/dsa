// go run budget_allocation.go

// Reduce binary integer linear programming to SAT.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func budgetClauses(a [][]int, b []int) [][]int {
	clauses := make([][]int, 0)
	for i, row := range a {
		ids := make([]int, 0, 3)
		coefficients := make([]int, 0, 3)
		for j, value := range row {
			if value != 0 {
				ids = append(ids, j+1)
				coefficients = append(coefficients, value)
				if len(coefficients) == 3 {
					break
				}
			}
		}
		if len(ids) == 0 {
			continue
		}
		for mask := 0; mask < 1<<len(ids); mask++ {
			sum := 0
			for k, coefficient := range coefficients {
				if mask&(1<<k) != 0 {
					sum += coefficient
				}
			}
			if sum <= b[i] {
				continue
			}
			clause := make([]int, len(ids))
			for k, id := range ids {
				if mask&(1<<k) == 0 {
					clause[k] = id
				} else {
					clause[k] = -id
				}
			}
			clauses = append(clauses, clause)
		}
	}
	if len(clauses) == 0 {
		clauses = append(clauses, []int{1, -1})
	}
	return clauses
}

func parseBudgetProblem(input string) (int, int, [][]int, []int) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, nil
	}
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, m)
		for j := range a[i] {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return 0, 0, nil, nil
			}
		}
	}
	b := make([]int, n)
	for i := range b {
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return 0, 0, nil, nil
		}
	}
	return n, m, a, b
}

func parseBudgetAndBuild(input string) ([][]int, int) {
	n, m, a, b := parseBudgetProblem(input)
	if n == 0 && m == 0 {
		return nil, 0
	}
	return budgetClauses(a, b), m
}

func solveBudgetText(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, m)
		for j := range a[i] {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return ""
			}
		}
	}
	b := make([]int, n)
	for i := range b {
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return ""
		}
	}
	clauses := budgetClauses(a, b)
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d %d\n", len(clauses), m)
	for _, clause := range clauses {
		for _, literal := range clause {
			fmt.Fprintf(&builder, "%d ", literal)
		}
		builder.WriteString("0\n")
	}
	return builder.String()
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	fmt.Print(solveBudgetText(string(data)))
}
