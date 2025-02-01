package main

import "testing"

func TestLongestCommonSubstring(t *testing.T) {
	tests := []struct {
		name       string
		s, t       string
		wantLength int
	}{
		{"course sample 1", "cool", "toolbox", 3},
		{"course sample 2", "aaa", "bb", 0},
		{"course sample 3", "aabaa", "babbaab", 3},
		{"identical", "abc", "abc", 3},
		{"single character", "x", "x", 1},
		{"empty first", "", "abc", 0},
		{"empty second", "abc", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, j, length := longestCommonSubstring(tt.s, tt.t)
			if length != tt.wantLength {
				t.Fatalf("longestCommonSubstring(%q, %q) length = %d, want %d", tt.s, tt.t, length, tt.wantLength)
			}
			if length > 0 && tt.s[i:i+length] != tt.t[j:j+length] {
				t.Fatalf("returned substrings differ: %q vs %q", tt.s[i:i+length], tt.t[j:j+length])
			}
		})
	}
}
