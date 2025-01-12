package main

import "testing"

func TestLargestNumber(t *testing.T) {
	tests := []struct {
		name    string
		numbers []string
		want    string
	}{
		{"course sample 1", []string{"21", "2"}, "221"},
		{"course sample 2", []string{"9", "4", "6", "1", "9"}, "99641"},
		{"course sample 3", []string{"23", "39", "92"}, "923923"},
		{"all zeros", []string{"0", "0"}, "0"},
		{"one before ten", []string{"10", "1"}, "110"},
		{"classic", []string{"3", "30", "34", "5", "9"}, "9534330"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := largestNumber(tt.numbers); got != tt.want {
				t.Fatalf("largestNumber(%v) = %q, want %q", tt.numbers, got, tt.want)
			}
		})
	}
}
