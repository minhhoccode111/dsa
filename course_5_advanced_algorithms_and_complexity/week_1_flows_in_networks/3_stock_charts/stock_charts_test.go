package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestMinChartsFixtures(t *testing.T) {
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
		stocks, err := parseStocks(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		want, err := strconv.Atoi(strings.TrimSpace(string(expected)))
		if err != nil {
			t.Fatal(err)
		}
		if got := minCharts(stocks); got != want {
			t.Fatalf("%s: minCharts = %d, want %d", entry.Name(), got, want)
		}
	}
}

func TestMinChartsExamples(t *testing.T) {
	tests := []struct {
		stocks [][]int
		want   int
	}{
		{[][]int{{1, 2}, {2, 3}, {3, 4}}, 1},
		{[][]int{{1, 1}, {1, 2}, {2, 2}}, 2},
		{[][]int{{1, 2, 3}, {2, 3, 4}, {1, 3, 2}}, 2},
	}
	for _, tt := range tests {
		if got := minCharts(tt.stocks); got != tt.want {
			t.Fatalf("minCharts(%v) = %d, want %d", tt.stocks, got, tt.want)
		}
	}
}
