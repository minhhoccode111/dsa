package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHeavyHittersFixtures(t *testing.T) {
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
		n, threshold, updates, queries, err := parseHeavy(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got := solveHeavy(n, threshold, updates, queries)
		want := strings.TrimSpace(string(expected))
		if got != want {
			t.Fatalf("%s: got %q, want %q", entry.Name(), got, want)
		}
	}
}

func TestHeavyHittersExact(t *testing.T) {
	updates := []heavyUpdate{
		{id: 1, value: 5},
		{id: 2, value: 4},
		{id: 1, value: -2},
		{id: 2, value: -4},
	}
	got := solveHeavy(2, 3, updates, []int{1, 2})
	if got != "1 0" {
		t.Fatalf("got %q, want %q", got, "1 0")
	}
}
