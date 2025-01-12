package main

import "testing"

func TestMaxDotProduct(t *testing.T) {
	tests := []struct {
		name          string
		first, second []int
		want          int64
	}{
		{"course sample 1", []int{23}, []int{39}, 897},
		{"course sample 2", []int{2, 3, 9}, []int{7, 4, 2}, 79},
		{"negative values", []int{-1, 2}, []int{3, -4}, 10},
		{"single negative", []int{-5}, []int{-7}, 35},
		{"zeros", []int{0, 0}, []int{1, 2}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxDotProduct(tt.first, tt.second); got != tt.want {
				t.Fatalf("maxDotProduct(%v, %v) = %d, want %d", tt.first, tt.second, got, tt.want)
			}
		})
	}
}
