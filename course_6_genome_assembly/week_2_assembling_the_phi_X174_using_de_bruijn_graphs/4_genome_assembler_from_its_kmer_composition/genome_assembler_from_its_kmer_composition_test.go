package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenomeFromKmersFixtures(t *testing.T) {
	entries, err := os.ReadDir("test")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".a") {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("test", entry.Name()))
			if err != nil {
				t.Fatal(err)
			}
			expected, err := os.ReadFile(filepath.Join("test", entry.Name()+".a"))
			if err != nil {
				t.Fatal(err)
			}
			kmers := parseKmers(string(input))
			got := assembleFromKmers(kmers)
			if len(got) != len(kmers) {
				t.Fatalf("genome length = %d, want %d", len(got), len(kmers))
			}
			doubled := got + got
			for _, kmer := range kmers {
				if !strings.Contains(doubled, kmer) {
					t.Fatalf("genome %q does not contain %q", got, kmer)
				}
			}
			want := strings.TrimSpace(string(expected))
			if len(got) == len(want) && !strings.Contains(got+got, want) {
				t.Fatalf("genome %q is not a rotation of %q", got, want)
			}
		})
	}
}

func TestFindEulerianCycleNestedCycles(t *testing.T) {
	// Cycle 0 is 0->1->5->0. Cycle 1 is 2->5->2 and shares only vertex 5.
	// Cycle 1 starts at vertex 2, which never appears on cycle 0.
	edges := [][2]int{{0, 1}, {1, 5}, {2, 5}, {5, 2}, {5, 0}}
	path, ok := FindEulerianCycle(6, edges)
	if !ok {
		t.Fatal("expected a cycle through both nested cycles")
	}
	if len(path) != len(edges)+1 || path[0] != path[len(path)-1] {
		t.Fatalf("path = %v", path)
	}
}

func TestGenomeFromKmersSimple(t *testing.T) {
	kmers := []string{"ATG", "TGC", "GCA", "CAT"}
	got := assembleFromKmers(kmers)
	if len(got) != len(kmers) {
		t.Fatalf("assembly length = %d, want %d", len(got), len(kmers))
	}
	doubled := got + got
	for _, kmer := range kmers {
		if !strings.Contains(doubled, kmer) {
			t.Fatalf("assembly %q does not contain %q", got, kmer)
		}
	}
}
