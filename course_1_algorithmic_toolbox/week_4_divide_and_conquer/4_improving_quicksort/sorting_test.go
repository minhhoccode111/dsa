package main

import (
	"reflect"
	"testing"
)

func TestSortThreeWay(t *testing.T) {
	tests := []struct {
		name string
		got  []int
		want []int
	}{
		{"course sample", []int{2, 3, 9, 2, 2}, []int{2, 2, 2, 3, 9}},
		{"already sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"reverse", []int{3, 2, 1}, []int{1, 2, 3}},
		{"all equal", []int{5, 5, 5, 5}, []int{5, 5, 5, 5}},
		{"empty", []int{}, []int{}},
		{"single", []int{4}, []int{4}},
		{"negative duplicates", []int{0, -1, 0, -1}, []int{-1, -1, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortThreeWay(tt.got)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("sortThreeWay = %v, want %v", tt.got, tt.want)
			}
		})
	}
}
