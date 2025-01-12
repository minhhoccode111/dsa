package main

import "testing"

func TestChange(t *testing.T) {
	tests := []struct {
		name  string
		money int
		want  int
	}{
		{"course sample 2", 2, 2},
		{"course sample 28", 28, 6},
		{"zero", 0, 0},
		{"nine", 9, 5},
		{"exact ten", 10, 1},
		{"ninety nine", 99, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := change(tt.money); got != tt.want {
				t.Fatalf("change(%d) = %d, want %d", tt.money, got, tt.want)
			}
		})
	}
}
