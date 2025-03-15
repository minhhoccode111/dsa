package lp

import (
	"math"
	"testing"
)

func TestSolveLinear(t *testing.T) {
	solution, ok := SolveLinear([][]float64{{1, 1}, {2, 3}}, []float64{3, 7})
	if !ok {
		t.Fatal("expected a solution")
	}
	if math.Abs(solution[0]-2) > 1e-9 || math.Abs(solution[1]-1) > 1e-9 {
		t.Fatalf("solution = %v, want [2 1]", solution)
	}
	if _, ok := SolveLinear([][]float64{{1, 2}, {2, 4}}, []float64{1, 2}); ok {
		t.Fatal("singular matrix reported as solvable")
	}
}

func TestSimplexMax(t *testing.T) {
	tests := []struct {
		name   string
		a      [][]float64
		b, c   []float64
		status SimplexStatus
		value  float64
	}{
		{"bounded", [][]float64{{1, 0}, {0, 1}}, []float64{2, 3}, []float64{1, 2}, SimplexBounded, 8},
		{"infeasible", [][]float64{{1}, {-1}}, []float64{1, -2}, []float64{1}, SimplexInfeasible, 0},
		{"unbounded", [][]float64{{-1}}, []float64{-1}, []float64{1}, SimplexUnbounded, 0},
		{"negative rhs", [][]float64{{1}, {-1}}, []float64{2, -1}, []float64{1}, SimplexBounded, 2},
		{"unbounded scaling", [][]float64{{1, 5, -9}, {-8, -10, 4}}, []float64{4, 21}, []float64{-1, 8, 8}, SimplexUnbounded, 0},
		{
			"beale cycling example",
			[][]float64{
				{0.25, -8, -1, 9},
				{0.5, -12, -0.5, 3},
				{0, 0, 1, 0},
			},
			[]float64{0, 0, 1},
			[]float64{0.75, -20, 0.5, -6},
			SimplexBounded,
			1.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, solution := SimplexMax(tt.a, tt.b, tt.c)
			if status != tt.status {
				t.Fatalf("status = %d, want %d", status, tt.status)
			}
			if status == SimplexBounded {
				value := 0.0
				for i, coefficient := range tt.c {
					value += coefficient * solution[i]
				}
				if math.Abs(value-tt.value) > 1e-6 {
					t.Fatalf("objective = %v, want %v", value, tt.value)
				}
			}
		})
	}
}
