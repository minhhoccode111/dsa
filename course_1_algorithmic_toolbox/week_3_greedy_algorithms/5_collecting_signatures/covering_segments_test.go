package main

import (
	"reflect"
	"testing"
)

func TestOptimalPoints(t *testing.T) {
	tests := []struct {
		name     string
		segments [][2]int
		want     []int
	}{
		{"course sample 1", [][2]int{{1, 3}, {2, 5}, {3, 6}}, []int{3}},
		{"course sample 2", [][2]int{{4, 7}, {1, 3}, {2, 5}, {5, 6}}, []int{3, 6}},
		{"nested segments", [][2]int{{1, 10}, {2, 3}, {4, 5}}, []int{3, 5}},
		{"touching endpoints", [][2]int{{1, 2}, {2, 3}}, []int{2}},
		{"single segment", [][2]int{{5, 9}}, []int{9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optimalPoints(tt.segments); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("optimalPoints(%v) = %v, want %v", tt.segments, got, tt.want)
			}
		})
	}
}
