package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestMaxFlowFixtures(t *testing.T) {
	dir := filepath.Join("test")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".a") {
			continue
		}
		input, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile(filepath.Join(dir, entry.Name()+".a"))
		if err != nil {
			t.Fatal(err)
		}
		graph, err := parseFlowGraph(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got := graph.maxFlow(0, graph.size()-1)
		want, err := strconv.ParseInt(strings.TrimSpace(string(expected)), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%s: maxFlow = %d, want %d", entry.Name(), got, want)
		}
	}
}

func TestMaxFlowSimple(t *testing.T) {
	graph, err := parseFlowGraph("2 1\n1 2 7\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := graph.maxFlow(0, 1); got != 7 {
		t.Fatalf("maxFlow = %d, want 7", got)
	}
}
