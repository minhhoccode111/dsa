package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestBubbleDetectionFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	k, threshold, reads, err := parseBubbleInput(string(input))
	if err != nil {
		t.Fatal(err)
	}
	want, err := strconv.Atoi(strings.TrimSpace(string(expected)))
	if err != nil {
		t.Fatal(err)
	}
	if got := countBubbles(k, threshold, reads); got != want {
		t.Fatalf("countBubbles = %d, want %d", got, want)
	}
}
