package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorProneAssemblerFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	reads := parseGenomeReads(string(input))
	got := assembleGenome(reads)
	want := strings.TrimSpace(string(expected))
	if len(got) != len(want) || !strings.Contains(got+got, want) {
		t.Fatalf("assembleGenome = %q, want a rotation of %q", got, want)
	}
}
