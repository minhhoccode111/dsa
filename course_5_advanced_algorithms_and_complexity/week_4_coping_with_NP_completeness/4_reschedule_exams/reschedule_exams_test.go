package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsa/course_5_advanced_algorithms_and_complexity/internal/twosat"
)

func TestRescheduleExamsFixtures(t *testing.T) {
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
		n, m, colors, edges, loop, err := parseReschedule(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got, ok := recolor(n, m, colors, edges, loop)
		want := strings.TrimSpace(string(expected))
		if want == "Impossible" {
			if ok {
				t.Fatalf("%s: got %q, want Impossible", entry.Name(), got)
			}
			continue
		}
		if !ok {
			t.Fatalf("%s: got Impossible, want %q", entry.Name(), want)
		}
		if len(got) != n {
			t.Fatalf("%s: recoloring length = %d, want %d", entry.Name(), len(got), n)
		}
		for i := 0; i < n; i++ {
			if got[i] == colors[i] {
				t.Fatalf("%s: vertex %d kept color %c", entry.Name(), i+1, colors[i])
			}
		}
		for _, edge := range edges {
			u, v := edge[0]-1, edge[1]-1
			if u >= 0 && u < n && v >= 0 && v < n && got[u] == got[v] {
				t.Fatalf("%s: edge %d-%d shares color %c", entry.Name(), u+1, v+1, got[u])
			}
		}
	}
}

func TestRescheduleExamsLoops(t *testing.T) {
	got, ok := recolor(2, 1, "RG", [][2]int{{1, 1}}, true)
	if ok || got != "" {
		t.Fatalf("loop should be impossible, got %q %v", got, ok)
	}
}

func TestTwoSATAssignment(t *testing.T) {
	clauses := []twosat.Clause{{1, 2}, {-1, 2}}
	assignment, ok := twosat.Solve(2, clauses)
	if !ok || !assignment[2] {
		t.Fatalf("assignment = %v, ok = %v", assignment, ok)
	}
}
