// Package lp provides Gaussian elimination and simplex helpers used by the
// linear-programming assignments.
package lp

import "math"

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
	status := simplexCore(tableau, objective, basis)
	if status == SimplexUnbounded {
		return SimplexInfeasible, nil
	}
	if objective[total] < -1e-7 {
		return SimplexInfeasible, nil
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

	cost = make([]float64, cols+rows)
	copy(cost, c)
	objective = buildObjective(tableau, basis, cost)
	status = simplexCore(tableau, objective, basis)
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

func simplexCore(tableau [][]float64, objective []float64, basis []int) SimplexStatus {
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
			ratio := tableau[i][cols] / coefficient
			if ratio < bestRatio-1e-9 || (math.Abs(ratio-bestRatio) <= 1e-9 && (leaving < 0 || basis[i] < basis[leaving])) {
				bestRatio = ratio
				leaving = i
			}
		}
		if leaving < 0 {
			return SimplexUnbounded
		}
		pivot(tableau, objective, basis, leaving, entering)
	}
	return SimplexUnbounded
}

func pivot(tableau [][]float64, objective []float64, basis []int, row, column int) {
	cols := len(tableau[0]) - 1
	value := tableau[row][column]
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
	}
	basis[row] = column
}
