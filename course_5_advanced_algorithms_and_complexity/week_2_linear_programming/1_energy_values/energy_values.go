// go run energy_values.go

// Infer ingredient energy values with Gaussian elimination.

package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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
	return SolveLinear(a, b)
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

// The helpers below are inlined from dsa/course_5_advanced_algorithms_and_complexity/internal/lp
// so that this file can be submitted to Coursera as a single Go file.

const epsilon = 1e-9

// SolveLinear solves a square linear system using Gauss-Jordan elimination
// with partial pivoting. It returns false when the matrix is singular.
func SolveLinear(a [][]float64, b []float64) ([]float64, bool) {
	n := len(a)
	if n == 0 {
		return []float64{}, true
	}
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n+1)
		copy(matrix[i], a[i])
		matrix[i][n] = b[i]
	}

	for col := 0; col < n; col++ {
		pivot := col
		for row := col + 1; row < n; row++ {
			if math.Abs(matrix[row][col]) > math.Abs(matrix[pivot][col]) {
				pivot = row
			}
		}
		if math.Abs(matrix[pivot][col]) < epsilon {
			return nil, false
		}
		matrix[col], matrix[pivot] = matrix[pivot], matrix[col]

		scale := matrix[col][col]
		for j := col; j <= n; j++ {
			matrix[col][j] /= scale
		}
		for row := 0; row < n; row++ {
			if row == col {
				continue
			}
			factor := matrix[row][col]
			if math.Abs(factor) < epsilon {
				continue
			}
			for j := col; j <= n; j++ {
				matrix[row][j] -= factor * matrix[col][j]
			}
		}
	}

	solution := make([]float64, n)
	for i := 0; i < n; i++ {
		solution[i] = matrix[i][n]
	}
	return solution, true
}

// SimplexStatus describes the outcome of a linear program.
type SimplexStatus int

const (
	// SimplexBounded means an optimal solution was found.
	SimplexBounded SimplexStatus = iota
	// SimplexInfeasible means no point satisfies all constraints.
	SimplexInfeasible
	// SimplexUnbounded means the objective can grow without bound.
	SimplexUnbounded
)

// SimplexMax maximizes c^T x subject to A x <= b and x >= 0.
func SimplexMax(a [][]float64, b, c []float64) (SimplexStatus, []float64) {
	rows := len(a)
	cols := len(c)
	// Scale each constraint row to a unit maximum coefficient. This improves
	// conditioning and keeps rounding errors from driving right-hand sides
	// negative during pivoting.
	if rows > 0 {
		scaledA := make([][]float64, rows)
		scaledB := make([]float64, rows)
		for i := 0; i < rows; i++ {
			scale := 1.0
			for _, value := range a[i] {
				if math.Abs(value) > scale {
					scale = math.Abs(value)
				}
			}
			scaledA[i] = make([]float64, len(a[i]))
			for j, value := range a[i] {
				scaledA[i][j] = value / scale
			}
			scaledB[i] = b[i] / scale
		}
		a = scaledA
		b = scaledB
	}
	if rows == 0 {
		for _, value := range c {
			if value > epsilon {
				return SimplexUnbounded, nil
			}
		}
		return SimplexBounded, make([]float64, cols)
	}

	artificialCount := 0
	for _, value := range b {
		if value < 0 {
			artificialCount++
		}
	}
	total := cols + rows + artificialCount
	tableau := make([][]float64, rows+1)
	for i := range tableau {
		tableau[i] = make([]float64, total+1)
	}
	basis := make([]int, rows)

	nextArtificial := cols + rows
	for i := 0; i < rows; i++ {
		if b[i] >= 0 {
			copy(tableau[i], a[i])
			tableau[i][cols+i] = 1
			tableau[i][total] = b[i]
			basis[i] = cols + i
		} else {
			for j := 0; j < cols; j++ {
				tableau[i][j] = -a[i][j]
			}
			tableau[i][cols+i] = -1
			tableau[i][nextArtificial] = 1
			tableau[i][total] = -b[i]
			basis[i] = nextArtificial
			nextArtificial++
		}
	}

	cost := make([]float64, total)
	for j := cols + rows; j < total; j++ {
		cost[j] = -1
	}
	objective := buildObjective(tableau, basis, cost)
	status := simplexCore(tableau, objective, basis, false)
	if status == SimplexUnbounded {
		return SimplexInfeasible, nil
	}
	// Phase 1 is feasible only when every artificial variable is zero. Check
	// the basic artificial values instead of the objective constant, which
	// accumulates rounding error during pivoting.
	for i := 0; i < rows; i++ {
		if basis[i] < cols+rows {
			continue
		}
		if tableau[i][total] > 1e-5 {
			return SimplexInfeasible, nil
		}
		tableau[i][total] = 0
	}
	// Pivot any artificial variable that is still basic out of the basis.
	for i := 0; i < len(tableau)-1; i++ {
		if basis[i] < cols+rows {
			continue
		}
		entering := -1
		for j := 0; j < cols+rows; j++ {
			if math.Abs(tableau[i][j]) > epsilon {
				entering = j
				break
			}
		}
		if entering >= 0 {
			pivot(tableau, objective, basis, i, entering)
		}
	}

	// Drop artificial columns. Rows whose artificial variable could not be
	// pivoted out are redundant and can be removed.
	newRows := make([][]float64, 0, rows)
	newBasis := make([]int, 0, rows)
	for i := 0; i < rows; i++ {
		if basis[i] >= cols+rows {
			continue
		}
		row := make([]float64, cols+rows+1)
		copy(row, tableau[i][:cols+rows])
		row[cols+rows] = tableau[i][total]
		newRows = append(newRows, row)
		newBasis = append(newBasis, basis[i])
	}
	if len(newRows) == 0 {
		return SimplexBounded, make([]float64, cols)
	}

	tableau = make([][]float64, len(newRows)+1)
	for i := range tableau {
		tableau[i] = make([]float64, cols+rows+1)
	}
	for i := range newRows {
		copy(tableau[i], newRows[i])
	}
	basis = newBasis

	// Pivoting can leave tiny negative right-hand sides. Clamp them to zero
	// so that phase 2 starts from a feasible basis.
	for i := 0; i < len(tableau)-1; i++ {
		if tableau[i][cols+rows] < 0 && tableau[i][cols+rows] > -1e-7 {
			tableau[i][cols+rows] = 0
		}
	}

	cost = make([]float64, cols+rows)
	copy(cost, c)
	objective = buildObjective(tableau, basis, cost)
	status = simplexCore(tableau, objective, basis, true)
	if status == SimplexUnbounded {
		return SimplexUnbounded, nil
	}

	solution := make([]float64, cols)
	for i, variable := range basis {
		if variable < cols {
			solution[variable] = tableau[i][cols+rows]
		}
	}
	return SimplexBounded, solution
}

func buildObjective(tableau [][]float64, basis []int, cost []float64) []float64 {
	rows := len(tableau) - 1
	cols := len(tableau[0]) - 1
	objective := tableau[rows]
	for j := range objective {
		objective[j] = 0
	}
	for j := 0; j < cols; j++ {
		objective[j] = -cost[j]
	}
	for i := 0; i < rows; i++ {
		coefficient := cost[basis[i]]
		if math.Abs(coefficient) < epsilon {
			continue
		}
		for j := 0; j < cols; j++ {
			objective[j] += coefficient * tableau[i][j]
		}
		objective[cols] += coefficient * tableau[i][cols]
	}
	return objective
}

func simplexCore(tableau [][]float64, objective []float64, basis []int, unboundedOK bool) SimplexStatus {
	rows := len(tableau) - 1
	cols := len(tableau[0]) - 1
	for iteration := 0; iteration < 200000; iteration++ {
		// Choose the smallest index with negative reduced cost. Combined with
		// the smallest-basis-variable ratio tie-break below, this is Bland's
		// rule, which prevents cycling on degenerate problems.
		entering := -1
		for j := 0; j < cols; j++ {
			if objective[j] < -1e-9 {
				entering = j
				break
			}
		}
		if entering < 0 {
			return SimplexBounded
		}

		leaving := -1
		bestRatio := math.Inf(1)
		for i := 0; i < rows; i++ {
			coefficient := tableau[i][entering]
			if coefficient <= 1e-9 {
				continue
			}
			rhs := tableau[i][cols]
			if rhs < 0 && rhs > -1e-7 {
				rhs = 0
			}
			ratio := rhs / coefficient
			if ratio < bestRatio-1e-9 || (math.Abs(ratio-bestRatio) <= 1e-9 && (leaving < 0 || basis[i] < basis[leaving])) {
				bestRatio = ratio
				leaving = i
			}
		}
		if leaving < 0 {
			if unboundedOK {
				return SimplexUnbounded
			}
			// Phase 1 cannot be unbounded: its objective is the negated sum
			// of non-negative artificial variables. Treat the column as a
			// numerical artefact and ignore its reduced cost.
			objective[entering] = 0
			continue
		}
		pivot(tableau, objective, basis, leaving, entering)
	}
	return SimplexUnbounded
}

func pivot(tableau [][]float64, objective []float64, basis []int, row, column int) {
	cols := len(tableau[0]) - 1
	value := tableau[row][column]
	if tableau[row][cols] < 0 && tableau[row][cols] > -1e-7 {
		tableau[row][cols] = 0
	}
	for j := 0; j < cols+1; j++ {
		tableau[row][j] /= value
	}
	for i := 0; i < len(tableau); i++ {
		if i == row {
			continue
		}
		factor := tableau[i][column]
		if math.Abs(factor) < 1e-12 {
			tableau[i][column] = 0
			continue
		}
		for j := 0; j < cols+1; j++ {
			tableau[i][j] -= factor * tableau[row][j]
		}
		if i < len(tableau)-1 && tableau[i][cols] < 0 && tableau[i][cols] > -1e-7 {
			tableau[i][cols] = 0
		}
	}
	basis[row] = column
}
