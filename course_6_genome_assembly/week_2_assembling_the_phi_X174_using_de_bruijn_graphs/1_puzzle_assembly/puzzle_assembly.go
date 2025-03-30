// go run puzzle_assembly.go

// Arrange colored square pieces into a valid grid.

package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func parsePieces(input string) [][4]string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	pieces := make([][4]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "(")
		line = strings.TrimSuffix(line, ")")
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			continue
		}
		var piece [4]string
		for i := range piece {
			piece[i] = strings.TrimSpace(parts[i])
		}
		pieces = append(pieces, piece)
	}
	return pieces
}

func solvePuzzle(pieces [][4]string) ([][]int, bool) {
	count := len(pieces)
	n := int(math.Round(math.Sqrt(float64(count))))
	if n*n != count || n == 0 {
		return nil, false
	}
	grid := make([][]int, n)
	for r := range grid {
		grid[r] = make([]int, n)
		for c := range grid[r] {
			grid[r][c] = -1
		}
	}
	used := make([]bool, count)

	var search func(position int) bool
	search = func(position int) bool {
		if position == count {
			return true
		}
		r, c := position/n, position%n
		for id, piece := range pieces {
			if used[id] {
				continue
			}
			if r == 0 && piece[0] != "black" {
				continue
			}
			if r > 0 && piece[0] != pieces[grid[r-1][c]][2] {
				continue
			}
			if c == 0 && piece[1] != "black" {
				continue
			}
			if c > 0 && piece[1] != pieces[grid[r][c-1]][3] {
				continue
			}
			if r == n-1 && piece[2] != "black" {
				continue
			}
			if c == n-1 && piece[3] != "black" {
				continue
			}
			used[id] = true
			grid[r][c] = id
			if search(position + 1) {
				return true
			}
			grid[r][c] = -1
			used[id] = false
		}
		return false
	}

	if !search(0) {
		return nil, false
	}
	return grid, true
}

func formatPuzzleGrid(pieces [][4]string, grid [][]int) string {
	var builder strings.Builder
	for r, row := range grid {
		if r > 0 {
			builder.WriteByte('\n')
		}
		for c, id := range row {
			if c > 0 {
				builder.WriteByte(';')
			}
			piece := pieces[id]
			fmt.Fprintf(&builder, "(%s,%s,%s,%s)", piece[0], piece[1], piece[2], piece[3])
		}
	}
	return builder.String()
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	pieces := parsePieces(string(data))
	grid, ok := solvePuzzle(pieces)
	if ok {
		fmt.Println(formatPuzzleGrid(pieces, grid))
	}
}
