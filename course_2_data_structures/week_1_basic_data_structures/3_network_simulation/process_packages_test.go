package main

import (
	"reflect"
	"testing"
)

func TestProcessRequests(t *testing.T) {
	tests := []struct {
		name       string
		bufferSize int
		requests   []Request
		want       []int
	}{
		{"empty", 1, nil, []int{}},
		{
			"drop second package",
			1,
			[]Request{{0, 1}, {0, 1}},
			[]int{0, -1},
		},
		{
			"process after previous finishes",
			1,
			[]Request{{0, 1}, {1, 1}},
			[]int{0, 1},
		},
		{
			"one slot",
			1,
			[]Request{{0, 1}, {0, 1}, {0, 1}},
			[]int{0, -1, -1},
		},
		{
			"two slots",
			2,
			[]Request{{0, 1}, {0, 1}, {0, 1}},
			[]int{0, 1, -1},
		},
		{
			"zero buffer",
			0,
			[]Request{{0, 1}},
			[]int{-1},
		},
		{
			"starts after busy",
			1,
			[]Request{{0, 5}, {1, 1}},
			[]int{0, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processRequests(tt.bufferSize, tt.requests); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("processRequests(%d, %v) = %v, want %v", tt.bufferSize, tt.requests, got, tt.want)
			}
		})
	}
}
