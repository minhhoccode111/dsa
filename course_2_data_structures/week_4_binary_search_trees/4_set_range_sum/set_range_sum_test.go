package main

import "testing"

func TestRangeSumSet(t *testing.T) {
	set := newRangeSumSet()

	set.insert(1)
	set.insert(1)
	set.insert(2)

	if !set.contains(1) || !set.contains(2) {
		t.Fatal("inserted values are missing")
	}
	if set.contains(3) {
		t.Fatal("unexpected value 3")
	}
	if got := set.rangeSum(1, 2); got != 3 {
		t.Fatalf("rangeSum(1, 2) = %d, want 3", got)
	}

	set.erase(1)
	if set.contains(1) {
		t.Fatal("value 1 was not erased")
	}
	if got := set.rangeSum(0, 5); got != 2 {
		t.Fatalf("rangeSum(0, 5) = %d, want 2", got)
	}
	if got := set.rangeSum(3, 5); got != 0 {
		t.Fatalf("rangeSum(3, 5) = %d, want 0", got)
	}
}
