package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPuzzleAssemblyFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	pieces := parsePieces(string(input))
	grid, ok := solvePuzzle(pieces)
	if !ok {
		t.Fatal("no solution found")
	}
	n := len(grid)
	if n*n != len(pieces) {
		t.Fatalf("grid has %d cells, want %d", n*n, len(pieces))
	}
	used := make([]bool, len(pieces))
	colors := make([][4]string, len(pieces))
	for i, piece := range pieces {
		colors[i] = piece
	}
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			id := grid[r][c]
			if id < 0 || id >= len(pieces) || used[id] {
				t.Fatalf("invalid piece id at %d,%d", r, c)
			}
			used[id] = true
			if r == 0 && colors[id][0] != "black" {
				t.Fatalf("top edge of %d,%d is not black", r, c)
			}
			if c == 0 && colors[id][1] != "black" {
				t.Fatalf("left edge of %d,%d is not black", r, c)
			}
			if r == n-1 && colors[id][2] != "black" {
				t.Fatalf("bottom edge of %d,%d is not black", r, c)
			}
			if c == n-1 && colors[id][3] != "black" {
				t.Fatalf("right edge of %d,%d is not black", r, c)
			}
			if c > 0 && colors[id][1] != colors[grid[r][c-1]][3] {
				t.Fatalf("horizontal mismatch at %d,%d", r, c)
			}
			if r > 0 && colors[id][0] != colors[grid[r-1][c]][2] {
				t.Fatalf("vertical mismatch at %d,%d", r, c)
			}
		}
	}
	for i, wasUsed := range used {
		if !wasUsed {
			t.Fatalf("piece %d not used", i)
		}
	}
	if got := formatPuzzleGrid(pieces, grid); got != strings.TrimSpace(string(expected)) {
		t.Logf("solution differs from fixture but is valid:\n%s", got)
	}
}
