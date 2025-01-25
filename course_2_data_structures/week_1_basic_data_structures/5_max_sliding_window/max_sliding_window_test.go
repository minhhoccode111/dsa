package main

import (
	"reflect"
	"testing"
)

func TestMaxSlidingWindow(t *testing.T) {
	tests := []struct {
		name     string
		sequence []int
		m        int
		want     []int
	}{
		{"course sample", []int{2, 7, 3, 1, 5, 2, 6, 2}, 4, []int{7, 7, 5, 6, 6}},
		{"window one", []int{1, 3, 2}, 1, []int{1, 3, 2}},
		{"window all", []int{1, 3, 2}, 3, []int{3}},
		{"duplicates", []int{1, 3, 1, 2, 0, 5}, 3, []int{3, 3, 2, 5}},
		{"single", []int{5}, 1, []int{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxSlidingWindow(tt.sequence, tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("maxSlidingWindow(%v, %d) = %v, want %v", tt.sequence, tt.m, got, tt.want)
			}
		})
	}
}
