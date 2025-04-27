package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestKUniversalFixtureValidity(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	k, err := strconv.Atoi(strings.TrimSpace(string(input)))
	if err != nil {
		t.Fatal(err)
	}
	got := universalString(k)
	if len(got) != 1<<k {
		t.Fatalf("length = %d, want %d", len(got), 1<<k)
	}
	seen := make(map[string]bool)
	doubled := got + got
	for i := 0; i < len(got); i++ {
		seen[doubled[i:i+k]] = true
	}
	if len(seen) != 1<<k {
		t.Fatalf("only %d distinct %d-mers", len(seen), k)
	}
	want := strings.TrimSpace(string(expected))
	if len(want) != len(got) {
		t.Fatalf("fixture length %d != %d", len(want), len(got))
	}
}

func TestFindEulerianCycleNestedCycles(t *testing.T) {
	// Cycle 0 is 0->1->5->0. Cycle 1 is 2->5->2 and shares only vertex 5.
	// Cycle 1 starts at vertex 2, which never appears on cycle 0.
	edges := [][2]int{{0, 1}, {1, 5}, {2, 5}, {5, 2}, {5, 0}}
	path, ok := FindEulerianCycle(6, edges)
	if !ok {
		t.Fatal("expected a cycle through both nested cycles")
	}
	if len(path) != len(edges)+1 || path[0] != path[len(path)-1] {
		t.Fatalf("path = %v", path)
	}
}

func TestKUniversalSmall(t *testing.T) {
	for k := 1; k <= 5; k++ {
		got := universalString(k)
		if len(got) != 1<<k {
			t.Fatalf("k=%d length=%d", k, len(got))
		}
		seen := make(map[string]bool)
		doubled := got + got
		for i := 0; i < len(got); i++ {
			seen[doubled[i:i+k]] = true
		}
		if len(seen) != 1<<k {
			t.Fatalf("k=%d missing k-mers", k)
		}
	}
}
