package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorFreeFixtures(t *testing.T) {
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
		reads := readReads(string(input))
		got := assembleErrorFree(reads)
		want := strings.TrimSpace(string(expected))
		if !sameCircularGenome(got, want) {
			t.Fatalf("%s: assembleErrorFree = %q, want a rotation of %q", entry.Name(), got, want)
		}
	}
}

func TestErrorFreeChain(t *testing.T) {
	reads := []string{"AAC", "ACG", "CGT"}
	if got := assembleErrorFree(reads); got != "AACGT" {
		t.Fatalf("assembly = %q, want AACGT", got)
	}
}

func TestErrorFreeIncompleteChain(t *testing.T) {
	reads := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	if got := assembleErrorFree(reads); got != "" {
		t.Fatalf("assembly with no read overlaps = %q, want empty", got)
	}
}

func sameCircularGenome(got, want string) bool {
	if len(got) != len(want) {
		return false
	}
	return strings.Contains(got+got, want)
}
