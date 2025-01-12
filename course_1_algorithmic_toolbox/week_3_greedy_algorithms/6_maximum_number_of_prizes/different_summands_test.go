package main

import (
	"reflect"
	"testing"
)

func TestOptimalSummands(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"course sample 6", 6, []int{1, 2, 3}},
		{"course sample 8", 8, []int{1, 2, 5}},
		{"course sample 2", 2, []int{2}},
		{"one", 1, []int{1}},
		{"zero", 0, []int{}},
		{"larger", 15, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optimalSummands(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("optimalSummands(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
