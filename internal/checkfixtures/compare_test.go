package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareWhitespace(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
		ok   bool
	}{
		{"same", "1 2\n", "1 2\n", true},
		{"extra spaces", "1   2\n", "1 2\n", true},
		{"leading and trailing", "  1 2  \n", "1 2\n", true},
		{"trailing newline", "1 2", "1 2\n", true},
		{"crlf", "1 2\r\n", "1 2\n", true},
		{"different line split", "1\n2\n", "1 2\n", false},
		{"internal blank line", "1\n\n2\n", "1\n2\n", true},
		{"different value", "1 3\n", "1 2\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := compareOutput([]byte(tt.got), []byte(tt.want), Options{}, modeNormal)
			if ok != tt.ok {
				t.Fatalf("compareOutput(%q, %q) = %v (%s), want %v", tt.got, tt.want, ok, detail, tt.ok)
			}
		})
	}
}

func TestCompareFloat(t *testing.T) {
	opts := Options{Float: 1e-6}
	tests := []struct {
		name string
		got  string
		want string
		ok   bool
	}{
		{"precision", "1.414213562\n", "1.414213\n", true},
		{"outside tolerance", "1.500000\n", "1.414213\n", false},
		{"integer exact", "5\n", "4\n", false},
		{"inf versus zero", "inf\n", "0.000000000\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := compareOutput([]byte(tt.got), []byte(tt.want), opts, modeNormal)
			if ok != tt.ok {
				t.Fatalf("compareOutput(%q, %q) = %v (%s), want %v", tt.got, tt.want, ok, detail, tt.ok)
			}
		})
	}
}

func TestCompareRotation(t *testing.T) {
	opts := Options{Float: 1e-6}
	tests := []struct {
		name string
		got  string
		want string
		ok   bool
	}{
		{"string rotation", "AACGTTCG\n", "ACGTTCGA\n", true},
		{"token rotation", "1 4 2 3 4 2 3\n", "2 3 4 2 3 1 4\n", true},
		{"not a rotation", "AACGTTCG\n", "ACGTTCGT\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := compareOutput([]byte(tt.got), []byte(tt.want), opts, modeRotationLines)
			if ok != tt.ok {
				t.Fatalf("compareOutput(%q, %q) = %v (%s), want %v", tt.got, tt.want, ok, detail, tt.ok)
			}
		})
	}
}

func TestCompareModes(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
		mode compareMode
	}{
		{"unordered lines", "b\na\n", "a\nb\n", modeUnorderedLines},
		{"unordered tokens", "4 2 0 1\n", "0 1 2 4\n", modeUnorderedTokens},
		{"drop first then unordered", "text\nb\na\n", "a\nb\n", modeDropFirstUnordered},
		{"smoke", "anything\n", "else\n", modeSmoke},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := compareOutput([]byte(tt.got), []byte(tt.want), Options{}, tt.mode)
			if !ok {
				t.Fatalf("compareOutput(%q, %q) = false (%s), want true", tt.got, tt.want, detail)
			}
		})
	}
}

func TestCompareExact(t *testing.T) {
	if ok, _ := compareOutput([]byte("1 2\n"), []byte("1 2\n"), Options{Exact: true}, modeNormal); !ok {
		t.Fatal("identical exact output rejected")
	}
	if ok, _ := compareOutput([]byte("1  2\n"), []byte("1 2\n"), Options{Exact: true}, modeNormal); ok {
		t.Fatal("whitespace difference accepted in exact mode")
	}
}

func TestFindTests(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "test")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, value string) {
		if err := os.WriteFile(filepath.Join(testDir, name), []byte(value), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("1", "2 3\n")
	write("1.a", "5\n")
	write("2", "not a fixture\n")
	write("dataset.txt", "ignored\n")

	tests, err := findTests(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(tests) != 1 || tests[0].Name != "1" {
		t.Fatalf("findTests = %+v, want only fixture 1", tests)
	}
}
