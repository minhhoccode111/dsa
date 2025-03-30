package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEulerianCycleFixtures(t *testing.T) {
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
		n, edges, err := parseEulerian(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got, ok := formatEulerian(n, edges)
		want := strings.TrimSpace(string(expected))
		if !ok {
			if want != "0" {
				t.Fatalf("%s: got no cycle, want %q", entry.Name(), want)
			}
			continue
		}
		if strings.HasPrefix(want, "0") {
			t.Fatalf("%s: got a cycle, want none", entry.Name())
		}
		path, _ := findEulerian(n, edges)
		checkEulerianPath(t, n, edges, path)
		if got != want {
			t.Logf("%s: valid cycle differs from fixture: got %q, want %q", entry.Name(), got, want)
		}
	}
}

func checkEulerianPath(t *testing.T, n int, edges [][2]int, path []int) {
	t.Helper()
	if len(path) != len(edges)+1 {
		t.Fatalf("path length = %d, want %d", len(path), len(edges)+1)
	}
	used := make([]bool, len(edges))
	for i := 0; i+1 < len(path); i++ {
		found := -1
		for id, edge := range edges {
			if !used[id] && edge[0] == path[i] && edge[1] == path[i+1] {
				found = id
				break
			}
		}
		if found < 0 {
			t.Fatalf("path uses missing edge %d->%d", path[i], path[i+1])
		}
		used[found] = true
	}
	for i, wasUsed := range used {
		if !wasUsed {
			t.Fatalf("edge %d unused", i)
		}
	}
	_ = n
}
