// go run ad_allocation.go

// Maximize advertising revenue with the simplex method.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/lp"
)

func parseAdAllocation(input string) (int, int, [][]float64, []float64, []float64, error) {
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

func solveAdAllocation(reader io.Reader, writer io.Writer) error {
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

	status, solution := lp.SimplexMax(a, b, c)
	switch status {
	case lp.SimplexInfeasible:
		fmt.Fprintln(writer, "No solution")
	case lp.SimplexUnbounded:
		fmt.Fprintln(writer, "Infinity")
	default:
		parts := make([]string, len(solution))
		for i, value := range solution {
			parts[i] = fmt.Sprintf("%.18f", value)
		}
		fmt.Fprintln(writer, "Bounded solution")
		fmt.Fprintln(writer, strings.Join(parts, " "))
	}
	return nil
}

func main() {
	if err := solveAdAllocation(bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		return
	}
}
