// go run energy_values.go

// Infer ingredient energy values with Gaussian elimination.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/lp"
)

func parseEquation(input string) (int, [][]float64, []float64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, nil, err
	}
	a := make([][]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return 0, nil, nil, err
			}
		}
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return 0, nil, nil, err
		}
	}
	return n, a, b, nil
}

func solveEquation(n int, a [][]float64, b []float64) ([]float64, bool) {
	if n == 0 {
		return nil, true
	}
	return lp.SolveLinear(a, b)
}

func formatValues(values []float64) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprintf("%.6f", value)
	}
	return strings.Join(parts, " ")
}

func solveEnergy(reader io.Reader, writer io.Writer) error {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return err
	}
	a := make([][]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &a[i][j]); err != nil {
				return err
			}
		}
		if _, err := fmt.Fscan(reader, &b[i]); err != nil {
			return err
		}
	}
	if n == 0 {
		return nil
	}
	values, ok := solveEquation(n, a, b)
	if !ok {
		return nil
	}
	fmt.Fprintln(writer, formatValues(values))
	return nil
}

func main() {
	if err := solveEnergy(bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		return
	}
}
