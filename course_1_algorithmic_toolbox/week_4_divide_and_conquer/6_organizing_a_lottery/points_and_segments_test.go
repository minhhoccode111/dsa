package main

import (
	"reflect"
	"testing"
)

func TestPointsCovered(t *testing.T) {
	tests := []struct {
		name         string
		starts, ends []int
		points       []int
		want         []int
	}{
		{
			"course sample 1",
			[]int{0, 7}, []int{5, 10}, []int{1, 6, 11},
			[]int{1, 0, 0},
		},
		{
			"course sample 3",
			[]int{0, -3, 7}, []int{5, 2, 10}, []int{1, 6},
			[]int{2, 0},
		},
		{
			"endpoints are included",
			[]int{0, 5}, []int{5, 10}, []int{0, 5, 10},
			[]int{1, 2, 1},
		},
		{
			"no points",
			[]int{0}, []int{5}, []int{},
			[]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pointsCovered(tt.starts, tt.ends, tt.points); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("pointsCovered(%v, %v, %v) = %v, want %v", tt.starts, tt.ends, tt.points, got, tt.want)
			}
		})
	}
}
