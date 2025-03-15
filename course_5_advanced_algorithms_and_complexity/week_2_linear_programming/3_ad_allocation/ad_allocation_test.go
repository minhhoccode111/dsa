package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/lp"
)

func TestAdAllocationFixtures(t *testing.T) {
	entries, err := os.ReadDir("test")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".a") {
			continue
		}
		input, err := os.ReadFile(filepath.Join("test", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile(filepath.Join("test", entry.Name()+".a"))
		if err != nil {
			t.Fatal(err)
		}
		n, m, a, b, c, err := parseAdAllocation(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		status, solution := lp.SimplexMax(a, b, c)
		lines := strings.Split(strings.TrimSpace(string(expected)), "\n")
		switch lines[0] {
		case "Infinity":
			if status != lp.SimplexUnbounded {
				t.Fatalf("%s: status = %d, want Infinity", entry.Name(), status)
			}
		case "No solution":
			if status != lp.SimplexInfeasible {
				t.Fatalf("%s: status = %d, want No solution", entry.Name(), status)
			}
		case "Bounded solution":
			if status != lp.SimplexBounded {
				t.Fatalf("%s: status = %d, want bounded", entry.Name(), status)
			}
			fields := strings.Fields(lines[1])
			if len(fields) != m {
				t.Fatalf("%s: expected %d values, got %d", entry.Name(), m, len(fields))
			}
			expectedObjective := 0.0
			for j, field := range fields {
				want, err := strconv.ParseFloat(field, 64)
				if err != nil {
					t.Fatal(err)
				}
				expectedObjective += c[j] * want
				if n >= 0 && math.IsNaN(want) {
					t.Fatalf("bad expected value")
				}
			}
			objective := 0.0
			for j, value := range solution {
				objective += c[j] * value
			}
			if math.Abs(objective-expectedObjective) > 1e-2 {
				t.Fatalf("%s: objective = %v, want %v", entry.Name(), objective, expectedObjective)
			}
			for i, row := range a {
				sum := 0.0
				for j, coefficient := range row {
					sum += coefficient * solution[j]
				}
				if sum > b[i]+1e-3 {
					t.Fatalf("%s: constraint %d violated: %v > %v", entry.Name(), i, sum, b[i])
				}
			}
			for j, value := range solution {
				if value < -1e-6 {
					t.Fatalf("%s: negative variable x[%d] = %v", entry.Name(), j, value)
				}
			}
		default:
			t.Fatalf("%s: unexpected expectation %q", entry.Name(), lines[0])
		}
	}
}
