package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestKUniversalFixtureValidity(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	k, err := strconv.Atoi(strings.TrimSpace(string(input)))
	if err != nil {
		t.Fatal(err)
	}
	got := universalString(k)
	if len(got) != 1<<k {
		t.Fatalf("length = %d, want %d", len(got), 1<<k)
	}
	seen := make(map[string]bool)
	doubled := got + got
	for i := 0; i < len(got); i++ {
		seen[doubled[i:i+k]] = true
	}
	if len(seen) != 1<<k {
		t.Fatalf("only %d distinct %d-mers", len(seen), k)
	}
	want := strings.TrimSpace(string(expected))
	if len(want) != len(got) {
		t.Fatalf("fixture length %d != %d", len(want), len(got))
	}
}

func TestKUniversalSmall(t *testing.T) {
	for k := 1; k <= 5; k++ {
		got := universalString(k)
		if len(got) != 1<<k {
			t.Fatalf("k=%d length=%d", k, len(got))
		}
		seen := make(map[string]bool)
		doubled := got + got
		for i := 0; i < len(got); i++ {
			seen[doubled[i:i+k]] = true
		}
		if len(seen) != 1<<k {
			t.Fatalf("k=%d missing k-mers", k)
		}
	}
}
