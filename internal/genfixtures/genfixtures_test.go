package main

import "testing"

func TestOutputsEquivalent(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same integers", "1 2 3\n", "1 2 3", true},
		{"float precision", "1.414213\n", "1.414213562\n", true},
		{"different integer", "1 2\n", "1 3\n", false},
		{"different token count", "1 2\n", "1\n", false},
		{"inf versus zero", "inf\n", "0.000000000\n", false},
		{"same float", "0.200000\n", "0.2\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outputsEquivalent([]byte(tt.a), []byte(tt.b)); got != tt.want {
				t.Fatalf("outputsEquivalent(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestShapeAllowsInput(t *testing.T) {
	seeds := []Seed{
		{Input: []byte("AAA\r\n1\r\nAA\r\n"), Answer: []byte("0 1")},
		{Input: []byte("AA\r\n1\r\nT\r\n"), Answer: nil},
	}
	shape, ok := shapeFromSeeds(seeds)
	if !ok {
		t.Fatal("shapeFromSeeds returned no meaningful shape")
	}
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"same shape", "AAAA\n1\nAT\n", true},
		{"extra line", "AAA\n1\nAA\n10\n", false},
		{"numeric string line", "AAA\n1\n10\n", false},
		{"unknown character", "AAA\n1\n5AA\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shapeAllowsInput([]byte(tt.input), shape); got != tt.want {
				t.Fatalf("shapeAllowsInput(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestShapeLineRelation(t *testing.T) {
	seeds := []Seed{
		{Input: []byte("1\n2 -1 -1\n"), Answer: []byte("CORRECT")},
		{Input: []byte("3\n2 1 2\n1 -1 -1\n3 -1 -1\n"), Answer: []byte("CORRECT")},
	}
	shape, ok := shapeFromSeeds(seeds)
	if !ok {
		t.Fatal("shapeFromSeeds returned no meaningful shape")
	}
	if shape.LineRelation != lineRelationPlusOne {
		t.Fatalf("LineRelation = %v, want plus-one", shape.LineRelation)
	}
	valid := "2\n2 1 -1\n1 -1 -1\n"
	invalid := "2 2 1 -1 1 -1 -1\n"
	invalidHeader := "2 1 -1\n2 1 -1\n1 -1 -1\n"
	if !shapeAllowsInput([]byte(valid), shape) {
		t.Fatalf("shapeAllowsInput(%q) = false, want true", valid)
	}
	if shapeAllowsInput([]byte(invalid), shape) {
		t.Fatalf("shapeAllowsInput(%q) = true, want false", invalid)
	}
	if shapeAllowsInput([]byte(invalidHeader), shape) {
		t.Fatalf("shapeAllowsInput(%q) = true, want false", invalidHeader)
	}
}

func TestShapeLineStringLength(t *testing.T) {
	seeds := []Seed{
		{Input: []byte("1\nATG\n"), Answer: []byte("ATG")},
		{Input: []byte("1\nTGA\n"), Answer: []byte("TGA")},
	}
	shape, ok := shapeFromSeeds(seeds)
	if !ok {
		t.Fatal("shapeFromSeeds returned no meaningful shape")
	}
	if !shapeAllowsInput([]byte("1\nTGT\n"), shape) {
		t.Fatal("same-length read should be accepted")
	}
	if shapeAllowsInput([]byte("1\nTG\n"), shape) {
		t.Fatal("shorter read should be rejected")
	}
}

func TestInputSignature(t *testing.T) {
	values := []string{"1", "1\n", "1 ", " 1\r\n"}
	for _, value := range values {
		if got := inputSignature([]byte(value)); got != "1" {
			t.Fatalf("inputSignature(%q) = %q, want %q", value, got, "1")
		}
	}
}

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		line string
		want lineKind
	}{
		{"", lineKindNone},
		{"1 -2 3", lineKindNumeric},
		{"ACGT $", lineKindString},
		{"add 1 name", lineKindMixed},
	}
	for _, tt := range tests {
		if got := classifyLine(tt.line); got != tt.want {
			t.Fatalf("classifyLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}
