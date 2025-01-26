package main

import (
	"reflect"
	"testing"
)

func isMinHeap(data []int) bool {
	for i := range data {
		left, right := 2*i+1, 2*i+2
		if left < len(data) && data[left] < data[i] {
			return false
		}
		if right < len(data) && data[right] < data[i] {
			return false
		}
	}
	return true
}

func TestBuildHeap(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  [][2]int
	}{
		{"course sample", []int{5, 4, 3, 2, 1}, [][2]int{{1, 4}, {0, 1}, {1, 3}}},
		{"already heap", []int{1, 2, 3}, [][2]int{}},
		{"needs one swap", []int{3, 2, 1}, [][2]int{{0, 2}}},
		{"single", []int{7}, [][2]int{}},
		{"empty", []int{}, [][2]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append([]int(nil), tt.input...)
			if got := buildHeap(data); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("buildHeap(%v) = %v, want %v", tt.input, got, tt.want)
			}
			if !isMinHeap(data) {
				t.Fatalf("buildHeap(%v) produced non-heap %v", tt.input, data)
			}
		})
	}
}
