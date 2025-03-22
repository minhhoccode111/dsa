package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/twosat"
)

func TestCircuitDesignFixtures(t *testing.T) {
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
		n, m, clauses, err := parseCircuit(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		lines := strings.Split(strings.TrimSpace(string(expected)), "\n")
		assignment, ok := twosat.Solve(n, clauses)
		if lines[0] == "UNSATISFIABLE" {
			if ok {
				t.Fatalf("%s: got SATISFIABLE, want UNSATISFIABLE", entry.Name())
			}
			if bruteCircuitSat(n, clauses) {
				t.Fatalf("%s: brute force says satisfiable", entry.Name())
			}
			continue
		}
		if !ok {
			t.Fatalf("%s: got UNSATISFIABLE, want SATISFIABLE", entry.Name())
		}
		if !circuitAssignmentValid(n, m, clauses, assignment) {
			t.Fatalf("%s: solver assignment invalid: %v", entry.Name(), assignment)
		}
		if !bruteCircuitSat(n, clauses) {
			t.Fatalf("%s: brute force says unsatisfiable", entry.Name())
		}
	}
}

func TestCircuitDesignSimple(t *testing.T) {
	clauses := []twosat.Clause{{1, 1}}
	assignment, ok := twosat.Solve(1, clauses)
	if !ok || !assignment[1] {
		t.Fatalf("got %v %v, want x1=true", assignment, ok)
	}
	clauses = []twosat.Clause{{-1, -1}}
	assignment, ok = twosat.Solve(1, clauses)
	if !ok || assignment[1] {
		t.Fatalf("got %v %v, want x1=false", assignment, ok)
	}
	clauses = []twosat.Clause{{1, 1}, {-1, -1}}
	if _, ok := twosat.Solve(1, clauses); ok {
		t.Fatal("expected unsatisfiable")
	}
}

func circuitAssignmentValid(n, m int, clauses []twosat.Clause, assignment []bool) bool {
	if len(assignment) != n+1 {
		return false
	}
	_ = m
	for _, clause := range clauses {
		satisfied := false
		for _, literal := range clause {
			value := false
			if literal > 0 {
				value = assignment[literal]
			} else {
				value = !assignment[-literal]
			}
			if value {
				satisfied = true
			}
		}
		if !satisfied {
			return false
		}
	}
	return true
}

func bruteCircuitSat(n int, clauses []twosat.Clause) bool {
	for mask := 0; mask < 1<<n; mask++ {
		assignment := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			assignment[i] = mask&(1<<(i-1)) != 0
		}
		if circuitAssignmentValid(n, len(clauses), clauses, assignment) {
			return true
		}
	}
	return false
}
