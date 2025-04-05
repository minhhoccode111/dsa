package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestTipRemovalFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	reads := parseTipReads(string(input))
	want, err := strconv.Atoi(strings.TrimSpace(string(expected)))
	if err != nil {
		t.Fatal(err)
	}
	if got := removeTips(reads); got != want {
		t.Fatalf("removeTips = %d, want %d", got, want)
	}
}
