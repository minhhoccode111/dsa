package main

import "testing"

func TestMergingTables(t *testing.T) {
	db := newDatabase([]int{1, 1, 1, 1, 1})

	merges := []struct {
		dst, src int
		wantMax  int64
	}{
		{2, 4, 2},
		{1, 3, 2},
		{0, 3, 3},
		{4, 3, 5},
		{2, 2, 5},
	}

	for i, merge := range merges {
		db.Merge(merge.dst, merge.src)
		if got := db.MaxRowCount(); got != merge.wantMax {
			t.Fatalf("after merge %d (%d, %d): MaxRowCount() = %d, want %d", i, merge.dst, merge.src, got, merge.wantMax)
		}
	}

	if db.Merge(0, 0) {
		t.Fatal("merge of a table with itself should return false")
	}
}
