package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/sat"
)

func TestBudgetAllocationFixtures(t *testing.T) {
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
		got := solveBudgetText(string(input))
		if got != string(expected) {
			t.Fatalf("%s: got %q, want %q", entry.Name(), got, string(expected))
		}
	}
}

func TestBudgetAllocationCustom(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"satisfiable", "2 3\n1 2 3\n0 -1 0\n6\n0\n"},
		{"unsatisfiable", "2 1\n1\n-1\n0\n0\n"},
		{"trivial", "1 1\n0\n0\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clauses, variables := parseBudgetAndBuild(tt.input)
			got, _ := sat.Solve(clauses, variables)
			n, m, a, b := parseBudgetProblem(tt.input)
			want := false
			for mask := 0; mask < 1<<m; mask++ {
				feasible := true
				for i := 0; i < n; i++ {
					sum := 0
					for j := 0; j < m; j++ {
						if a[i][j] != 0 && mask&(1<<j) != 0 {
							sum += a[i][j]
						}
					}
					if sum > b[i] {
						feasible = false
						break
					}
				}
				if feasible {
					want = true
					break
				}
			}
			if got != want {
				t.Fatalf("formula satisfiable = %v, brute force = %v", got, want)
			}
		})
	}
}
