package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindMatchingFixtures(t *testing.T) {
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
		n, m, adj, err := parseCrewInput(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got := formatMatching(findMatching(n, m, adj))
		want := strings.TrimSpace(string(expected))
		if got != want {
			t.Fatalf("%s: matching = %q, want %q", entry.Name(), got, want)
		}
	}
}

func TestFindMatchingMaximumCardinality(t *testing.T) {
	adj := [][]bool{
		{true, true, false},
		{true, false, false},
		{false, true, true},
	}
	matching := findMatching(len(adj), len(adj[0]), adj)
	used := make(map[int]bool)
	assigned := 0
	for flight, crew := range matching {
		if crew == -1 {
			continue
		}
		if crew < 0 || crew >= len(adj[0]) || used[crew] || !adj[flight][crew] {
			t.Fatalf("invalid matching %v", matching)
		}
		used[crew] = true
		assigned++
	}
	if assigned != 3 {
		t.Fatalf("assigned %d flights, want 3: %v", assigned, matching)
	}
}
