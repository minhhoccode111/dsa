// go run merging_tables.go

// Simulate merging tables with a disjoint-set union structure.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Database struct {
	counts      []int64
	ranks       []int
	parents     []int
	maxRowCount int64
}

func newDatabase(counts []int) *Database {
	n := len(counts)
	db := &Database{
		counts:  make([]int64, n),
		ranks:   make([]int, n),
		parents: make([]int, n),
	}
	for i, count := range counts {
		db.counts[i] = int64(count)
		db.parents[i] = i
		db.ranks[i] = 1
		if int64(count) > db.maxRowCount {
			db.maxRowCount = int64(count)
		}
	}
	return db
}

func (db *Database) find(table int) int {
	root := table
	for root != db.parents[root] {
		root = db.parents[root]
	}
	for table != root {
		next := db.parents[table]
		db.parents[table] = root
		table = next
	}
	return root
}

func (db *Database) Merge(dst, src int) bool {
	srcRoot := db.find(src)
	dstRoot := db.find(dst)
	if srcRoot == dstRoot {
		return false
	}

	if db.ranks[srcRoot] > db.ranks[dstRoot] {
		srcRoot, dstRoot = dstRoot, srcRoot
	}

	db.parents[srcRoot] = dstRoot
	db.counts[dstRoot] += db.counts[srcRoot]
	db.counts[srcRoot] = 0

	if db.ranks[srcRoot] == db.ranks[dstRoot] {
		db.ranks[dstRoot]++
	}

	if db.counts[dstRoot] > db.maxRowCount {
		db.maxRowCount = db.counts[dstRoot]
	}

	return true
}

func (db *Database) MaxRowCount() int64 {
	return db.maxRowCount
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var nTables, nQueries int
	if _, err := fmt.Fscan(reader, &nTables, &nQueries); err != nil {
		return
	}

	counts := make([]int, nTables)
	for i := range counts {
		if _, err := fmt.Fscan(reader, &counts[i]); err != nil {
			return
		}
	}

	db := newDatabase(counts)
	for i := 0; i < nQueries; i++ {
		var dst, src int
		if _, err := fmt.Fscan(reader, &dst, &src); err != nil {
			return
		}
		db.Merge(dst-1, src-1)
		fmt.Println(db.MaxRowCount())
	}
}
