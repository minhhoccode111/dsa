package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorProneFixtures(t *testing.T) {
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
		reads := parseReads(string(input))
		got := assembleErrorProne(reads)
		want := strings.TrimSpace(string(expected))
		if !sameCircularGenome(got, want) {
			t.Fatalf("%s: assembleErrorProne = %q, want a rotation of %q", entry.Name(), got, want)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	reads := []string{"AACGT", "AACGG", "TTTTT"}
	got := removeDuplicates(reads, 5, 1)
	if len(got) != 2 || got[0] != "AACGT" || got[1] != "TTTTT" {
		t.Fatalf("removeDuplicates = %v", got)
	}
}

func TestErrorProneIncompleteChain(t *testing.T) {
	reads := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	if got := assembleErrorProne(reads); got != "" {
		t.Fatalf("assembly with no read overlaps = %q, want empty", got)
	}
}

func sameCircularGenome(got, want string) bool {
	if len(got) != len(want) {
		return false
	}
	return strings.Contains(got+got, want)
}
