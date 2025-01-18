package main

import "testing"

func TestMajorityElement(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		want     int
	}{
		{"course sample majority", []int{2, 3, 9, 2, 2}, 1},
		{"course sample no majority", []int{1, 2, 3, 1}, 0},
		{"single element", []int{7}, 1},
		{"exactly half is not majority", []int{1, 2, 1, 2}, 0},
		{"strict majority", []int{1, 2, 2}, 1},
		{"empty", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := majorityElement(tt.elements); got != tt.want {
				t.Fatalf("majorityElement(%v) = %d, want %d", tt.elements, got, tt.want)
			}
		})
	}
}
