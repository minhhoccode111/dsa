package main

import "testing"

func TestMinRefills(t *testing.T) {
	tests := []struct {
		name     string
		distance int
		tank     int
		stops    []int
		want     int
	}{
		{"course sample 1", 950, 400, []int{200, 375, 550, 750}, 2},
		{"course sample 2", 10, 3, []int{1, 2, 5, 9}, -1},
		{"course sample 3", 200, 250, []int{100, 150}, 0},
		{"no stops and out of range", 10, 3, nil, -1},
		{"no stops but in range", 10, 10, nil, 0},
		{"zero distance", 0, 5, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minRefills(tt.distance, tt.tank, tt.stops); got != tt.want {
				t.Fatalf("minRefills(%d, %d, %v) = %d, want %d", tt.distance, tt.tank, tt.stops, got, tt.want)
			}
		})
	}
}
