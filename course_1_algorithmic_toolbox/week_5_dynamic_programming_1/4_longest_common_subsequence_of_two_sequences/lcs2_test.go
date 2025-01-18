package main

import "testing"

func TestLCS2(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want int
	}{
		{"course sample 1", []int{2, 7, 5}, []int{2, 5}, 2},
		{"course sample 2", []int{7}, []int{1, 2, 3, 4}, 0},
		{"course sample 3", []int{2, 7, 8, 3}, []int{5, 2, 8, 7}, 2},
		{"identical", []int{1, 2, 3}, []int{1, 2, 3}, 3},
		{"empty first", nil, []int{1, 2}, 0},
		{"duplicates", []int{1, 1, 1}, []int{1, 1}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcs2(tt.a, tt.b); got != tt.want {
				t.Fatalf("lcs2(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
