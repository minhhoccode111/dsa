package main

import "testing"

func TestTopologicalChecker(t *testing.T) {
	input := []byte("4 3\n1 2\n4 1\n3 1\n")
	expected := []byte("4 3 1 2\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "4 3 1 2\n", true},
		{"wrong but valid order", "3 4 1 2\n", true},
		{"reversed edge", "2 3 4 1\n", false},
		{"duplicate", "4 3 1 1\n", false},
		{"missing", "4 3 1\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkTopologicalOrder(input, expected, []byte(tt.output))
			if ok != tt.ok {
				t.Fatalf("checkTopologicalOrder(%q) = %v (%s), want %v", tt.output, ok, detail, tt.ok)
			}
		})
	}
}

func TestKUniversalChecker(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		ok     bool
	}{
		{"valid k=3", "3\n", "00010111\n", true},
		{"valid rotation", "3\n", "00101110\n", true},
		{"wrong length", "3\n", "0001011\n", false},
		{"bad character", "3\n", "0001011x\n", false},
		{"missing kmer", "3\n", "00000000\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkKUniversal([]byte(tt.input), nil, []byte(tt.output))
			if ok != tt.ok {
				t.Fatalf("checkKUniversal(%q) = %v (%s), want %v", tt.output, ok, detail, tt.ok)
			}
		})
	}
}

func TestNonSharedChecker(t *testing.T) {
	input := []byte("CCAAGCTGCTAGAGG CATGCTGGGCTGGCT\n")
	expected := []byte("AA\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "AA\n", true},
		{"valid alternative", "ATG\n", false},
		{"occurs in both", "GG\n", false},
		{"too short", "A\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkNonSharedSubstring(input, expected, []byte(tt.output))
			if ok != tt.ok {
				t.Fatalf("checkNonSharedSubstring(%q) = %v (%s), want %v", tt.output, ok, detail, tt.ok)
			}
		})
	}
}

func TestLongestCommonSubstringChecker(t *testing.T) {
	input := []byte("cool toolbox\naaa bb\n")
	expected := []byte("1 1 3\n0 1 0\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "1 1 3\n0 1 0\n", true},
		{"wrong length", "1 1 2\n0 1 0\n", false},
		{"not common", "0 0 3\n0 1 0\n", false},
		{"out of range", "99 1 3\n0 1 0\n", false},
		{"missing line", "1 1 3\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkLongestCommonSubstring(input, expected, []byte(tt.output))
			if ok != tt.ok {
				t.Fatalf("checkLongestCommonSubstring(%q) = %v (%s), want %v", tt.output, ok, detail, tt.ok)
			}
		})
	}
}
