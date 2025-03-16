package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/sat"
)

func TestGSMNetworkFixtures(t *testing.T) {
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
		n, m, edges, err := parseGSM(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		clauses, variables := gsmClauses(n, m, edges)
		gotSat, _ := sat.Solve(clauses, variables)
		wantClauses, wantVars := sat.ParseCNF(string(expected))
		wantSat, _ := sat.Solve(wantClauses, wantVars)
		if gotSat != wantSat {
			t.Fatalf("%s: formula satisfiable = %v, want %v", entry.Name(), gotSat, wantSat)
		}
		if gotSat != threeColorable(n, edges) {
			t.Fatalf("%s: reduction and brute-force colorability disagree", entry.Name())
		}
		if len(clauses) > 5000 || variables > 3000 || len(clauses) == 0 || variables == 0 {
			t.Fatalf("%s: formula limits exceeded: %d clauses, %d variables", entry.Name(), len(clauses), variables)
		}
	}
}

func TestGSMSingleVertex(t *testing.T) {
	clauses, variables := gsmClauses(1, 0, nil)
	ok, _ := sat.Solve(clauses, variables)
	if !ok {
		t.Fatal("one vertex should be 3-colorable")
	}
}

func threeColorable(n int, edges [][2]int) bool {
	colors := make([]int, n)
	var search func(int) bool
	search = func(vertex int) bool {
		if vertex == n {
			return true
		}
		for color := 0; color < 3; color++ {
			valid := true
			for _, edge := range edges {
				u, v := edge[0]-1, edge[1]-1
				other := -1
				if u == vertex {
					other = v
				} else if v == vertex {
					other = u
				}
				if other >= 0 && colors[other] == color {
					valid = false
					break
				}
			}
			if valid {
				colors[vertex] = color
				if search(vertex + 1) {
					return true
				}
				colors[vertex] = 0
			}
		}
		return false
	}
	return search(0)
}
