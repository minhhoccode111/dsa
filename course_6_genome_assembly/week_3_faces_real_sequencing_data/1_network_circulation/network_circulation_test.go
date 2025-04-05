package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNetworkCirculationFixtures(t *testing.T) {
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
		n, m, edges, err := parseCirculation(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		flows, ok := findCirculation(n, m, edges)
		wantLines := strings.Split(strings.TrimSpace(string(expected)), "\n")
		wantYes := wantLines[0] == "YES"
		if ok != wantYes {
			t.Fatalf("%s: solvable = %v, want %v", entry.Name(), ok, wantYes)
		}
		if !ok {
			continue
		}
		if len(flows) != m {
			t.Fatalf("%s: got %d flows, want %d", entry.Name(), len(flows), m)
		}
		balance := make([]int64, n)
		for i, edge := range edges {
			if flows[i] < int64(edge.lower) || flows[i] > int64(edge.upper) {
				t.Fatalf("%s: edge %d flow %d outside [%d,%d]", entry.Name(), i, flows[i], edge.lower, edge.upper)
			}
			balance[edge.from] -= flows[i]
			balance[edge.to] += flows[i]
		}
		for vertex, value := range balance {
			if value != 0 {
				t.Fatalf("%s: vertex %d imbalance %d", entry.Name(), vertex+1, value)
			}
		}
	}
}

func TestNetworkCirculationTrivial(t *testing.T) {
	edges := []circulationEdge{
		{from: 0, to: 1, lower: 1, upper: 1},
		{from: 1, to: 0, lower: 1, upper: 1},
	}
	flows, ok := findCirculation(2, 2, edges)
	if !ok || len(flows) != 2 || flows[0] != 1 || flows[1] != 1 {
		t.Fatalf("flows=%v ok=%v", flows, ok)
	}
}
