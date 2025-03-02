package main

import (
	"sort"
	"testing"
)

func TestInverseBWT(t *testing.T) {
	tests := []struct {
		bwt  string
		want string
	}{
		{"AC$A", "ACA$"},
		{"AGGGAA$", "GAGAGA$"},
		{"$", "$"},
	}

	for _, tt := range tests {
		if got := inverseBWT(tt.bwt); got != tt.want {
			t.Fatalf("inverseBWT(%q) = %q, want %q", tt.bwt, got, tt.want)
		}
	}
}

func TestBWTRoundTrip(t *testing.T) {
	texts := []string{
		"A$",
		"AA$",
		"ACACACAC$",
		"AGACATA$",
		"GAGAGA$",
		"ACGTACGT$",
	}

	for _, text := range texts {
		encoded := bwt(text)
		if got := inverseBWT(encoded); got != text {
			t.Fatalf("inverseBWT(bwt(%q)) = %q, want %q", text, got, text)
		}
	}
}

func naiveBWT(text string) string {
	rotations := make([]string, len(text))
	for i := range rotations {
		rotations[i] = text[i:] + text[:i]
	}
	sort.Strings(rotations)
	result := make([]byte, len(text))
	for i, rotation := range rotations {
		result[i] = rotation[len(rotation)-1]
	}
	return string(result)
}
