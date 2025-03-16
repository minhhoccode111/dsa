package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/sat"
)

func TestCleaningApartmentFixtures(t *testing.T) {
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
		n, m, edges, err := parseApartment(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		clauses, variables := apartmentClauses(n, m, edges)
		gotSat, _ := sat.Solve(clauses, variables)
		wantClauses, wantVars := sat.ParseCNF(string(expected))
		wantSat, _ := sat.Solve(wantClauses, wantVars)
		if gotSat != wantSat {
			t.Fatalf("%s: formula satisfiable = %v, want %v", entry.Name(), gotSat, wantSat)
		}
		if gotSat != hasHamiltonianPath(n, edges) {
			t.Fatalf("%s: reduction and brute-force Hamiltonian path disagree", entry.Name())
		}
		if len(clauses) == 0 || variables != n*n {
			t.Fatalf("%s: unexpected formula shape", entry.Name())
		}
	}
}

func TestCleaningApartmentTrivial(t *testing.T) {
	clauses, variables := apartmentClauses(1, 0, nil)
	ok, _ := sat.Solve(clauses, variables)
	if !ok {
		t.Fatal("single room should have a Hamiltonian path")
	}
}

func hasHamiltonianPath(n int, edges [][2]int) bool {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for _, edge := range edges {
		adj[edge[0]-1][edge[1]-1] = true
		adj[edge[1]-1][edge[0]-1] = true
	}
	used := make([]bool, n)
	var search func(vertex, depth int) bool
	search = func(vertex, depth int) bool {
		if depth == n {
			return true
		}
		for next := 0; next < n; next++ {
			if !used[next] && adj[vertex][next] {
				used[next] = true
				if search(next, depth+1) {
					return true
				}
				used[next] = false
			}
		}
		return false
	}
	for start := 0; start < n; start++ {
		used[start] = true
		if search(start, 1) {
			return true
		}
		used[start] = false
	}
	return false
}
