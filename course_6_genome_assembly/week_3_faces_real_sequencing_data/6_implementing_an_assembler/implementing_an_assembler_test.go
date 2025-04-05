package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImplementingAnAssemblerFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	reads := parseAssemblerReads(string(input))
	got := assembleContigs(reads)
	want := strings.TrimSpace(string(expected))
	if got != want {
		t.Fatalf("assembleContigs =\n%s\nwant\n%s", got, want)
	}
}
