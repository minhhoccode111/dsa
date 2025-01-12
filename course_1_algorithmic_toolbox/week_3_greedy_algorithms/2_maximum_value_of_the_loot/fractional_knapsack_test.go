package main

import (
	"math"
	"testing"
)

func TestFractionalKnapsack(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		weights  []int
		values   []int
		want     float64
	}{
		{"course sample 1", 50, []int{20, 50, 30}, []int{60, 100, 120}, 180},
		{"course sample 2", 10, []int{30}, []int{500}, 500.0 / 3.0},
		{"zero capacity", 0, []int{10}, []int{100}, 0},
		{"takes a fraction", 5, []int{10}, []int{100}, 50},
		{"uses best ratio first", 10, []int{5, 10}, []int{10, 15}, 17.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fractionalKnapsack(tt.capacity, tt.weights, tt.values)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("fractionalKnapsack(%d, %v, %v) = %v, want %v", tt.capacity, tt.weights, tt.values, got, tt.want)
			}
		})
	}
}
