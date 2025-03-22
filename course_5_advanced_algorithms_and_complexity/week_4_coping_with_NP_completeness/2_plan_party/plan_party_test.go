package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestPlanPartyFixtures(t *testing.T) {
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
		n, weights, edges, err := parseParty(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		want, err := strconv.Atoi(strings.TrimSpace(string(expected)))
		if err != nil {
			t.Fatal(err)
		}
		if got := maxPartyFun(n, weights, edges); got != int64(want) {
			t.Fatalf("%s: maxPartyFun = %d, want %d", entry.Name(), got, want)
		}
	}
}

func TestPlanPartyBruteForce(t *testing.T) {
	tests := []struct {
		weights []int
		edges   [][2]int
	}{
		{[]int{5}, nil},
		{[]int{1, 2}, [][2]int{{1, 2}}},
		{[]int{1, 5, 3, 7, 5}, [][2]int{{5, 4}, {2, 3}, {4, 2}, {1, 2}}},
		{[]int{3, 4, 5, 6, 7, 8}, [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}}},
	}
	for _, tt := range tests {
		n := len(tt.weights)
		want := 0
		for mask := 0; mask < 1<<n; mask++ {
			valid := true
			for _, edge := range tt.edges {
				if mask&(1<<(edge[0]-1)) != 0 && mask&(1<<(edge[1]-1)) != 0 {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			sum := 0
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					sum += tt.weights[i]
				}
			}
			if sum > want {
				want = sum
			}
		}
		if got := maxPartyFun(n, tt.weights, tt.edges); got != int64(want) {
			t.Fatalf("maxPartyFun(%v, %v) = %d, want %d", tt.weights, tt.edges, got, want)
		}
	}
}
