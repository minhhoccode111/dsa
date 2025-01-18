package main

import "testing"

func TestLCS3(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c []int
		want    int
	}{
		{"course sample 1", []int{1, 2, 3}, []int{2, 1, 3}, []int{1, 3, 5}, 2},
		{"course sample 2", []int{8, 3, 2, 1, 7}, []int{8, 2, 1, 3, 8, 10, 7}, []int{6, 8, 3, 1, 4, 7}, 3},
		{"identical", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, 3},
		{"no common", []int{1}, []int{2}, []int{3}, 0},
		{"empty first", nil, []int{1}, []int{1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcs3(tt.a, tt.b, tt.c); got != tt.want {
				t.Fatalf("lcs3(%v, %v, %v) = %d, want %d", tt.a, tt.b, tt.c, got, tt.want)
			}
		})
	}
}
