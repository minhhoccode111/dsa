package main

import "testing"

func TestBWT(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{"AA$", "AA$"},
		{"ACACACAC$", "CCCC$AAAA"},
		{"AGACATA$", "ATG$CAAA"},
		{"$", "$"},
	}

	for _, tt := range tests {
		if got := bwt(tt.text); got != tt.want {
			t.Fatalf("bwt(%q) = %q, want %q", tt.text, got, tt.want)
		}
	}
}
