package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestEnergyValuesFixtures(t *testing.T) {
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
		n, a, b, err := parseEquation(string(input))
		if err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		got, ok := solveEquation(n, a, b)
		wantFields := strings.Fields(string(expected))
		if len(wantFields) == 0 {
			if ok && len(got) != 0 {
				t.Fatalf("%s: got %v, want no output", entry.Name(), got)
			}
			continue
		}
		if !ok {
			t.Fatalf("%s: singular system", entry.Name())
		}
		if len(got) != len(wantFields) {
			t.Fatalf("%s: got %d values, want %d", entry.Name(), len(got), len(wantFields))
		}
		for i, field := range wantFields {
			want, err := strconv.ParseFloat(field, 64)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(got[i]-want) > 1e-6 {
				t.Fatalf("%s: value %d = %v, want %v", entry.Name(), i, got[i], want)
			}
		}
	}
}

func TestEnergyValuesKnown(t *testing.T) {
	got, ok := solveEquation(2, [][]float64{{1, 1}, {2, 3}}, []float64{3, 7})
	if !ok {
		t.Fatal("expected a solution")
	}
	if math.Abs(got[0]-2) > 1e-9 || math.Abs(got[1]-1) > 1e-9 {
		t.Fatalf("got %v, want [2 1]", got)
	}
}
