package main

import "testing"

func TestMaxArithmeticValue(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
	}{
		{"course sample", "5-8+7*4-8+9", 200},
		{"single number", "7", 7},
		{"addition", "1+2", 3},
		{"multiplication before addition", "2*3+4", 14},
		{"subtraction and multiplication", "1-2*3", -3},
		{"addition and multiplication", "1+2*3", 9},
		{"chained subtraction", "2-1-1", 2},
		{"chained multiplication", "1*2*3", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxArithmeticValue(tt.expression); got != tt.want {
				t.Fatalf("maxArithmeticValue(%q) = %d, want %d", tt.expression, got, tt.want)
			}
		})
	}
}
