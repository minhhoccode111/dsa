package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestSchoolBusFixtures(t *testing.T) {
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
		n, graph, err := parseBus(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		cost, path, ok := shortestBusRoute(n, graph)
		fields := strings.Fields(string(expected))
		wantCost, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			if wantCost != -1 {
				t.Fatalf("%s: got no route, want cost %d", entry.Name(), wantCost)
			}
			continue
		}
		if wantCost == -1 {
			t.Fatalf("%s: got route, want -1", entry.Name())
		}
		if cost != wantCost {
			t.Fatalf("%s: cost = %d, want %d", entry.Name(), cost, wantCost)
		}
		if len(path) != n {
			t.Fatalf("%s: path length = %d, want %d", entry.Name(), len(path), n)
		}
		seen := make([]bool, n)
		var actual int64
		for i, vertex := range path {
			if vertex < 0 || vertex >= n || seen[vertex] {
				t.Fatalf("%s: invalid path %v", entry.Name(), path)
			}
			seen[vertex] = true
			next := path[(i+1)%n]
			if graph[vertex][next] == busInfinity {
				t.Fatalf("%s: path uses missing edge %d-%d", entry.Name(), vertex+1, next+1)
			}
			actual += graph[vertex][next]
		}
		if actual != cost {
			t.Fatalf("%s: path cost = %d, reported %d", entry.Name(), actual, cost)
		}
	}
}

func TestSchoolBusRandomAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for round := 0; round < 100; round++ {
		n := 2 + rng.Intn(6)
		graph := make([][]int64, n)
		for i := range graph {
			graph[i] = make([]int64, n)
			for j := range graph[i] {
				graph[i][j] = busInfinity
			}
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if rng.Intn(2) == 0 {
					weight := int64(1 + rng.Intn(20))
					graph[i][j] = weight
					graph[j][i] = weight
				}
			}
		}
		want := bruteRoute(n, graph)
		got, _, ok := shortestBusRoute(n, graph)
		if want == -1 {
			if ok {
				t.Fatalf("round %d: got route %d, want none", round, got)
			}
		} else if !ok || got != want {
			t.Fatalf("round %d: got %d/%v, want %d", round, got, ok, want)
		}
	}
}

func bruteRoute(n int, graph [][]int64) int64 {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	best := int64(-1)
	var generate func(int)
	generate = func(index int) {
		if index == n {
			var sum int64
			for i := 0; i < n; i++ {
				next := perm[(i+1)%n]
				if graph[perm[i]][next] == busInfinity {
					return
				}
				sum += graph[perm[i]][next]
			}
			if best == -1 || sum < best {
				best = sum
			}
			return
		}
		for i := index; i < n; i++ {
			perm[index], perm[i] = perm[i], perm[index]
			generate(index + 1)
			perm[index], perm[i] = perm[i], perm[index]
		}
	}
	generate(0)
	return best
}
