package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestComputePrefixFunction(t *testing.T) {
	tests := []struct {
		pattern string
		want    []int
	}{
		{"", nil},
		{"A", []int{0}},
		{"AAAA", []int{0, 1, 2, 3}},
		{"ABABACA", []int{0, 0, 1, 2, 3, 0, 1}},
		{"ABCDE", []int{0, 0, 0, 0, 0}},
	}

	for _, tt := range tests {
		if got := computePrefixFunction(tt.pattern); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("computePrefixFunction(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}

func TestFindPattern(t *testing.T) {
	tests := []struct {
		pattern, text string
		want          []int
	}{
		{"GT", "TACG", nil},
		{"ATA", "ATATA", []int{0, 2}},
		{"ATAT", "GATATATGCATATACTT", []int{1, 3, 9}},
		{"A", "AAAA", []int{0, 1, 2, 3}},
		{"ACG", "AC", nil},
	}

	for _, tt := range tests {
		got := findPattern(tt.pattern, tt.text)
		if len(got) == 0 && len(tt.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("findPattern(%q, %q) = %v, want %v", tt.pattern, tt.text, got, tt.want)
		}
	}
}

func TestWritePositions(t *testing.T) {
	tests := []struct {
		positions []int
		want      string
	}{
		{nil, "\n"},
		{[]int{0}, "0\n"},
		{[]int{0, 2, 2, 10}, "0 2 2 10\n"},
		{[]int{999999, 1000000}, "999999 1000000\n"},
	}

	for _, tt := range tests {
		var out bytes.Buffer
		writePositions(&out, tt.positions)
		if got := out.String(); got != tt.want {
			t.Fatalf("writePositions(%v) = %q, want %q", tt.positions, got, tt.want)
		}
	}
}
