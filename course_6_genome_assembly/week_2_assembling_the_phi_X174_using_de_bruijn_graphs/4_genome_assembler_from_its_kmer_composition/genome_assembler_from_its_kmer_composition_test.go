package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenomeFromKmersFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
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
