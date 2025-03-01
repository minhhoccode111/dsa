package strs

import (
	"reflect"
	"sort"
	"testing"
)

func TestSuffixArray(t *testing.T) {
	tests := []struct {
		text string
		want []int
	}{
		{"", nil},
		{"A", []int{0}},
		{"AA$", []int{2, 1, 0}},
		{"GAC$", []int{3, 1, 2, 0}},
		{"GAGAGAGA$", []int{8, 7, 5, 3, 1, 6, 4, 2, 0}},
		{"AACGATAGCGGTAGA$", []int{15, 14, 0, 1, 12, 6, 4, 2, 8, 13, 3, 7, 9, 10, 11, 5}},
	}

	for _, tt := range tests {
		if got := SuffixArray([]byte(tt.text)); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("SuffixArray(%q) = %v, want %v", tt.text, got, tt.want)
		}
	}
}

func TestSuffixArrayMatchesBruteForce(t *testing.T) {
	texts := []string{
		"banana$",
		"mississippi$",
		"abracadabra$",
		"aaaaaa$",
		"ACGTACGT$",
		"zxyxzyx$",
		"$",
	}

	for _, text := range texts {
		want := bruteSuffixArray(text)
		if got := SuffixArray([]byte(text)); !reflect.DeepEqual(got, want) {
			t.Fatalf("SuffixArray(%q) = %v, want %v", text, got, want)
		}
	}
}

func bruteSuffixArray(text string) []int {
	sa := make([]int, len(text))
	for i := range sa {
		sa[i] = i
	}
	sort.Slice(sa, func(i, j int) bool {
		return text[sa[i]:] < text[sa[j]:]
	})
	return sa
}

func TestLCP(t *testing.T) {
	text := []byte("banana$")
	sa := SuffixArray(text)
	want := []int{0, 0, 1, 3, 0, 0, 2}
	if got := LCP(text, sa); !reflect.DeepEqual(got, want) {
		t.Fatalf("LCP(%q) = %v, want %v", text, got, want)
	}
}

func TestRMQ(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	rmq := NewRMQ(values)
	tests := []struct {
		left, right, want int
	}{
		{0, 1, 5},
		{1, 2, 1},
		{0, 5, 1},
		{2, 5, 2},
		{3, 4, 2},
		{2, 3, 4},
	}
	for _, tt := range tests {
		if got := rmq.Min(tt.left, tt.right); got != tt.want {
			t.Fatalf("Min(%d, %d) = %d, want %d", tt.left, tt.right, got, tt.want)
		}
	}
}

func TestBuildSuffixTreeEdgeLabels(t *testing.T) {
	tests := []struct {
		text string
		want []string
	}{
		{"A$", []string{"A$", "$"}},
		{"ACA$", []string{"A", "CA$", "$", "CA$", "$"}},
		{"ATAAATG$", []string{"A", "T", "G$", "$", "AAATG$", "G$", "T", "A", "ATG$", "TG$", "AAATG$", "G$"}},
	}

	for _, tt := range tests {
		text := []byte(tt.text)
		tree := BuildSuffixTree(text, SuffixArray(text), LCP(text, SuffixArray(text)))
		got := tree.EdgeLabels()
		sort.Strings(got)
		want := append([]string(nil), tt.want...)
		sort.Strings(want)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Tree(%q) edge labels = %v, want %v", tt.text, got, want)
		}
	}
}
