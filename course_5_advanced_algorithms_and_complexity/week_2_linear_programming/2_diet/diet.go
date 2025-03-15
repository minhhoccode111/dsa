// go run diet.go

// Solve the optimal diet linear program with the simplex method.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/lp"
)

const (
	dietNoSolution = -1
	dietBounded    = 0
	dietInfinity   = 1
)

func parseDiet(input string) (int, int, [][]float64, []float64, []float64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, 0, nil, nil, nil, err
	}
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return 0, 0, nil, nil, nil, err
			}
		}
	}
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return 0, 0, nil, nil, nil, err
		}
	}
	c := make([]float64, m)
	for j := 0; j < m; j++ {
		if _, err := fmt.Fscan(reader, &c[j]); err != nil {
			return 0, 0, nil, nil, nil, err
		}
	}
	return n, m, a, b, c, nil
}

func solveDiet(n, m int, a [][]float64, b, c []float64) (int, []float64) {
	status, solution := lp.SimplexMax(a, b, c)
	switch status {
	case lp.SimplexInfeasible:
		return dietNoSolution, nil
	case lp.SimplexUnbounded:
		return dietInfinity, nil
	default:
		return dietBounded, solution
	}
}

func formatDiet(values []float64) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprintf("%.18f", value)
	}
	return strings.Join(parts, " ")
}

func solveDietIO(reader io.Reader, writer io.Writer) error {
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return err
	}
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return err
			}
		}
	}
	b := make([]float64, n)
	for i := range b {
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return err
		}
	}
	c := make([]float64, m)
	for j := range c {
		if _, err := fmt.Fscan(reader, &c[j]); err != nil {
			return err
		}
	}

	status, solution := solveDiet(n, m, a, b, c)
	switch status {
	case dietNoSolution:
		fmt.Fprintln(writer, "No solution")
	case dietInfinity:
		fmt.Fprintln(writer, "Infinity")
	default:
		fmt.Fprintln(writer, "Bounded solution")
		fmt.Fprintln(writer, formatDiet(solution))
	}
	return nil
}

func main() {
	if err := solveDietIO(bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		return
	}
}
