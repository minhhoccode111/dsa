package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDietFixtures(t *testing.T) {
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
		n, m, a, b, c, err := parseDiet(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		status, solution := solveDiet(n, m, a, b, c)
		lines := strings.Split(strings.TrimSpace(string(expected)), "\n")
		switch lines[0] {
		case "No solution":
			if status != dietNoSolution {
				t.Fatalf("%s: status = %d, want No solution", entry.Name(), status)
			}
		case "Infinity":
			if status != dietInfinity {
				t.Fatalf("%s: status = %d, want Infinity", entry.Name(), status)
			}
		case "Bounded solution":
			if status != dietBounded {
				t.Fatalf("%s: status = %d, want Bounded solution", entry.Name(), status)
			}
			fields := strings.Fields(lines[1])
			if len(fields) != m {
				t.Fatalf("%s: expected %d values, got %d", entry.Name(), m, len(fields))
			}
			for i, field := range fields {
				want, err := strconv.ParseFloat(field, 64)
				if err != nil {
					t.Fatal(err)
				}
				if math.Abs(solution[i]-want) > 1e-5 {
					t.Fatalf("%s: x[%d] = %v, want %v", entry.Name(), i, solution[i], want)
				}
			}
		default:
			t.Fatalf("%s: unexpected expectation %q", entry.Name(), lines[0])
		}
	}
}

func TestDietSimple(t *testing.T) {
	status, solution := solveDiet(1, 1, [][]float64{{26}}, []float64{4362}, []float64{27})
	if status != dietBounded {
		t.Fatalf("status = %d, want bounded", status)
	}
	if math.Abs(solution[0]-4362.0/26.0) > 1e-6 {
		t.Fatalf("solution = %v", solution)
	}
}
