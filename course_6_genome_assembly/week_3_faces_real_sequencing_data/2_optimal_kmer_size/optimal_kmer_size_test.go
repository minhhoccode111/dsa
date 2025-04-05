package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestOptimalKmerSizeFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	reads := parseSequences(string(input))
	want, err := strconv.Atoi(strings.TrimSpace(string(expected)))
	if err != nil {
		t.Fatal(err)
	}
	if got := optimalKmerSize(reads); got != want {
		t.Fatalf("optimalKmerSize = %d, want %d", got, want)
	}
}
