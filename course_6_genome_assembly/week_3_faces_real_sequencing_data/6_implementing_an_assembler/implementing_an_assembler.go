// go run implementing_an_assembler.go

// Assemble contigs from single reads or read-pairs using de Bruijn graphs.

package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

type readPair struct {
	first    []byte
	second   []byte
	distance int
}

type assemblerDataset struct {
	reads   [][]byte
	pairs   []readPair
	paired  bool
	kMaxCap int
	// resolveForksAlways keeps indel-fork resolution on for the polish
	// rounds even after error correction has lowered the detected error
	// fraction below the normal threshold.
	resolveForksAlways bool
}

func parseAssemblerDataset(input []byte) *assemblerDataset {
	dataset := &assemblerDataset{}
	start := 0
	firstData := -1
	for start < len(input) {
		line, next, ok := nextInputLine(input, start)
		start = next
		if !ok {
			break
		}
		if len(line) == 0 {
			continue
		}
		firstData = start
		break
	}
	if firstData < 0 {
		return dataset
	}
	start = firstData
	for start < len(input) {
		line, next, ok := nextInputLine(input, start)
		start = next
		if !ok {
			break
		}
		if len(line) == 0 {
			continue
		}
		if bytes.IndexByte(line, '|') >= 0 {
			dataset.paired = true
			break
		}
	}
	start = firstData
	for start < len(input) {
		line, next, ok := nextInputLine(input, start)
		start = next
		if !ok {
			break
		}
		if len(line) == 0 {
			continue
		}
		if !dataset.paired {
			upperInPlace(line)
			dataset.reads = append(dataset.reads, line)
			continue
		}
		firstSeparator := bytes.IndexByte(line, '|')
		if firstSeparator < 0 {
			continue
		}
		secondSeparator := bytes.IndexByte(line[firstSeparator+1:], '|')
		first := bytes.TrimSpace(line[:firstSeparator])
		var second []byte
		var distanceField []byte
		if secondSeparator < 0 {
			second = bytes.TrimSpace(line[firstSeparator+1:])
		} else {
			secondSeparator += firstSeparator + 1
			second = bytes.TrimSpace(line[firstSeparator+1 : secondSeparator])
			distanceField = bytes.TrimSpace(line[secondSeparator+1:])
		}
		upperInPlace(first)
		upperInPlace(second)
		if len(first) == 0 || len(second) == 0 {
			continue
		}
		distance := 0
		if len(distanceField) > 0 {
			if parsed, err := strconv.Atoi(string(distanceField)); err == nil {
				distance = parsed
			}
		}
		dataset.pairs = append(dataset.pairs, readPair{first: first, second: second, distance: distance})
		dataset.reads = append(dataset.reads, first, second)
	}
	return dataset
}

func nextInputLine(input []byte, start int) ([]byte, int, bool) {
	if start >= len(input) {
		return nil, start, false
	}
	end := bytes.IndexByte(input[start:], '\n')
	if end < 0 {
		end = len(input)
	} else {
		end += start
	}
	return bytes.TrimSpace(input[start:end]), end + 1, true
}

func upperInPlace(data []byte) {
	for i, base := range data {
		if base >= 'a' && base <= 'z' {
			data[i] = base - ('a' - 'A')
		}
	}
}

func parseAssemblerReads(input string) []string {
	dataset := parseAssemblerDataset([]byte(input))
	reads := make([]string, len(dataset.reads))
	for i, read := range dataset.reads {
		reads[i] = string(read)
	}
	return reads
}

type kmer struct {
	lo uint64
	hi uint64
}

func encodeKmer(seq []byte, k int) (kmer, bool) {
	if k <= 0 || len(seq) < k {
		return kmer{}, false
	}
	var code kmer
	for _, base := range seq[:k] {
		value, ok := baseValue(base)
		if !ok {
			return kmer{}, false
		}
		code = code.appendCode(value, k)
	}
	return code, true
}

func baseValue(base byte) (byte, bool) {
	switch base {
	case 'A':
		return 0, true
	case 'C':
		return 1, true
	case 'G':
		return 2, true
	case 'T':
		return 3, true
	default:
		return 0, false
	}
}

func decodeKmer(code kmer, k int) []byte {
	seq := make([]byte, k)
	for i := k - 1; i >= 0; i-- {
		seq[i] = "ACGT"[code.lo&3]
		code = code.prefix()
	}
	return seq
}

func (code kmer) prefix() kmer {
	return kmer{lo: (code.lo >> 2) | (code.hi << 62), hi: code.hi >> 2}
}

func (code kmer) suffix(k int) kmer {
	bits := 2 * (k - 1)
	if bits <= 0 {
		return kmer{}
	}
	if bits < 64 {
		return kmer{lo: code.lo & ((uint64(1) << uint(bits)) - 1)}
	}
	if bits == 64 {
		return kmer{lo: code.lo}
	}
	return kmer{lo: code.lo, hi: code.hi & ((uint64(1) << uint(bits-64)) - 1)}
}

func (code kmer) appendBase(base byte, k int) kmer {
	value, ok := baseValue(base)
	if !ok {
		return code
	}
	return code.appendCode(value, k)
}

func (code kmer) appendCode(value byte, k int) kmer {
	hi := (code.hi << 2) | (code.lo >> 62)
	lo := (code.lo << 2) | uint64(value&3)
	bits := 2 * k
	if bits <= 64 {
		if bits < 64 {
			lo &= (uint64(1) << uint(bits)) - 1
		}
		hi = 0
	} else {
		hi &= (uint64(1) << uint(bits-64)) - 1
	}
	return kmer{lo: lo, hi: hi}
}

func (code kmer) reverseComplement(k int) kmer {
	var result kmer
	for i := 0; i < k; i++ {
		base := byte(code.lo & 3)
		result = result.appendCode(3-base, k)
		code = code.prefix()
	}
	return result
}

func canonicalKmer(code kmer, k int) kmer {
	rc := code.reverseComplement(k)
	if kmerLess(rc, code) {
		return rc
	}
	return code
}

func kmerLess(a, b kmer) bool {
	if a.hi != b.hi {
		return a.hi < b.hi
	}
	return a.lo < b.lo
}

type countedKMers struct {
	keys              []kmer
	counts            []uint32
	totalInstances    uint64
	estimatedCoverage float64
	fillMinSolid      uint32
	indexSlots        []int32
	indexMask         uint64
}

func (c *countedKMers) buildIndex() {
	if len(c.keys) == 0 || len(c.indexSlots) > 0 {
		return
	}
	size := 16
	for size < len(c.keys)*2 {
		size <<= 1
	}
	slots := make([]int32, size)
	for i := range slots {
		slots[i] = -1
	}
	mask := uint64(size - 1)
	for i, key := range c.keys {
		slot := int(hashKmer(key) & mask)
		for slots[slot] >= 0 {
			slot = (slot + 1) & int(mask)
		}
		slots[slot] = int32(i)
	}
	c.indexSlots = slots
	c.indexMask = mask
}

type kmerTable struct {
	keys   []kmer
	counts []uint32
	used   []uint64
	mask   uint64
	size   int
}

func newKmerTable(capacity int) *kmerTable {
	size := 16
	for size < capacity {
		size <<= 1
	}
	return &kmerTable{
		keys:   make([]kmer, size),
		counts: make([]uint32, size),
		used:   make([]uint64, (size+63)/64),
		mask:   uint64(size - 1),
	}
}

func (t *kmerTable) usedAt(index int) bool {
	return t.used[index>>6]&(uint64(1)<<uint(index&63)) != 0
}

func (t *kmerTable) setUsed(index int) {
	t.used[index>>6] |= uint64(1) << uint(index&63)
}

func (t *kmerTable) slot(key kmer) (int, bool) {
	index := int(hashKmer(key) & t.mask)
	for {
		if !t.usedAt(index) {
			return index, false
		}
		if t.keys[index] == key {
			return index, true
		}
		index = (index + 1) & int(t.mask)
	}
}

func (t *kmerTable) add(key kmer) int {
	if (t.size+1)*4 >= len(t.keys)*3 {
		t.grow()
	}
	index, found := t.slot(key)
	if !found {
		t.keys[index] = key
		t.setUsed(index)
		t.size++
	}
	return index
}

func (t *kmerTable) grow() {
	next := newKmerTable(len(t.keys) * 2)
	for index := range t.keys {
		if !t.usedAt(index) {
			continue
		}
		key := t.keys[index]
		slot, _ := next.slot(key)
		next.keys[slot] = key
		next.setUsed(slot)
		next.counts[slot] = t.counts[index]
		next.size++
	}
	*t = *next
}

func hashKmer(key kmer) uint64 {
	h := key.lo ^ (key.hi + 0x9e3779b97f4a7c15 + (key.lo << 6) + (key.lo >> 2))
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	return h
}

func kmerBloomHashes(code kmer) (uint64, uint64) {
	first := hashKmer(code)
	second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
	return first, second
}

type kmerBloom struct {
	bits []uint64
	mask uint64
}

func newKmerBloom(pow2Bits uint) *kmerBloom {
	bloomBits := uint64(1) << pow2Bits
	return &kmerBloom{
		bits: make([]uint64, bloomBits/64),
		mask: bloomBits - 1,
	}
}

func (b *kmerBloom) add(code kmer) {
	first, second := kmerBloomHashes(code)
	i1 := first & b.mask
	b.bits[i1>>6] |= uint64(1) << uint(i1&63)
	i2 := second & b.mask
	b.bits[i2>>6] |= uint64(1) << uint(i2&63)
}

func (b *kmerBloom) contains(code kmer) bool {
	first, second := kmerBloomHashes(code)
	i1 := first & b.mask
	if b.bits[i1>>6]&(uint64(1)<<uint(i1&63)) == 0 {
		return false
	}
	i2 := second & b.mask
	return b.bits[i2>>6]&(uint64(1)<<uint(i2&63)) != 0
}

func (t *kmerTable) find(key kmer) (int, bool) {
	index, found := t.slot(key)
	if !found {
		return 0, false
	}
	return index, true
}

func countSolidKMers(dataset *assemblerDataset, k int, canonical bool) *countedKMers {
	return countKMers(dataset, k, canonical, false, false)
}

func countKMers(dataset *assemblerDataset, k int, canonical, bothStrands, releaseReads bool) *countedKMers {
	instances := kmerInstances(dataset, k)
	if bothStrands {
		instances *= 2
	}
	if len(dataset.reads) <= 100 || instances < 500000 {
		return countDirectMode(dataset, k, canonical, bothStrands)
	}
	return countBloomMode(dataset, k, canonical, bothStrands, instances, releaseReads)
}

func countDirectMode(dataset *assemblerDataset, k int, canonical, bothStrands bool) *countedKMers {
	table := newKmerTable(1024)
	total := uint64(0)
	eachKmerMode(dataset, k, canonical, bothStrands, func(code kmer) {
		index := table.add(code)
		table.counts[index]++
		total++
	})
	return extractCounts(table, 1, total, k)
}

func countDirectLarge(dataset *assemblerDataset, k int, canonical bool, total uint64, releaseReads bool) *countedKMers {
	return countDirectLargeMode(dataset, k, canonical, false, total, releaseReads)
}

func countDirectLargeMode(dataset *assemblerDataset, k int, canonical, bothStrands bool, total uint64, releaseReads bool) *countedKMers {
	if bothStrands {
		total *= 2
	}
	initial := 1 << 20
	if total > 20_000_000 {
		initial = 1 << 23
	} else if total > 8_000_000 {
		initial = 1 << 22
	}
	table := newKmerTable(initial)
	eachKmerMode(dataset, k, canonical, bothStrands, func(code kmer) {
		index := table.add(code)
		table.counts[index]++
	})
	if releaseReads {
		dataset.reads = nil
	}
	runtime.GC()
	debug.FreeOSMemory()
	return extractCounts(table, 1, total, k)
}

func kmerInstances(dataset *assemblerDataset, k int) uint64 {
	total := uint64(0)
	for _, read := range dataset.reads {
		if len(read) >= k {
			total += uint64(len(read) - k + 1)
		}
	}
	return total
}

func countBloomMode(dataset *assemblerDataset, k int, canonical, bothStrands bool, total uint64, releaseReads bool) *countedKMers {
	bloom := newKmerBloom(29)
	initial := 1 << 20
	if total > 20_000_000 {
		initial = 1 << 23
	} else if total > 8_000_000 {
		initial = 1 << 22
	}
	candidate := newKmerTable(initial)
	eachKmerMode(dataset, k, canonical, bothStrands, func(code kmer) {
		if bloom.contains(code) {
			candidate.add(code)
		}
		bloom.add(code)
	})
	eachKmerMode(dataset, k, canonical, bothStrands, func(code kmer) {
		if index, ok := candidate.find(code); ok {
			candidate.counts[index]++
		}
	})
	if releaseReads {
		dataset.reads = nil
	}
	bloom = nil
	runtime.GC()
	debug.FreeOSMemory()
	return extractCounts(candidate, 2, total, k)
}

func canonPair(fwd, rc kmer) kmer {
	if kmerLess(rc, fwd) {
		return rc
	}
	return fwd
}

// scanReadKMers walks read in a rolling 2-bit pass, invoking visit with each
// valid k-mer and its reverse complement; returns false if stopped on invalid bases.
func scanReadKMers(read []byte, k int, resetOnInvalid bool, visit func(fwd, rc kmer) bool) bool {
	if k <= 0 || len(read) < k {
		return true
	}
	var code, rc kmer
	run := 0
	for i := 0; i < len(read); i++ {
		v, ok := baseValue(read[i])
		if !ok {
			if !resetOnInvalid {
				return false
			}
			run = 0
			code = kmer{}
			rc = kmer{}
			continue
		}
		code = code.appendCode(v, k)
		rc = prependCode(rc.prefix(), 3-v, k)
		run++
		if run < k {
			continue
		}
		if !visit(code, rc) {
			return true
		}
	}
	return true
}

func eachKmer(dataset *assemblerDataset, k int, canonical bool, visit func(kmer)) {
	eachKmerMode(dataset, k, canonical, false, visit)
}

func eachKmerMode(dataset *assemblerDataset, k int, canonical, bothStrands bool, visit func(kmer)) {
	if k <= 0 {
		return
	}
	for _, read := range dataset.reads {
		scanReadKMers(read, k, true, func(code, rc kmer) bool {
			if canonical {
				visit(canonPair(code, rc))
				return true
			}
			visit(code)
			if bothStrands && rc != code {
				visit(rc)
			}
			return true
		})
	}
}

type kmerCountSorter struct {
	keys   []kmer
	counts []uint32
}

func (s kmerCountSorter) Len() int           { return len(s.keys) }
func (s kmerCountSorter) Less(i, j int) bool { return kmerLess(s.keys[i], s.keys[j]) }
func (s kmerCountSorter) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
	s.counts[i], s.counts[j] = s.counts[j], s.counts[i]
}

func extractCounts(table *kmerTable, minCount uint32, total uint64, k int) *countedKMers {
	result := &countedKMers{totalInstances: total}
	for index := range table.keys {
		if table.usedAt(index) && table.counts[index] >= minCount {
			result.keys = append(result.keys, table.keys[index])
			result.counts = append(result.counts, table.counts[index])
		}
	}
	sort.Sort(kmerCountSorter{keys: result.keys, counts: result.counts})
	if len(result.keys) > 0 {
		result.estimatedCoverage = float64(total) / float64(len(result.keys))
		const maxHistogramCount = 4096
		histogram := make([]int, maxHistogramCount+1)
		for _, count := range result.counts {
			if count <= maxHistogramCount {
				histogram[count]++
			}
		}
		mode := 0
		modeCount := 0
		for count := 3; count < len(histogram); count++ {
			if histogram[count] > modeCount {
				modeCount = histogram[count]
				mode = count
			}
		}
		if mode > 0 {
			result.estimatedCoverage = float64(mode)
		}
		// Robust dominant-coverage estimate: maximize c * (#keys with count
		// >= c). Junk k-mers inflate key counts at low c but carry little
		// depth, so the true coverage wins this product even when error
		// keys outnumber real ones many times over.
		const minEstimate = 3
		if len(result.counts) > 1000 {
			suffix := make([]int64, maxHistogramCount+2)
			for _, count := range result.counts {
				if count > maxHistogramCount {
					count = maxHistogramCount
				}
				suffix[count]++
			}
			startC := minEstimate
			// Skip monotonically decaying error tail when low-count keys dominate.
			if histogram[2] > 0 && histogram[3] > 2*histogram[5] {
				for startC < 20 && histogram[startC] > histogram[startC+1] {
					startC++
				}
			}
			for c := maxHistogramCount; c >= startC; c-- {
				suffix[c] += suffix[c+1]
			}
			best := 0
			bestScore := int64(0)
			for c := startC; c <= maxHistogramCount; c++ {
				score := int64(c) * suffix[c]
				if score > bestScore {
					bestScore = score
					best = c
				}
			}
			if best > 0 {
				result.estimatedCoverage = float64(best)
			}
		}
	}
	return result
}

func (c *countedKMers) index(code kmer) int {
	if len(c.indexSlots) > 0 {
		slot := int(hashKmer(code) & c.indexMask)
		for {
			idx := c.indexSlots[slot]
			if idx < 0 {
				return -1
			}
			if c.keys[idx] == code {
				return int(idx)
			}
			slot = (slot + 1) & int(c.indexMask)
		}
	}
	low, high := 0, len(c.keys)
	for low < high {
		middle := (low + high) / 2
		if kmerLess(c.keys[middle], code) {
			low = middle + 1
		} else {
			high = middle
		}
	}
	if low < len(c.keys) && c.keys[low] == code {
		return low
	}
	return -1
}

func (c *countedKMers) count(code kmer) uint32 {
	index := c.index(code)
	if index < 0 {
		return 0
	}
	return c.counts[index]
}

func estimateErrorFraction(dataset *assemblerDataset, counts *countedKMers, k int) float64 {
	if len(dataset.reads) == 0 {
		return 0
	}
	step := len(dataset.reads) / 1000
	if step < 1 {
		step = 1
	}
	total := 0
	absent := 0
	for index := 0; index < len(dataset.reads); index += step {
		read := dataset.reads[index]
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			total++
			if counts.count(code) == 0 {
				absent++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(absent) / float64(total)
}

type assemblyConfig struct {
	k            int
	minCount     uint32
	filter       bool
	filterRatio  float64
	resolveForks bool
	bothStrands  bool
}

type unitig struct {
	edges    []int32
	coverage float64
	closed   bool
}

func prependCode(vertex kmer, value byte, k int) kmer {
	bits := 2 * (k - 1)
	if bits < 64 {
		return kmer{lo: vertex.lo | (uint64(value&3) << uint(bits))}
	}
	if bits == 64 {
		return kmer{lo: vertex.lo, hi: uint64(value & 3)}
	}
	return kmer{lo: vertex.lo, hi: vertex.hi | (uint64(value&3) << uint(bits-64))}
}

func (c *countedKMers) vertexInLive(vertex kmer, k int, removed []bool) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		index := c.index(prependCode(vertex, base, k))
		if index >= 0 && (removed == nil || !removed[index]) {
			degree++
		}
	}
	return degree
}

func (c *countedKMers) vertexOutLive(vertex kmer, k int, removed []bool) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		index := c.index(vertex.appendCode(base, k))
		if index >= 0 && (removed == nil || !removed[index]) {
			degree++
		}
	}
	return degree
}

func (c *countedKMers) uniqueOutLive(vertex kmer, k int, removed []bool) (kmer, bool) {
	var next kmer
	found := false
	for base := byte(0); base < 4; base++ {
		candidate := vertex.appendCode(base, k)
		index := c.index(candidate)
		if index >= 0 && (removed == nil || !removed[index]) {
			if found {
				return kmer{}, false
			}
			next = candidate
			found = true
		}
	}
	return next, found
}

func buildUnitigs(counts *countedKMers, k int, removed []bool, avoidHairpin bool) []unitig {
	used := make([]bool, len(counts.keys))
	units := make([]unitig, 0)
	walk := func(start int) unitig {
		unit := unitig{}
		current := start
		for {
			used[current] = true
			unit.edges = append(unit.edges, int32(current))
			destination := counts.keys[current].suffix(k)
			if counts.vertexInLive(destination, k, removed) != 1 || counts.vertexOutLive(destination, k, removed) != 1 {
				break
			}
			next, ok := counts.uniqueOutLive(destination, k, removed)
			if !ok {
				break
			}
			// In double-stranded pooled graphs, do not fold across a palindromic vertex
			// back onto the reverse complement of the current edge (hairpin chimera).
			if avoidHairpin && next == counts.keys[current].reverseComplement(k) {
				break
			}
			nextIndex := counts.index(next)
			if nextIndex < 0 {
				break
			}
			if used[nextIndex] {
				unit.closed = nextIndex == start
				break
			}
			current = nextIndex
		}
		total := 0.0
		for _, edge := range unit.edges {
			total += float64(counts.counts[edge])
		}
		if len(unit.edges) > 0 {
			unit.coverage = total / float64(len(unit.edges))
		}
		return unit
	}
	for index, key := range counts.keys {
		if removed[index] || used[index] {
			continue
		}
		source := key.prefix()
		if counts.vertexInLive(source, k, removed) == 1 && counts.vertexOutLive(source, k, removed) == 1 {
			continue
		}
		units = append(units, walk(index))
	}
	for index := range counts.keys {
		if removed[index] || used[index] {
			continue
		}
		units = append(units, walk(index))
	}
	return units
}

func unitigString(counts *countedKMers, unit unitig, k int) string {
	if len(unit.edges) == 0 {
		return ""
	}
	first := counts.keys[unit.edges[0]]
	result := decodeKmer(first, k)
	for _, edge := range unit.edges[1:] {
		result = append(result, byte("ACGT"[counts.keys[edge].lo&3]))
	}
	if unit.closed && len(result) >= k-1 {
		result = result[:len(result)-(k-1)]
	}
	return string(result)
}

func minRotation(seq string) string {
	n := len(seq)
	if n == 0 {
		return seq
	}
	doubled := seq + seq
	i, j, k := 0, 1, 0
	for i < n && j < n && k < n {
		a := doubled[i+k]
		b := doubled[j+k]
		if a == b {
			k++
			continue
		}
		if a > b {
			i += k + 1
			if i <= j {
				i = j + 1
			}
		} else {
			j += k + 1
			if j <= i {
				j = i + 1
			}
		}
		k = 0
	}
	start := i
	if j < start {
		start = j
	}
	return doubled[start : start+n]
}

func canonicalCircular(seq string) string {
	forward := minRotation(seq)
	reverse := minRotation(string(reverseComplementSequence([]byte(seq))))
	if reverse < forward {
		return reverse
	}
	return forward
}

func assembleAtK(counts *countedKMers, cfg assemblyConfig) []string {
	if cfg.k < 2 || len(counts.keys) == 0 {
		return nil
	}
	counts.buildIndex()
	removed := make([]bool, len(counts.keys))
	protected := make([]bool, len(counts.keys))
	var rcIdx []int32
	if cfg.bothStrands {
		rcIdx = make([]int32, len(counts.keys))
		for i, key := range counts.keys {
			rcIdx[i] = int32(counts.index(key.reverseComplement(cfg.k)))
		}
	}
	symmetrize := func() {
		if !cfg.bothStrands {
			return
		}
		for i := range counts.keys {
			if rc := rcIdx[i]; rc >= 0 {
				if removed[i] {
					removed[rc] = true
				}
				if protected[i] {
					protected[rc] = true
				}
			}
		}
	}
	filter := cfg.filter && counts.estimatedCoverage >= 4
	maxLength := 2 * cfg.k
	if maxLength < 500 {
		maxLength = 500
	}
	threshold := 0.30 * counts.estimatedCoverage
	if cfg.filterRatio > 0 {
		threshold = cfg.filterRatio * counts.estimatedCoverage
	}
	if counts.estimatedCoverage < 10 {
		if threshold < 2.1 {
			threshold = 2.1
		}
	} else if threshold < 3 {
		threshold = 3
	}
	if filter {
		if cfg.resolveForks {
			resolveForks(counts, cfg.k, removed, protected, maxLength)
			symmetrize()
			// Indel-heavy data splits every homopolymer into parallel shifted
			// graph paths; collapse each bubble to its best-supported path and
			// shield the winner from the coverage filter so reads stay connected.
			collapseBubbles(counts, cfg.k, removed, protected, maxLength)
			symmetrize()
		}
	}
	for pass := 0; pass < 20 && filter; pass++ {
		units := buildUnitigs(counts, cfg.k, removed, cfg.bothStrands)
		changed := false
		for _, unit := range units {
			if len(unit.edges) == 0 || len(unit.edges) > maxLength {
				continue
			}
			if unit.coverage < threshold && !unitProtected(unit, protected) {
				for _, edge := range unit.edges {
					removed[edge] = true
				}
				changed = true
			}
		}
		symmetrize()
		if !changed {
			break
		}
	}
	units := buildUnitigs(counts, cfg.k, removed, cfg.bothStrands)
	contigs := make([]string, 0, len(units))
	seen := make(map[string]bool, len(units))
	for _, unit := range units {
		contig := unitigString(counts, unit, cfg.k)
		if contig == "" {
			continue
		}
		if unit.closed {
			contig = canonicalCircular(contig)
		}
		if seen[contig] {
			continue
		}
		seen[contig] = true
		contigs = append(contigs, contig)
	}
	sort.Strings(contigs)
	return contigs
}

func unitProtected(unit unitig, protected []bool) bool {
	for _, edge := range unit.edges {
		if protected[edge] {
			return true
		}
	}
	return false
}

// resolveForks kills indel-induced bubbles: from every branch vertex it walks
// each outgoing path until the next complex vertex (bounded by maxLength) and
// groups the walks by their end vertex. Walks that reconverge form a bubble;
// only its best-covered path survives, shielded from the coverage filter.
// True repeat branches do not reconverge inside the cap and stay intact, and
// endpoint identity is never compared, so fragmented divergence nodes do not
// hide the bubbles.
func resolveForks(counts *countedKMers, k int, removed, protected []bool, maxLength int) {
	for pass := 0; pass < 6; pass++ {
		outAdj := make(map[kmer][]int32)
		inDeg := make(map[kmer]int)
		for i := range counts.keys {
			if removed[i] {
				continue
			}
			key := counts.keys[i]
			outAdj[key.prefix()] = append(outAdj[key.prefix()], int32(i))
			inDeg[key.suffix(k)]++
		}
		changed := false
		type walkResult struct {
			end   kmer
			cov   float64
			edges []int32
		}
		for _, edges := range outAdj {
			if len(edges) < 2 {
				continue
			}
			results := make([]walkResult, 0, len(edges))
			for _, start := range edges {
				path := []int32{start}
				total := float64(counts.counts[start])
				cur := counts.keys[start].suffix(k)
				for step := 1; step < maxLength; step++ {
					if inDeg[cur] > 1 {
						break // reconvergence: bubble end
					}
					outs := outAdj[cur]
					if len(outs) == 0 {
						break
					}
					// Follow the dominant continuation so stray low-count
					// junk edges cannot cut the class walk short.
					next := outs[0]
					for _, candidate := range outs[1:] {
						if counts.counts[candidate] > counts.counts[next] {
							next = candidate
						}
					}
					path = append(path, next)
					total += float64(counts.counts[next])
					cur = counts.keys[next].suffix(k)
				}
				results = append(results, walkResult{
					end:   cur,
					cov:   total / float64(len(path)),
					edges: path,
				})
			}
			byEnd := make(map[kmer][]int)
			for idx, result := range results {
				if len(result.edges) < maxLength {
					byEnd[result.end] = append(byEnd[result.end], idx)
				}
			}
			for _, group := range byEnd {
				if len(group) < 2 {
					continue
				}
				best := group[0]
				for _, gi := range group[1:] {
					if results[gi].cov > results[best].cov {
						best = gi
					}
				}
				for _, gi := range group {
					if gi == best || (results[gi].cov >= 0.85*results[best].cov && results[gi].cov >= 0.75*counts.estimatedCoverage) {
						continue
					}
					for _, e := range results[gi].edges {
						removed[e] = true
					}
					changed = true
				}
				for _, e := range results[best].edges {
					protected[e] = true
				}
			}
		}
		if !changed {
			break
		}
	}
}

type bubbleKey struct {
	start kmer
	end   kmer
}

// collapseBubbles removes all but the highest-coverage unitig among short
// parallel paths that share both endpoints. Error-induced shifts produce such
// bubbles; true long repeats do not reconverge inside the length limit.
func collapseBubbles(counts *countedKMers, k int, removed, protected []bool, maxLength int) {
	for round := 0; round < 2; round++ {
		units := buildUnitigs(counts, k, removed, false)
		groups := make(map[bubbleKey][]int)
		for i := range units {
			unit := units[i]
			if unit.closed || len(unit.edges) == 0 || len(unit.edges) > maxLength {
				continue
			}
			start := counts.keys[unit.edges[0]].prefix()
			end := counts.keys[unit.edges[len(unit.edges)-1]].suffix(k)
			key := bubbleKey{start: start, end: end}
			groups[key] = append(groups[key], i)
		}
		changed := false
		for _, members := range groups {
			if len(members) < 2 {
				continue
			}
			best := -1
			bestCoverage := -1.0
			bestCanon := ""
			for _, i := range members {
				if units[i].coverage > bestCoverage {
					bestCoverage = units[i].coverage
					best = i
					bestCanon = ""
				} else if units[i].coverage == bestCoverage && best >= 0 {
					if bestCanon == "" {
						s := unitigString(counts, units[best], k)
						r := string(reverseComplementSequence([]byte(s)))
						bestCanon = s
						if r < bestCanon {
							bestCanon = r
						}
					}
					s := unitigString(counts, units[i], k)
					r := string(reverseComplementSequence([]byte(s)))
					c := s
					if r < c {
						c = r
					}
					if c < bestCanon {
						bestCanon = c
						best = i
					}
				}
			}
			if best < 0 {
				continue
			}
			for _, i := range members {
				if i == best {
					continue
				}
				loser := true
				for _, edge := range units[i].edges {
					if protected[edge] {
						loser = false
						break
					}
				}
				if !loser {
					continue
				}
				for _, edge := range units[i].edges {
					if !protected[edge] {
						removed[edge] = true
					}
				}
				changed = true
			}
			for _, edge := range units[best].edges {
				protected[edge] = true
			}
		}
		if !changed {
			break
		}
	}
}

func reverseComplementSequence(read []byte) []byte {
	result := make([]byte, len(read))
	for i, base := range read {
		var complement byte
		switch base {
		case 'A':
			complement = 'T'
		case 'C':
			complement = 'G'
		case 'G':
			complement = 'C'
		case 'T':
			complement = 'A'
		default:
			complement = 'N'
		}
		result[len(read)-1-i] = complement
	}
	return result
}

func isDoubleStranded(dataset *assemblerDataset, counts *countedKMers, k int) bool {
	if len(dataset.reads) == 0 {
		return false
	}
	step := len(dataset.reads) / 1000
	if step < 1 {
		step = 1
	}
	total := 0
	present := 0
	for index := 0; index < len(dataset.reads); index += step {
		read := reverseComplementSequence(dataset.reads[index])
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			total++
			if counts.count(code) > 0 {
				present++
			}
		}
	}
	if total == 0 {
		return false
	}
	return float64(present)/float64(total) > 0.20
}

func canonicalContigs(contigs []string) []string {
	unique := make(map[string]string, len(contigs))
	for _, contig := range contigs {
		if contig == "" {
			continue
		}
		key := canonicalCircular(contig)
		if existing, ok := unique[key]; !ok || contig < existing {
			unique[key] = contig
		}
	}
	result := make([]string, 0, len(unique))
	for _, contig := range unique {
		result = append(result, contig)
	}
	sort.Strings(result)
	return result
}

func medianReadLength(dataset *assemblerDataset) int {
	lengths := make([]int, 0, len(dataset.reads))
	for _, read := range dataset.reads {
		if len(read) > 0 {
			lengths = append(lengths, len(read))
		}
	}
	if len(lengths) == 0 {
		return 0
	}
	sort.Ints(lengths)
	return lengths[len(lengths)/2]
}

func orientPairedReads(dataset *assemblerDataset) {
	if !dataset.paired || len(dataset.pairs) == 0 {
		return
	}
	const k = 31
	limit := len(dataset.pairs)
	if limit > 30000 {
		limit = 30000
	}
	seen := newKmerBloom(25)
	for _, pair := range dataset.pairs[:limit] {
		scanReadKMers(pair.first, k, true, func(code, _ kmer) bool {
			seen.add(code)
			return true
		})
	}
	forward, reverse, total := 0, 0, 0
	for _, pair := range dataset.pairs[:limit] {
		scanReadKMers(pair.second, k, true, func(code, _ kmer) bool {
			total++
			if seen.contains(code) {
				forward++
			}
			return true
		})
		rc := reverseComplementSequence(pair.second)
		scanReadKMers(rc, k, true, func(code, _ kmer) bool {
			if seen.contains(code) {
				reverse++
			}
			return true
		})
	}
	if total == 0 || reverse <= 3*forward || reverse < 1000 {
		return
	}
	for index := range dataset.pairs {
		rc := reverseComplementSequence(dataset.pairs[index].second)
		copy(dataset.pairs[index].second, rc)
	}
}

func detectDoubleStranded(dataset *assemblerDataset, k int) bool {
	if len(dataset.reads) < 2 || k < 2 {
		return false
	}
	indexDS := sampleDataset(dataset, 200000)
	seen := newKmerBloom(26)
	eachKmerMode(indexDS, k, false, false, func(code kmer) {
		seen.add(code)
	})
	queryDS := sampleDataset(dataset, 2000)
	total := 0
	present := 0
	for _, read := range queryDS.reads {
		scanReadKMers(read, k, true, func(_, rc kmer) bool {
			total++
			if seen.contains(rc) {
				present++
			}
			return true
		})
	}
	if total == 0 {
		return false
	}
	return float64(present)/float64(total) > 0.20
}

func computeKBounds(median, cap int) (int, int) {
	kMax := median - 1
	if kMax > 61 {
		kMax = 61
	}
	if kMax < 3 {
		kMax = 3
	}
	if kMax > median {
		kMax = median
	}
	if cap > 0 && kMax > cap {
		kMax = cap
	}
	detectionK := 31
	if detectionK > kMax {
		detectionK = kMax
	}
	if detectionK > median {
		detectionK = median
	}
	if median < 95 && detectionK > 23 {
		detectionK = 23
	}
	return kMax, detectionK
}

func maybeCorrectDataset(dataset *assemblerDataset, median, kMax, detectionK int, errorFraction float64) (float64, bool) {
	if errorFraction <= 0.10 || len(dataset.reads) < 500 {
		return errorFraction, false
	}
	correctionK := 29
	if median < 95 {
		correctionK = 23
	}
	if correctionK > kMax {
		correctionK = kMax
	}
	if correctionK > median {
		correctionK = median
	}
	if correctionK < 2 {
		return errorFraction, true
	}
	solid := solidKMersForCorrection(dataset, correctionK, 3)
	bloom := buildKmerBloom(solid.keys)
	correctReads(dataset, correctionK, bloom)
	solid = nil
	bloom = nil
	runtime.GC()
	debug.FreeOSMemory()
	return detectErrorFraction(dataset, detectionK), true
}

func chooseAssemblyK(readsCount, median, kMax, detectionK int, errorFraction float64, correctionRan bool) int {
	chosenK := kMax
	if errorFraction > 0.10 {
		chosenK = 31
		if kMax >= 37 && median >= 37 {
			chosenK = 37
		}
		// Error-heavy data (indels shred reads into short clean stretches):
		// shrink k so a solid k-mer stays likely inside one clean stretch.
		if detectionK > 1 {
			perBase := 1 - math.Pow(1-errorFraction, 1/float64(detectionK))
			if perBase > 1e-4 {
				kErr := int(math.Log(0.45) / math.Log(1-perBase))
				if kErr < 3 {
					kErr = 3
				}
				if kErr < 21 {
					kErr = 21
				}
				if kErr < chosenK {
					chosenK = kErr
				}
			}
		}
	} else if correctionRan {
		if errorFraction <= 0.02 {
			if readsCount < 50000 && chosenK > 51 {
				chosenK = 51
			}
		} else if median < 80 && chosenK > 27 {
			chosenK = 27
		} else if chosenK > 35 {
			chosenK = 35
		}
	}
	if median < 95 && chosenK > median*2/3 && median*2/3 >= 21 {
		chosenK = median * 2 / 3
		if chosenK%2 == 0 {
			chosenK--
		}
	}
	if chosenK > kMax {
		chosenK = kMax
	}
	if chosenK > median {
		chosenK = median
	}
	if chosenK < 3 {
		chosenK = kMax
	}
	return chosenK
}

func assembleMultiKErrorHeavy(dataset *assemblerDataset, kMax, fallbackK int, bothStrands bool) ([]string, *countedKMers, int) {
	bestN50, bestBig := -1, -1
	chosenK := fallbackK
	var counts *countedKMers
	var contigs []string
	for cand := 21; cand <= 37 && cand <= kMax; cand += 6 {
		candidate := countKMers(dataset, cand, false, bothStrands, false)
		candidate.fillMinSolid = 2
		fRatio := 0.35
		if candidate.estimatedCoverage >= 18 {
			fRatio = 0.40
		}
		assembly := assembleAtK(candidate, assemblyConfig{
			k: cand, minCount: 2, filter: true, filterRatio: fRatio, resolveForks: true, bothStrands: bothStrands,
		})
		n50 := bigN50(assembly)
		big := bigBases(assembly)
		if n50 > bestN50 || (n50 == bestN50 && big > bestBig) {
			bestN50, bestBig = n50, big
			counts, contigs, chosenK = candidate, assembly, cand
		}
	}
	dataset.reads = nil
	runtime.GC()
	debug.FreeOSMemory()
	if len(contigs) == 0 {
		contigs = assembleAtK(counts, assemblyConfig{k: chosenK, minCount: 1, bothStrands: bothStrands})
	}
	return contigs, counts, chosenK
}

func assembleDataset(dataset *assemblerDataset) ([]string, *countedKMers, int) {
	if len(dataset.reads) == 0 {
		return nil, nil, 0
	}
	orientPairedReads(dataset)
	median := medianReadLength(dataset)
	if median < 2 {
		return nil, nil, 0
	}
	kMax, detectionK := computeKBounds(median, dataset.kMaxCap)
	if kMax < 2 {
		return nil, nil, 0
	}
	if len(dataset.reads) < 100 {
		reads := make([]string, len(dataset.reads))
		for i, read := range dataset.reads {
			reads[i] = string(read)
		}
		return UnitigContigs(collectKMers(reads, kMax, true)), nil, kMax
	}
	bothStrands := detectionK >= 2 && detectDoubleStranded(dataset, detectionK)
	errorFraction := 0.0
	if detectionK >= 2 {
		errorFraction = detectErrorFraction(dataset, detectionK)
	}
	errorFraction, correctionRan := maybeCorrectDataset(dataset, median, kMax, detectionK, errorFraction)
	runtime.GC()
	debug.FreeOSMemory()
	sample := sampleDataset(dataset, 1000)
	chosenK := chooseAssemblyK(len(dataset.reads), median, kMax, detectionK, errorFraction, correctionRan)
	if chosenK < 2 {
		return nil, nil, 0
	}
	if errorFraction > 0.10 && !dataset.paired && len(dataset.reads) > 1000000 {
		dataset.reads = dataset.reads[:1000000]
	}
	var counts *countedKMers
	var contigs []string
	if errorFraction <= 0.03 || (!correctionRan && errorFraction <= 0.10) {
		counts = countDirectLargeMode(dataset, chosenK, false, bothStrands, kmerInstances(dataset, chosenK), true)
		if correctionRan {
			counts.fillMinSolid = 2
		} else {
			counts.fillMinSolid = 3
		}
	} else if errorFraction <= 0.10 {
		counts = countKMers(dataset, chosenK, false, bothStrands, true)
		counts.fillMinSolid = 2
	} else if kMax >= 25 && len(dataset.reads) <= 400000 {
		contigs, counts, chosenK = assembleMultiKErrorHeavy(dataset, kMax, chosenK, bothStrands)
	} else {
		counts = countKMers(dataset, chosenK, false, bothStrands, true)
		counts.fillMinSolid = 2
	}
	canonical := bothStrands || isDoubleStranded(sample, counts, chosenK)
	if contigs == nil {
		filterRatio := 0.30
		if errorFraction > 0.50 {
			filterRatio = 0.40
		}
		contigs = assembleAtK(counts, assemblyConfig{
			k: chosenK, minCount: 2, filter: true, filterRatio: filterRatio,
			resolveForks: errorFraction > 0.03 || dataset.resolveForksAlways, bothStrands: bothStrands,
		})
		if len(contigs) == 0 {
			contigs = assembleAtK(counts, assemblyConfig{k: chosenK, minCount: 1, bothStrands: bothStrands})
		}
	}
	if canonical {
		return canonicalContigs(contigs), counts, chosenK
	}
	return contigs, counts, chosenK
}

func detectErrorFraction(dataset *assemblerDataset, k int) float64 {
	seen := newKmerBloom(28)
	twice := newKmerBloom(28)
	thrice := newKmerBloom(28)
	limit := len(dataset.reads)
	if limit > 600000 {
		limit = 600000
	}
	for _, read := range dataset.reads[:limit] {
		scanReadKMers(read, k, true, func(code, rc kmer) bool {
			canon := canonPair(code, rc)
			if twice.contains(canon) {
				thrice.add(canon)
			}
			if seen.contains(canon) {
				twice.add(canon)
			}
			seen.add(canon)
			return true
		})
	}
	step := limit / 1000
	if step < 1 {
		step = 1
	}
	absent := 0
	total := 0
	for index := 0; index < limit; index += step {
		scanReadKMers(dataset.reads[index], k, true, func(code, rc kmer) bool {
			total++
			if !thrice.contains(canonPair(code, rc)) {
				absent++
			}
			return true
		})
	}
	if total == 0 {
		return 0
	}
	return float64(absent) / float64(total)
}

func solidKMersForCorrection(dataset *assemblerDataset, k int, minCount uint32) *countedKMers {
	total := kmerInstances(dataset, k)
	if total < 500000 || minCount <= 1 {
		table := newKmerTable(1024)
		eachKmer(dataset, k, true, func(code kmer) {
			index := table.add(code)
			table.counts[index]++
		})
		return extractCounts(table, minCount, total, k)
	}
	counts := countBloomMode(dataset, k, true, false, total, false)
	if minCount <= 2 {
		return counts
	}
	filteredKeys := make([]kmer, 0, len(counts.keys))
	filteredCounts := make([]uint32, 0, len(counts.counts))
	for i, c := range counts.counts {
		if c >= minCount {
			filteredKeys = append(filteredKeys, counts.keys[i])
			filteredCounts = append(filteredCounts, c)
		}
	}
	counts.keys = filteredKeys
	counts.counts = filteredCounts
	return counts
}

func buildKmerBloom(keys []kmer) *kmerBloom {
	bloom := newKmerBloom(30)
	for _, code := range keys {
		bloom.add(code)
	}
	return bloom
}

func correctReads(dataset *assemblerDataset, k int, bloom *kmerBloom) (int, int) {
	corrected := 0
	attempted := 0
	for _, read := range dataset.reads {
		if len(read) < k {
			continue
		}
		fixes := 0
		var code, rc kmer
		run := 0
		i := 0
		for i < len(read) && fixes < 4 {
			v, ok := baseValue(read[i])
			if !ok {
				run = 0
				code = kmer{}
				rc = kmer{}
				i++
				continue
			}
			code = code.appendCode(v, k)
			rc = prependCode(rc.prefix(), 3-v, k)
			run++
			i++
			if run < k {
				continue
			}
			canon := code
			if kmerLess(rc, code) {
				canon = rc
			}
			if bloom.contains(canon) {
				continue
			}
			start := i - k
			attempted++
			if correctReadAt(read, start, k, bloom) {
				corrected++
				fixes++
				i = start + 1
				run = 0
				code = kmer{}
				rc = kmer{}
			} else {
				i = start + k/2
				if i < start+1 {
					i = start + 1
				}
				run = 0
				code = kmer{}
				rc = kmer{}
			}
		}
	}
	return corrected, attempted
}

func tryMutatePosition(read []byte, position, k int, bloom *kmerBloom) bool {
	original := read[position]
	for _, base := range []byte("ACGT") {
		if base == original {
			continue
		}
		read[position] = base
		if localKmersSolid(read, position, k, bloom) {
			return true
		}
	}
	read[position] = original
	return false
}

func correctReadAt(read []byte, bad, k int, bloom *kmerBloom) bool {
	end := bad + k - 1
	if end >= len(read) {
		end = len(read) - 1
	}
	if bad > 0 && tryMutatePosition(read, end, k, bloom) {
		return true
	}
	for position := bad; position <= end; position++ {
		if bad > 0 && position == end {
			continue
		}
		if tryMutatePosition(read, position, k, bloom) {
			return true
		}
	}
	return false
}

func localKmersSolid(read []byte, position, k int, bloom *kmerBloom) bool {
	low := position - k + 1
	if low < 0 {
		low = 0
	}
	high := position + k - 1
	if high >= len(read) {
		high = len(read) - 1
	}
	if high-low+1 < k {
		return false
	}
	solid := true
	ok := scanReadKMers(read[low:high+1], k, false, func(code, rc kmer) bool {
		if !bloom.contains(canonPair(code, rc)) {
			solid = false
			return false
		}
		return true
	})
	return ok && solid
}

func sampleDataset(dataset *assemblerDataset, limit int) *assemblerDataset {
	step := len(dataset.reads) / limit
	if step < 1 {
		step = 1
	}
	sample := &assemblerDataset{paired: dataset.paired}
	for index := 0; index < len(dataset.reads) && len(sample.reads) < limit; index += step {
		read := make([]byte, len(dataset.reads[index]))
		copy(read, dataset.reads[index])
		sample.reads = append(sample.reads, read)
	}
	return sample
}

type contigHit struct {
	contig int
	start  int
	strand int8
}

func buildContigIndex(contigs []string, k int) map[uint64]contigHit {
	index := make(map[uint64]contigHit)
	ambiguous := make(map[uint64]bool)
	for contigID, contig := range contigs {
		for start := 0; start+k <= len(contig); start++ {
			code, ok := encodeKmer([]byte(contig[start:start+k]), k)
			if !ok {
				continue
			}
			key := code.lo
			if ambiguous[key] {
				continue
			}
			if _, exists := index[key]; exists {
				delete(index, key)
				ambiguous[key] = true
				continue
			}
			index[key] = contigHit{contig: contigID, start: start, strand: 1}
		}
	}
	return index
}

func mapRead(read []byte, index map[uint64]contigHit, k int) (contigHit, bool) {
	if len(read) < k {
		return contigHit{}, false
	}
	votes := make(map[contigHit]int)
	step := k / 4
	if step < 1 {
		step = 1
	}
	for start := 0; start+k <= len(read); start += step {
		code, ok := encodeKmer(read[start:start+k], k)
		if !ok {
			continue
		}
		if hit, found := index[code.lo]; found {
			votes[contigHit{contig: hit.contig, start: hit.start - start, strand: 1}]++
		}
	}
	reverse := reverseComplementSequence(read)
	for start := 0; start+k <= len(reverse); start += step {
		code, ok := encodeKmer(reverse[start:start+k], k)
		if !ok {
			continue
		}
		if hit, found := index[code.lo]; found {
			readStart := hit.start + (len(read) - 1 - start)
			votes[contigHit{contig: hit.contig, start: readStart, strand: -1}]++
		}
	}
	best := contigHit{}
	bestCount := 0
	for hit, count := range votes {
		if count > bestCount {
			best = hit
			bestCount = count
		}
	}
	secondCount := 0
	for hit, count := range votes {
		if hit.contig == best.contig && hit.strand == best.strand && abs(hit.start-best.start) <= 15 {
			continue
		}
		if count > secondCount {
			secondCount = count
		}
	}
	if bestCount < 2 || bestCount < 2*secondCount {
		return contigHit{}, false
	}
	return best, true
}

type contigLink struct {
	from     int
	fromRev  bool
	fromEnd  int8
	to       int
	toRev    bool
	toEnd    int8
	support  int
	gaps     []int
	gapTotal int

	// fillKind: 0 none, 1 graph walk, 2 read threading
	fillKind    uint8
	fill        string
	fillBackoff int // bases trimmed from the left contig's end anchor
	fillJ       int // bases trimmed from the right contig's head anchor
	fillGap     int
}

type linkEnd struct {
	contig int
	end    int8
	rev    bool
}

func validMapping(hit contigHit, contigLength, readLength int) bool {
	if hit.contig < 0 || readLength <= 0 {
		return false
	}
	if hit.strand > 0 {
		return hit.start >= 0 && hit.start+readLength <= contigLength
	}
	return hit.start-readLength+1 >= 0 && hit.start < contigLength
}

type mappedReadPair struct {
	first  contigHit
	second contigHit
	ok     bool
}

type linkKey struct {
	from    int
	fromRev bool
	to      int
	toRev   bool
}

func validScaffoldLink(link *contigLink) (int, bool) {
	if link.support < 3 {
		return 0, false
	}
	sort.Ints(link.gaps)
	medGap := link.gaps[len(link.gaps)/2]
	if medGap < -250 || medGap > 2500 {
		return 0, false
	}
	return medGap, true
}

func collectScaffoldLinks(mapped []mappedReadPair, contigs []string, insert int, fr bool) map[linkKey]*contigLink {
	links := make(map[linkKey]*contigLink)
	for _, item := range mapped {
		if !item.ok || item.first.contig == item.second.contig {
			continue
		}
		from, to := item.first, item.second
		fromRev := from.strand < 0
		toRev := to.strand < 0
		if fr {
			toRev = to.strand > 0
		}
		fromStart := from.start
		if fromRev {
			fromStart = len(contigs[from.contig]) - 1 - fromStart
		}
		toStart := to.start
		if toRev {
			toStart = len(contigs[to.contig]) - 1 - toStart
		}
		gap := insert - (len(contigs[from.contig]) - fromStart) - toStart
		fromEnd := int8(1)
		if fromRev {
			fromEnd = 0
		}
		toEnd := int8(0)
		if toRev {
			toEnd = 1
		}
		addLink := func(f int, fRev bool, fEnd int8, t int, tRev bool, tEnd int8) {
			key := linkKey{from: f, fromRev: fRev, to: t, toRev: tRev}
			entry := links[key]
			if entry == nil {
				entry = &contigLink{
					from: f, fromRev: fRev, fromEnd: fEnd,
					to: t, toRev: tRev, toEnd: tEnd,
				}
				links[key] = entry
			}
			entry.support++
			entry.gaps = append(entry.gaps, gap)
			entry.gapTotal += gap
		}
		addLink(from.contig, fromRev, fromEnd, to.contig, toRev, toEnd)
		addLink(to.contig, !toRev, toEnd, from.contig, !fromRev, fromEnd)
	}
	return links
}

func tallyLinkSupport(links map[linkKey]*contigLink, contigs []string, repeatContig []bool, minLen int, normalizeGap bool) (map[linkEnd]int, map[linkEnd]int, map[linkEnd]int, map[linkEnd]int) {
	update := func(best, second map[linkEnd]int, key linkEnd, support int) {
		if support > best[key] {
			second[key] = best[key]
			best[key] = support
		} else if support > second[key] {
			second[key] = support
		}
	}
	outBest, outSecond := make(map[linkEnd]int), make(map[linkEnd]int)
	inBest, inSecond := make(map[linkEnd]int), make(map[linkEnd]int)
	for _, link := range links {
		if repeatContig != nil && (repeatContig[link.from] || repeatContig[link.to]) {
			continue
		}
		medGap, ok := validScaffoldLink(link)
		if !ok {
			continue
		}
		if normalizeGap {
			link.gapTotal = medGap * link.support
		}
		if len(contigs[link.from]) < minLen || len(contigs[link.to]) < minLen {
			continue
		}
		update(outBest, outSecond, linkEnd{link.from, link.fromEnd, link.fromRev}, link.support)
		update(inBest, inSecond, linkEnd{link.to, link.toEnd, link.toRev}, link.support)
	}
	return outBest, outSecond, inBest, inSecond
}

func findRepeatContigs(contigCount int, outBest, outSecond, inBest, inSecond map[linkEnd]int) []bool {
	repeatContig := make([]bool, contigCount)
	for c := 0; c < contigCount; c++ {
		ambIn, ambOut := false, false
		for _, rev := range []bool{false, true} {
			for _, end := range []int8{0, 1} {
				k := linkEnd{c, end, rev}
				if inSecond[k] >= 3 && inBest[k] < 2*inSecond[k] {
					ambIn = true
				}
				if outSecond[k] >= 3 && outBest[k] < 2*outSecond[k] {
					ambOut = true
				}
			}
		}
		if ambIn && ambOut {
			repeatContig[c] = true
		}
	}
	return repeatContig
}

func buildScaffoldLinks(mapped []mappedReadPair, contigs []string, insert int, fr bool) []contigLink {
	links := collectScaffoldLinks(mapped, contigs, insert, fr)
	outBest, outSecond, inBest, inSecond := tallyLinkSupport(links, contigs, nil, 250, true)
	repeatContig := findRepeatContigs(len(contigs), outBest, outSecond, inBest, inSecond)
	outBest, outSecond, inBest, inSecond = tallyLinkSupport(links, contigs, repeatContig, 0, false)

	selected := make([]contigLink, 0, len(links))
	for _, link := range links {
		if repeatContig[link.from] || repeatContig[link.to] {
			continue
		}
		if _, ok := validScaffoldLink(link); !ok {
			continue
		}
		fromKey := linkEnd{link.from, link.fromEnd, link.fromRev}
		toKey := linkEnd{link.to, link.toEnd, link.toRev}
		if link.support != outBest[fromKey] || link.support != inBest[toKey] {
			continue
		}
		if (outSecond[fromKey] >= 2 && link.support < 2*outSecond[fromKey]) || (inSecond[toKey] >= 2 && link.support < 2*inSecond[toKey]) {
			continue
		}
		selected = append(selected, *link)
	}
	sort.Slice(selected, func(i, j int) bool {
		if selected[i].support != selected[j].support {
			return selected[i].support > selected[j].support
		}
		gi := selected[i].gapTotal / selected[i].support
		gj := selected[j].gapTotal / selected[j].support
		if gi != gj {
			return gi < gj
		}
		if selected[i].from != selected[j].from {
			return selected[i].from < selected[j].from
		}
		if selected[i].to != selected[j].to {
			return selected[i].to < selected[j].to
		}
		if selected[i].fromRev != selected[j].fromRev {
			return !selected[i].fromRev
		}
		return !selected[i].toRev && selected[j].toRev
	})
	return selected
}

func joinContigPair(from, to string, gap, seedK int) (string, bool) {
	if gap < 0 {
		overlap := -gap
		if overlap <= len(from) && overlap <= len(to) && overlap >= seedK && strings.HasSuffix(from, to[:overlap]) {
			return from + to[overlap:], true
		}
	}
	if gap >= -120 && gap <= 40 {
		limit := 200
		if len(from) < limit {
			limit = len(from)
		}
		if len(to) < limit {
			limit = len(to)
		}
		for ov := limit; ov >= seedK; ov-- {
			if abs(ov+gap) <= 45 && strings.HasSuffix(from, to[:ov]) {
				return from + to[ov:], true
			}
		}
	}
	// Positive gaps are only joined with sequence found by uniqueFillBetween;
	// N-filled joins never improve NGA50 and risk misassemblies.
	return "", false
}

func orientedContig(contigs []string, contig int, reverse bool) string {
	sequence := contigs[contig]
	if reverse {
		return string(reverseComplementSequence([]byte(sequence)))
	}
	return sequence
}

const fillSlack = 120

const fillBackoff = 200

// solidAnchorTail returns the start of the deepest (k-1)-mer anchor inside the
// tail of text that has at least one solid outgoing edge, or -1 when none
// exists.
func solidAnchorTail(counts *countedKMers, k int, minSolid uint32, text string) int {
	low := len(text) - (k - 1) - fillBackoff
	if low < 0 {
		low = 0
	}
	for i := len(text) - (k - 1); i >= low; i-- {
		vertex, ok := encodeKmer([]byte(text[i:i+k-1]), k-1)
		if !ok {
			continue
		}
		for base := byte(0); base < 4; base++ {
			if counts.count(vertex.appendCode(base, k)) >= minSolid {
				return i
			}
		}
	}
	return -1
}

// solidAnchorHead returns the start of the shallowest (k-1)-mer anchor inside
// the head of text that has at least one solid incoming edge, or -1 when none
// exists.
func solidAnchorHead(counts *countedKMers, k int, minSolid uint32, text string) int {
	high := fillBackoff
	if high > len(text)-(k-1) {
		high = len(text) - (k - 1)
	}
	for j := 0; j <= high; j++ {
		vertex, ok := encodeKmer([]byte(text[j:j+k-1]), k-1)
		if !ok {
			continue
		}
		for base := byte(0); base < 4; base++ {
			if counts.count(prependCode(vertex, base, k)) >= minSolid {
				return j
			}
		}
	}
	return -1
}

// uniqueFillBetween searches the solid k-mer graph for paths from the source
// vertex (the tail of the left contig) to the target vertex (the head of the
// right contig). It returns the gap sequence only when exactly one such path
// of the expected length exists; ambiguous or missing paths return false.
func uniqueFillBetween(counts *countedKMers, k int, minSolid uint32, source, target kmer, gap int) (string, bool) {
	if counts == nil || k < 2 {
		return "", false
	}
	want := gap + k - 1
	if want < 1 {
		want = 1
	}
	minEdges := want - fillSlack
	if minEdges < 1 {
		minEdges = 1
	}
	maxEdges := want + fillSlack
	if maxEdges > 4000 {
		maxEdges = 4000
	}
	budget := int64(50000)
	solutions := make([]string, 0, 2)
	path := make([]byte, 0, maxEdges)
	var search func(v kmer, depth int)
	search = func(v kmer, depth int) {
		if len(solutions) > 1 || budget <= 0 {
			return
		}
		budget--
		if depth >= maxEdges {
			return
		}
		for base := byte(0); base < 4; base++ {
			edge := v.appendCode(base, k)
			if counts.count(edge) < minSolid {
				continue
			}
			if edge.suffix(k) == target {
				if depth+1 >= minEdges && depth+1 <= maxEdges {
					spelled := make([]byte, 0, depth+1)
					spelled = append(spelled, path...)
					spelled = append(spelled, base)
					solutions = append(solutions, string(spelled))
				}
				continue
			}
			path = append(path, base)
			search(edge.suffix(k), depth+1)
			path = path[:len(path)-1]
			if len(solutions) > 1 || budget <= 0 {
				return
			}
		}
	}
	search(source, 0)
	if len(solutions) != 1 {
		return "", false
	}
	bases := solutions[0]
	spelled := make([]byte, len(bases))
	for i, code := range []byte(bases) {
		spelled[i] = "ACGT"[code]
	}
	return string(spelled), true
}

func mergeScaffoldLinksChain(contigs []string, links []contigLink, seedK int, counts *countedKMers, k int) []string {
	if len(contigs) == 0 {
		return nil
	}
	n := len(contigs)
	outEdge := make([]int, 2*n)
	inEdge := make([]int, 2*n)
	for i := range outEdge {
		outEdge[i] = -1
		inEdge[i] = -1
	}
	edges := make([]contigLink, 0, len(links))
	for _, link := range links {
		fromNode := 2 * link.from
		if link.fromRev {
			fromNode++
		}
		toNode := 2 * link.to
		if link.toRev {
			toNode++
		}
		if fromNode == toNode {
			continue
		}
		if outEdge[fromNode] >= 0 || inEdge[toNode] >= 0 {
			continue
		}
		outEdge[fromNode] = len(edges)
		inEdge[toNode] = len(edges)
		edges = append(edges, link)
	}
	used := make([]bool, n)
	result := make([]string, 0, n)
	emitChain := func(startNode int) {
		nodes := make([]int, 0, 8)
		node := startNode
		for node >= 0 && !used[node/2] {
			used[node/2] = true
			nodes = append(nodes, node)
			edgeIndex := outEdge[node]
			if edgeIndex < 0 {
				break
			}
			next := 2 * edges[edgeIndex].to
			if edges[edgeIndex].toRev {
				next++
			}
			if used[next/2] {
				break
			}
			node = next
		}
		if len(nodes) == 0 {
			return
		}
		text := orientedContig(contigs, nodes[0]/2, nodes[0]%2 == 1)
		for index := 1; index < len(nodes); index++ {
			next := orientedContig(contigs, nodes[index]/2, nodes[index]%2 == 1)
			edgeIndex := outEdge[nodes[index-1]]
			gap := 0
			if edgeIndex >= 0 {
				gap = edges[edgeIndex].gapTotal / edges[edgeIndex].support
			}
			merged, ok := joinContigPair(text, next, gap, seedK)
			if !ok && edgeIndex >= 0 {
				link := &edges[edgeIndex]
				switch link.fillKind {
				case 1:
					if k >= 2 && len(text) >= link.fillBackoff && len(next) >= link.fillJ+k-1 {
						cut := len(text) - link.fillBackoff
						merged = text[:cut] + link.fill + next[link.fillJ+k-1:]
						ok = true
					}
				case 2:
					if link.fillGap > 0 {
						merged = text + link.fill + next
						ok = true
					} else if overlap := -link.fillGap; overlap < len(next) {
						merged = text + next[overlap:]
						ok = true
					}
				}
			}
			if !ok {
				result = append(result, text)
				text = next
				continue
			}
			text = merged
		}
		result = append(result, text)
	}
	for node := 0; node < 2*n; node++ {
		if used[node/2] || inEdge[node] >= 0 || outEdge[node] < 0 {
			continue
		}
		emitChain(node)
	}
	for node := 0; node < 2*n; node++ {
		if used[node/2] || outEdge[node] < 0 {
			continue
		}
		emitChain(node)
	}
	for index, contig := range contigs {
		if !used[index] {
			result = append(result, contig)
		}
	}
	sort.Strings(result)
	return result
}

func scaffoldContigs(contigs []string, pairs []readPair, seedK int, counts *countedKMers, k int) []string {
	if len(contigs) < 2 || len(pairs) == 0 {
		return contigs
	}
	lengths := make([]int, 0, len(pairs))
	for _, pair := range pairs {
		length := len(pair.first)
		if len(pair.second) < length {
			length = len(pair.second)
		}
		if length > 0 {
			lengths = append(lengths, length)
		}
	}
	if len(lengths) == 0 {
		return contigs
	}
	sort.Ints(lengths)
	median := lengths[len(lengths)/2]
	if seedK > median-1 {
		seedK = median - 1
	}
	if seedK < 8 {
		return contigs
	}
	index := buildContigIndex(contigs, seedK)
	mapped := make([]mappedReadPair, len(pairs))
	sameDistances := make([]int, 0)
	oppositeCount := 0
	sameCount := 0
	for i, pair := range pairs {
		first, okFirst := mapRead(pair.first, index, seedK)
		second, okSecond := mapRead(pair.second, index, seedK)
		if okFirst {
			okFirst = validMapping(first, len(contigs[first.contig]), len(pair.first))
		}
		if okSecond {
			okSecond = validMapping(second, len(contigs[second.contig]), len(pair.second))
		}
		mapped[i] = mappedReadPair{first: first, second: second, ok: okFirst && okSecond}
		if !okFirst || !okSecond {
			continue
		}
		if first.strand != second.strand {
			oppositeCount++
		} else {
			sameCount++
		}
		if first.contig != second.contig {
			continue
		}
		distance := second.start - first.start
		if distance < 0 {
			distance = -distance
		}
		sameDistances = append(sameDistances, distance)
	}
	insert := 0
	if len(sameDistances) > 0 {
		sort.Ints(sameDistances)
		insert = sameDistances[len(sameDistances)/2]
	}
	if insert <= 0 {
		insert = pairs[0].distance
	}
	options := []bool{false, true}
	if oppositeCount+sameCount >= 10 {
		options = []bool{oppositeCount > sameCount}
	}
	bestLinks := []contigLink(nil)
	bestScore := -1
	bestGap := 1 << 60
	for _, fr := range options {
		candidateInsert := insert
		if fr && len(sameDistances) == 0 {
			candidateInsert += median - 1
		}
		links := buildScaffoldLinks(mapped, contigs, candidateInsert, fr)
		score := 0
		gapScore := 0
		for _, link := range links {
			score += link.support
			gap := link.gapTotal / link.support
			if gap < 0 {
				gap = -gap
			}
			gapScore += gap
		}
		if score > bestScore || (score == bestScore && gapScore < bestGap) {
			bestLinks = links
			bestScore = score
			bestGap = gapScore
		}
	}
	computeLinkFills(contigs, bestLinks, counts, k)
	threadLinkFills(contigs, bestLinks, pairs)
	return mergeScaffoldLinksChain(contigs, bestLinks, seedK, counts, k)
}

// computeLinkFills tries a unique solid-graph walk between the oriented ends
// of every selected link and stores the spelled gap sequence on the link.
func computeLinkFills(contigs []string, links []contigLink, counts *countedKMers, k int) {
	if counts == nil || k < 2 {
		return
	}
	minSolid := uint32(3)
	if counts.fillMinSolid != 0 {
		minSolid = counts.fillMinSolid
	} else if counts.estimatedCoverage < 6 {
		minSolid = 2
	}
	for idx := range links {
		link := &links[idx]
		gap := link.gapTotal / link.support
		if gap < -fillSlack {
			continue
		}
		from := orientedContig(contigs, link.from, link.fromRev)
		to := orientedContig(contigs, link.to, link.toRev)
		if len(from) < k-1 || len(to) < k-1 {
			continue
		}
		i := solidAnchorTail(counts, k, minSolid, from)
		j := solidAnchorHead(counts, k, minSolid, to)
		if i < 0 || j < 0 {
			continue
		}
		source, sourceOK := encodeKmer([]byte(from[i:i+k-1]), k-1)
		target, targetOK := encodeKmer([]byte(to[j:j+k-1]), k-1)
		if !sourceOK || !targetOK {
			continue
		}
		adjustedGap := (len(from) - i - (k - 1)) + gap + j
		if fill, filled := uniqueFillBetween(counts, k, minSolid, source, target, adjustedGap); filled {
			link.fillKind = 1
			link.fill = fill
			link.fillBackoff = len(from) - (i + k - 1)
			link.fillJ = j
		}
	}
}

// threadLinkFills recovers junctions the graph walk cannot spell (starved or
// error-damaged k-mer coverage) directly from reads: a read spanning both
// contig ends proves the adjacency and supplies the gap bases.
func threadLinkFills(contigs []string, links []contigLink, pairs []readPair) {
	const anchorLen = 21
	const threadSpan = 120
	type anchorRef struct {
		link int
		head bool
	}
	anchors := make(map[uint64][]anchorRef)
	pending := 0
	for idx := range links {
		link := &links[idx]
		if link.fillKind != 0 {
			continue
		}
		gap := link.gapTotal / link.support
		if gap < -threadSpan || gap > threadSpan {
			continue
		}
		from := orientedContig(contigs, link.from, link.fromRev)
		to := orientedContig(contigs, link.to, link.toRev)
		if len(from) < anchorLen || len(to) < anchorLen {
			continue
		}
		tailCode, okTail := encodeKmer([]byte(from[len(from)-anchorLen:]), anchorLen)
		headCode, okHead := encodeKmer([]byte(to[:anchorLen]), anchorLen)
		if !okTail || !okHead {
			continue
		}
		anchors[tailCode.lo] = append(anchors[tailCode.lo], anchorRef{idx, false})
		anchors[headCode.lo] = append(anchors[headCode.lo], anchorRef{idx, true})
		pending++
	}
	if pending == 0 {
		return
	}
	type voteKey struct {
		link    int
		gapRead int
	}
	votes := make(map[voteKey]map[string]int)
	mask := (uint64(1) << (2 * anchorLen)) - 1
	type hit struct {
		link int
		pos  int
		head bool
	}
	tailHits := make([]hit, 0, 8)
	headHits := make([]hit, 0, 8)
	record := func(read []byte) {
		if len(read) < anchorLen {
			return
		}
		tailHits = tailHits[:0]
		headHits = headHits[:0]
		var code uint64
		run := 0
		for w := 0; w < len(read); w++ {
			value, ok := baseValue(read[w])
			if !ok {
				run = 0
				code = 0
				continue
			}
			code = ((code << 2) | uint64(value)) & mask
			run++
			if run < anchorLen {
				continue
			}
			windowStart := w + 1 - anchorLen
			if refs, found := anchors[code]; found {
				for _, ref := range refs {
					if ref.head {
						headHits = append(headHits, hit{ref.link, windowStart, true})
					} else {
						tailHits = append(tailHits, hit{ref.link, windowStart, false})
					}
				}
			}
		}
		for _, tail := range tailHits {
			for _, head := range headHits {
				if head.link != tail.link {
					continue
				}
				link := &links[tail.link]
				gap := link.gapTotal / link.support
				gapRead := head.pos - tail.pos - anchorLen
				if gapRead < -threadSpan || gapRead > threadSpan {
					continue
				}
				if abs(gapRead-gap) > 60 {
					continue
				}
				fill := ""
				if head.pos > tail.pos+anchorLen {
					fill = string(read[tail.pos+anchorLen : head.pos])
				}
				key := voteKey{tail.link, gapRead}
				if votes[key] == nil {
					votes[key] = make(map[string]int)
				}
				votes[key][fill]++
			}
		}
	}
	for _, pair := range pairs {
		record(pair.first)
		record(pair.second)
		record(reverseComplementSequence(pair.first))
		record(reverseComplementSequence(pair.second))
	}
	for key, fills := range votes {
		total := 0
		bestFill := ""
		bestCount := 0
		for fill, count := range fills {
			total += count
			if count > bestCount || (count == bestCount && fill < bestFill) {
				bestCount = count
				bestFill = fill
			}
		}
		if total >= 2 {
			link := &links[key.link]
			link.fillKind = 2
			link.fill = bestFill
			link.fillGap = key.gapRead
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func formatFASTA(contigs []string) string {
	ordered := append([]string(nil), contigs...)
	sort.Strings(ordered)
	var builder strings.Builder
	index := 0
	for _, contig := range ordered {
		if contig == "" {
			continue
		}
		index++
		fmt.Fprintf(&builder, ">CONTIG%d\n%s\n", index, contig)
	}
	return strings.TrimSuffix(builder.String(), "\n")
}

func collectKMers(reads []string, k int, preserveDuplicates bool) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			kmer := read[i : i+k]
			if !preserveDuplicates {
				if seen[kmer] {
					continue
				}
				seen[kmer] = true
			}
			result = append(result, kmer)
		}
	}
	return result
}

func assembleContigs(reads []string) string {
	dataset := &assemblerDataset{}
	for _, read := range reads {
		dataset.reads = append(dataset.reads, []byte(strings.ToUpper(read)))
	}
	contigs, _, _ := assembleDataset(dataset)
	return formatFASTA(contigs)
}

func readAllStdin() []byte {
	if info, err := os.Stdin.Stat(); err == nil && info.Size() > 0 {
		data := make([]byte, info.Size())
		read, err := io.ReadFull(os.Stdin, data)
		if err == nil || err == io.ErrUnexpectedEOF || err == io.EOF {
			return data[:read]
		}
	}
	data, _ := io.ReadAll(os.Stdin)
	return data
}

// Consensus polish: real IonTorrent-style data carries systematic indel
// errors that no graph filter can undo, so after a draft assembly every read
// is realigned to the draft contigs and replaced by the consensus projection
// of its alignment. Re-assembling the corrected reads yields far longer
// contigs, mirroring what production assemblers do with read correction.

type polishSeed struct {
	contig int32
	pos    int32
}

type polishIndex struct {
	k      int
	refs   map[kmer][]polishSeed
	contig []string
}

func buildPolishIndex(contigs []string, k int) *polishIndex {
	idx := &polishIndex{k: k, refs: make(map[kmer][]polishSeed, len(contigs)*100)}
	idx.contig = contigs
	for ci, contig := range contigs {
		if len(contig) < k {
			continue
		}
		for start := 0; start+k <= len(contig); start++ {
			code, ok := encodeKmer([]byte(contig[start:start+k]), k)
			if !ok {
				continue
			}
			idx.refs[code] = append(idx.refs[code], polishSeed{contig: int32(ci), pos: int32(start)})
		}
	}
	return idx
}

type polishMapping struct {
	contig  int32
	start   int32
	end     int32
	qStart  int32
	qEnd    int32
	forward bool
}

func voteBaseIndex(b byte) int {
	switch b {
	case 'A':
		return 0
	case 'C':
		return 1
	case 'G':
		return 2
	case 'T':
		return 3
	default:
		return -1
	}
}

func consensusBaseFromVotes(v []uint32, p int, fallback byte) byte {
	a, c, g, t := v[p*4], v[p*4+1], v[p*4+2], v[p*4+3]
	best := a
	base := byte('A')
	if c > best {
		best = c
		base = 'C'
	}
	if g > best {
		best = g
		base = 'G'
	}
	if t > best {
		best = t
		base = 'T'
	}
	if best == 0 {
		return fallback
	}
	return base
}

func fillPolishDP(seq []byte, window string, freeLeft, freeRight bool, buffer []int32) ([]int32, int, int, int) {
	L, W := len(seq), len(window)
	width := W + 1
	need := (L + 1) * width
	if len(buffer) < need {
		buffer = make([]int32, need)
	}
	dp := buffer[:need]
	for j := 0; j <= W; j++ {
		dp[j] = 0
	}
	for i := 1; i <= L; i++ {
		if freeLeft {
			dp[i*width] = 0
		} else {
			dp[i*width] = int32(-6 * i)
		}
	}
	for i := 1; i <= L; i++ {
		sb := seq[i-1]
		for j := 1; j <= W; j++ {
			score := int32(-4)
			if sb == window[j-1] {
				score = 3
			}
			best := dp[(i-1)*width+(j-1)] + score
			if gap := dp[(i-1)*width+j] - 6; gap > best {
				best = gap
			}
			if gap := dp[i*width+(j-1)] - 6; gap > best {
				best = gap
			}
			dp[i*width+j] = best
		}
	}
	bestScore, bestI, bestJ := dp[L*width+W], L, W
	for j := 0; j <= W; j++ {
		if dp[L*width+j] > bestScore {
			bestScore, bestI, bestJ = dp[L*width+j], L, j
		}
	}
	if freeRight {
		for i := 1; i <= L; i++ {
			if dp[i*width+W] > bestScore {
				bestScore, bestI, bestJ = dp[i*width+W], i, W
			}
		}
	}
	return buffer, width, bestI, bestJ
}

func tracebackPolishDP(seq []byte, window string, dp []int32, width, bestI, bestJ, cStart int, votes, delVotes []uint32, insVotes map[int32]map[string]int) (int, int, int, int, int) {
	i, j := bestI, bestJ
	matches, firstJ := 0, j
	var insBuf [32]byte
	insLen := 0
	flushIns := func(refPos int32) {
		if insLen > 0 && insVotes != nil {
			var fwd [32]byte
			for k := 0; k < insLen; k++ {
				fwd[k] = insBuf[insLen-1-k]
			}
			m := insVotes[refPos]
			if m == nil {
				m = make(map[string]int, 2)
				insVotes[refPos] = m
			}
			m[string(fwd[:insLen])]++
		}
		insLen = 0
	}
	for i > 0 && j > 0 {
		sb := seq[i-1]
		score := int32(-4)
		if sb == window[j-1] {
			score = 3
		}
		switch {
		case dp[(i-1)*width+(j-1)]+score == dp[i*width+j]:
			flushIns(int32(cStart + j - 1))
			if score == 3 {
				matches++
			}
			if votes != nil {
				if base := voteBaseIndex(sb); base >= 0 {
					votes[(cStart+j-1)*4+base]++
				}
			}
			i--
			j--
			firstJ = j
		case dp[(i-1)*width+j]-6 == dp[i*width+j]:
			if insLen < len(insBuf) {
				insBuf[insLen] = sb
				insLen++
			}
			i--
		default:
			flushIns(int32(cStart + j - 1))
			if delVotes != nil {
				delVotes[cStart+j-1]++
			}
			j--
			firstJ = j
		}
	}
	return firstJ + cStart, bestJ + cStart, i, bestI, matches
}

// polishAlignSpanIndel aligns seq against contig[cStart:cEnd], allowing free
// overhanging query ends at contig boundaries, records substitution, deletion,
// and insertion votes, and returns (refStart, refEnd, queryStart, queryEnd, matches, buffer).
func polishAlignSpanIndel(seq []byte, contig string, cStart, cEnd int, votes, delVotes []uint32, insVotes map[int32]map[string]int, buffer []int32) (int, int, int, int, int, []int32) {
	window := contig[cStart:cEnd]
	if len(seq) == 0 || len(window) == 0 {
		return 0, 0, 0, 0, 0, buffer
	}
	buffer, width, bestI, bestJ := fillPolishDP(seq, window, cStart == 0, cEnd == len(contig), buffer)
	dp := buffer[:(len(seq)+1)*width]
	rStart, rEnd, qStart, qEnd, matches := tracebackPolishDP(seq, window, dp, width, bestI, bestJ, cStart, votes, delVotes, insVotes)
	return rStart, rEnd, qStart, qEnd, matches, buffer
}

// polishAlignSpan aligns seq against contig[cStart:cEnd], allowing free
// overhanging query ends at contig boundaries, records per-position base votes,
// and returns (refStart, refEnd, queryStart, queryEnd, matches, buffer).
func polishAlignSpan(seq []byte, contig string, cStart, cEnd int, votes []uint32, buffer []int32) (int, int, int, int, int, []int32) {
	return polishAlignSpanIndel(seq, contig, cStart, cEnd, votes, nil, nil, buffer)
}

// polishAlign runs a full DP alignment of seq against a contig window and
// records per-position base votes. It returns the aligned reference span.
func polishAlign(seq []byte, contig string, cStart, cEnd int, votes []uint32, buffer []int32) (int, int, int, []int32) {
	jStart, jEnd, _, _, matches, buf := polishAlignSpanIndel(seq, contig, cStart, cEnd, votes, nil, nil, buffer)
	return jStart, jEnd, matches, buf
}

// polishSequence aligns seq (or its reverse complement) to the contig set
// and applies votes; it returns the consensus span for read replacement.
func polishSequence(seq []byte, idx *polishIndex, votes, delVotes [][]uint32, insVotes []map[int32]map[string]int, buffer []int32) (polishMapping, []int32, bool) {
	if len(seq) < idx.k+10 {
		return polishMapping{}, buffer, false
	}
	orient := func(s []byte, forward bool) (polishMapping, []int32, bool) {
		type hitRef struct {
			contig int32
			diag   int32
		}
		clusters := make(map[hitRef]int)
		bestCluster := hitRef{contig: -1}
		bestCount := 0
		for w := 0; w+idx.k <= len(s); w++ {
			code, ok := encodeKmer(s[w:w+idx.k], idx.k)
			if !ok {
				continue
			}
			refs := idx.refs[code]
			if len(refs) == 0 || len(refs) > 8 {
				continue
			}
			for _, ref := range refs {
				key := hitRef{contig: ref.contig, diag: int32(ref.pos) - int32(w)}
				clusters[key]++
				if clusters[key] > bestCount {
					bestCount = clusters[key]
					bestCluster = key
				}
			}
		}
		if bestCluster.contig < 0 || bestCount < 2 {
			return polishMapping{}, buffer, false
		}
		contig := idx.contig[bestCluster.contig]
		// window around the cluster diagonal
		cStart := int(bestCluster.diag) - 30
		if cStart < 0 {
			cStart = 0
		}
		cEnd := int(bestCluster.diag) + len(s) + 60
		if cEnd > len(contig) {
			cEnd = len(contig)
		}
		if cEnd-cStart < 40 {
			return polishMapping{}, buffer, false
		}
		var dv []uint32
		var iv map[int32]map[string]int
		if delVotes != nil {
			dv = delVotes[bestCluster.contig]
		}
		if insVotes != nil {
			iv = insVotes[bestCluster.contig]
		}
		jStart, jEnd, iStart, iEnd, matches, buffer := polishAlignSpanIndel(s, contig, cStart, cEnd, votes[bestCluster.contig], dv, iv, buffer)
		span := jEnd - jStart
		if span < 40 || matches*10 < span*7 {
			return polishMapping{}, buffer, false
		}
		return polishMapping{
			contig:  bestCluster.contig,
			start:   int32(jStart),
			end:     int32(jEnd),
			qStart:  int32(iStart),
			qEnd:    int32(iEnd),
			forward: forward,
		}, buffer, true
	}
	if m, b, ok := orient(seq, true); ok {
		return m, b, true
	}
	rcSeq := reverseComplementSequence(seq)
	return orient(rcSeq, false)
}

func bigBases(contigs []string) int {
	total := 0
	for _, contig := range contigs {
		if len(contig) >= 500 {
			total += len(contig)
		}
	}
	return total
}

// bigN50 is the N50 of contigs of at least 500 bases, a reference-free
// proxy for assembly contiguity.
func bigN50(contigs []string) int {
	lengths := make([]int, 0, len(contigs))
	total := 0
	for _, contig := range contigs {
		if len(contig) >= 500 {
			lengths = append(lengths, len(contig))
			total += len(contig)
		}
	}
	if total == 0 {
		return 0
	}
	sort.Sort(sort.Reverse(sort.IntSlice(lengths)))
	half := total / 2
	cumulative := 0
	for _, length := range lengths {
		cumulative += length
		if cumulative >= half {
			return length
		}
	}
	return 0
}

// buildIndelConsensus constructs the majority-vote consensus sequence per contig
// along with a coordinate map translating old contig offsets into consensus offsets.
func buildIndelConsensus(contigs []string, votes, delVotes [][]uint32, insVotes []map[int32]map[string]int) ([]string, [][]int32) {
	consensus := make([]string, len(contigs))
	coordMap := make([][]int32, len(contigs))
	for ci, contig := range contigs {
		v := votes[ci]
		dv := delVotes[ci]
		iv := insVotes[ci]
		out := make([]byte, 0, len(contig)+16)
		cmap := make([]int32, len(contig)+1)
		for p := 0; p < len(contig); p++ {
			cmap[p] = int32(len(out))
			a, c, g, t := v[p*4], v[p*4+1], v[p*4+2], v[p*4+3]
			d := dv[p]
			depth := a + c + g + t + d
			if !(depth >= 3 && d*2 > depth) {
				out = append(out, consensusBaseFromVotes(v, p, contig[p]))
			}
			if depth >= 3 {
				if insMap := iv[int32(p)]; insMap != nil {
					bestIns := ""
					bestCnt := 0
					for insStr, cnt := range insMap {
						if cnt > bestCnt || (cnt == bestCnt && insStr < bestIns) {
							bestCnt = cnt
							bestIns = insStr
						}
					}
					if bestCnt*2 > int(depth) {
						out = append(out, bestIns...)
					}
				}
			}
		}
		cmap[len(contig)] = int32(len(out))
		consensus[ci] = string(out)
		coordMap[ci] = cmap
	}
	return consensus, coordMap
}

func polishReadsAgainst(reads [][]byte, pairs []readPair, paired bool, contigs []string) *assemblerDataset {
	const seedK = 15
	idx := buildPolishIndex(contigs, seedK)
	votes := make([][]uint32, len(contigs))
	delVotes := make([][]uint32, len(contigs))
	insVotes := make([]map[int32]map[string]int, len(contigs))
	for i, contig := range contigs {
		votes[i] = make([]uint32, len(contig)*4)
		delVotes[i] = make([]uint32, len(contig))
		insVotes[i] = make(map[int32]map[string]int)
	}
	var buffer []int32
	mappings := make([]polishMapping, len(reads))
	mapped := make([]bool, len(reads))
	for i, read := range reads {
		m, b, ok := polishSequence(read, idx, votes, delVotes, insVotes, buffer)
		buffer = b
		if ok {
			mappings[i] = m
			mapped[i] = true
		}
	}
	pairMappingsFirst := make([]polishMapping, len(pairs))
	pairMappedFirst := make([]bool, len(pairs))
	pairMappingsSecond := make([]polishMapping, len(pairs))
	pairMappedSecond := make([]bool, len(pairs))
	if paired && len(reads) == 2*len(pairs) {
		for i := range pairs {
			pairMappingsFirst[i], pairMappedFirst[i] = mappings[2*i], mapped[2*i]
			pairMappingsSecond[i], pairMappedSecond[i] = mappings[2*i+1], mapped[2*i+1]
		}
	} else {
		for i := range pairs {
			m, b, ok := polishSequence(pairs[i].first, idx, votes, delVotes, insVotes, buffer)
			buffer = b
			if ok {
				pairMappingsFirst[i] = m
				pairMappedFirst[i] = true
			}
			m, b, ok = polishSequence(pairs[i].second, idx, votes, delVotes, insVotes, buffer)
			buffer = b
			if ok {
				pairMappingsSecond[i] = m
				pairMappedSecond[i] = true
			}
		}
	}
	consensus, coordMap := buildIndelConsensus(contigs, votes, delVotes, insVotes)
	project := func(m polishMapping, original []byte) []byte {
		seqOrient := original
		if !m.forward {
			seqOrient = reverseComplementSequence(original)
		}
		newStart := coordMap[m.contig][m.start]
		newEnd := coordMap[m.contig][m.end]
		core := consensus[m.contig][newStart:newEnd]
		if len(core) == 0 {
			return original
		}
		combined := make([]byte, 0, int(m.qStart)+len(core)+len(seqOrient)-int(m.qEnd))
		if m.qStart > 0 && int(m.qStart) <= len(seqOrient) {
			combined = append(combined, seqOrient[:m.qStart]...)
		}
		combined = append(combined, core...)
		if int(m.qEnd) >= 0 && int(m.qEnd) < len(seqOrient) {
			combined = append(combined, seqOrient[m.qEnd:]...)
		}
		if !m.forward {
			combined = reverseComplementSequence(combined)
		}
		return combined
	}
	polished := &assemblerDataset{paired: paired}
	for i, read := range reads {
		if mapped[i] {
			polished.reads = append(polished.reads, project(mappings[i], read))
		} else {
			polished.reads = append(polished.reads, read)
		}
	}
	for i := range pairs {
		pair := pairs[i]
		if pairMappedFirst[i] {
			pair.first = project(pairMappingsFirst[i], pair.first)
		}
		if pairMappedSecond[i] {
			pair.second = project(pairMappingsSecond[i], pair.second)
		}
		polished.pairs = append(polished.pairs, pair)
	}
	return polished
}

// maybePolish decides whether a consensus-polish pass is worth running: only
// for affordable datasets whose draft assembly is fragmented relative to the
// coverage-implied genome size.
func maybePolish(reads [][]byte, pairs []readPair, paired bool, draft []string, counts *countedKMers, kmerInstances int) *assemblerDataset {
	if len(draft) == 0 || counts == nil {
		return nil
	}
	if kmerInstances <= 0 || kmerInstances > 80_000_000 || len(reads) < 5000 || len(reads) > 400_000 {
		return nil
	}
	if counts.estimatedCoverage < 4 {
		return nil
	}
	estGenome := float64(kmerInstances) / counts.estimatedCoverage
	if float64(bigBases(draft)) >= 0.6*estGenome && !(len(reads) <= 100_000 && len(draft) > 4 && bigN50(draft) < 35000) {
		return nil
	}
	return polishReadsAgainst(reads, pairs, paired, draft)
}

func assembleAndScaffold(dataset *assemblerDataset) ([]string, *countedKMers, int) {
	contigs, counts, k := assembleDataset(dataset)
	if dataset.paired {
		contigs = scaffoldContigs(contigs, dataset.pairs, 21, counts, k)
	}
	return contigs, counts, k
}

func cloneReads(reads [][]byte) [][]byte {
	copied := make([][]byte, len(reads))
	for i, read := range reads {
		buf := make([]byte, len(read))
		copy(buf, read)
		copied[i] = buf
	}
	return copied
}

func countKMerInstances(reads [][]byte, k int) int {
	total := 0
	for _, read := range reads {
		if len(read) >= k {
			total += len(read) - k + 1
		}
	}
	return total
}

func runPolishLoop(contigs []string, counts *countedKMers, k int, readsSnapshot [][]byte, pairs []readPair, paired bool) []string {
	best := contigs
	bestBig := bigBases(best)
	bestN50 := bigN50(best)
	maxRounds := 3
	if len(readsSnapshot) <= 100_000 {
		maxRounds = 8
	}
	for round := 0; round < maxRounds; round++ {
		polished := maybePolish(readsSnapshot, pairs, paired, best, counts, countKMerInstances(readsSnapshot, k))
		if polished == nil {
			break
		}
		// Corrected reads end at draft contig boundaries, so very large k
		// would starve exactly those junctions; cap the second round and
		// keep indel-fork resolution on for residual raw-read shifts.
		polished.kMaxCap = 37
		if bestN50 < 2000 && len(readsSnapshot) <= 100_000 {
			polished.kMaxCap = 27
		}
		polished.resolveForksAlways = true
		var retryDataset *assemblerDataset
		if len(readsSnapshot) <= 100_000 {
			retryDataset = &assemblerDataset{
				reads:              cloneReads(polished.reads),
				pairs:              polished.pairs,
				paired:             polished.paired,
				resolveForksAlways: true,
			}
		}
		contigs2, counts2, k2 := assembleAndScaffold(polished)
		big2, n50_2 := bigBases(contigs2), bigN50(contigs2)
		// Guard against repeat-shattered or duplicated assemblies that
		// inflate raw base counts while lowering true contiguity. If k jumped
		// above the previous k and shattered shorter reads, retry capped at k.
		if (n50_2 <= bestN50 || big2 < bestBig*95/100) && retryDataset != nil && k2 > k && k >= 21 {
			retryDataset.kMaxCap = k
			contigs2, counts2, k2 = assembleAndScaffold(retryDataset)
			big2, n50_2 = bigBases(contigs2), bigN50(contigs2)
		}
		retryDataset = nil
		if n50_2 <= bestN50 || big2 < bestBig*95/100 {
			break
		}
		best, bestBig, bestN50 = contigs2, big2, n50_2
		counts, k = counts2, k2
	}
	return best
}

func main() {
	debug.SetGCPercent(50)
	debug.SetMemoryLimit(2000 << 20)
	input := readAllStdin()
	dataset := parseAssemblerDataset(input)
	input = nil
	runtime.GC()
	var readsSnapshot [][]byte
	if len(dataset.reads) <= 400_000 {
		readsSnapshot = make([][]byte, len(dataset.reads))
		copy(readsSnapshot, dataset.reads)
	}
	contigs, counts, k := assembleAndScaffold(dataset)
	contigs = runPolishLoop(contigs, counts, k, readsSnapshot, dataset.pairs, dataset.paired)
	if output := formatFASTA(contigs); output != "" {
		fmt.Println(output)
	}
}

// The helpers below are inlined from dsa/course_6_genome_assembly/internal/assembly
// so that this file can be submitted to Coursera as a single Go file.

// KMinusOneMers returns the sorted distinct prefixes and suffixes of length
// k-1 of the given k-mers.
func KMinusOneMers(kmers []string) []string {
	seen := make(map[string]bool)
	for _, kmer := range kmers {
		if len(kmer) < 2 {
			continue
		}
		seen[kmer[:len(kmer)-1]] = true
		seen[kmer[1:]] = true
	}
	result := make([]string, 0, len(seen))
	for mer := range seen {
		result = append(result, mer)
	}
	sort.Strings(result)
	return result
}

// DeBruijnEdges builds edges (prefix, suffix) for each k-mer.
func DeBruijnEdges(kmers, vertices []string) [][2]int {
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		if len(kmer) < 2 {
			continue
		}
		edges = append(edges, [2]int{index[kmer[:len(kmer)-1]], index[kmer[1:]]})
	}
	return edges
}

// UnitigContigs returns the strings spelled by maximal non-branching paths in
// the de Bruijn graph of the given k-mers. Duplicate edges produce duplicate
// contigs.
func UnitigContigs(kmers []string) []string {
	if len(kmers) == 0 {
		return nil
	}
	vertices := KMinusOneMers(kmers)
	index := make(map[string]int, len(vertices))
	for i, vertex := range vertices {
		index[vertex] = i
	}
	edges := make([][2]int, 0, len(kmers))
	for _, kmer := range kmers {
		edges = append(edges, [2]int{index[kmer[:len(kmer)-1]], index[kmer[1:]]})
	}
	inDegree := make([]int, len(vertices))
	outDegree := make([]int, len(vertices))
	for _, edge := range edges {
		outDegree[edge[0]]++
		inDegree[edge[1]]++
	}

	used := make([]bool, len(edges))
	adjacency := make([][]int, len(vertices))
	for id, edge := range edges {
		adjacency[edge[0]] = append(adjacency[edge[0]], id)
	}

	contigs := make([]string, 0)
	for startEdge := 0; startEdge < len(edges); startEdge++ {
		if used[startEdge] {
			continue
		}
		start := edges[startEdge][0]
		if inDegree[start] == 1 && outDegree[start] == 1 {
			continue
		}
		currentEdge := startEdge
		var builder strings.Builder
		builder.WriteString(vertices[start])
		for {
			if used[currentEdge] {
				break
			}
			used[currentEdge] = true
			edge := edges[currentEdge]
			builder.WriteString(vertices[edge[1]][len(vertices[edge[1]])-1:])
			next := edge[1]
			if inDegree[next] != 1 || outDegree[next] != 1 {
				break
			}
			found := -1
			for _, candidate := range adjacency[next] {
				if !used[candidate] {
					found = candidate
					break
				}
			}
			if found < 0 {
				break
			}
			currentEdge = found
		}
		contigs = append(contigs, builder.String())
	}
	// Pure cycles have no vertex with in-degree or out-degree different from
	// one, so they are not discovered by the branch-based pass above.
	for startEdge := 0; startEdge < len(edges); startEdge++ {
		if used[startEdge] {
			continue
		}
		start := edges[startEdge][0]
		var builder strings.Builder
		builder.WriteString(vertices[start])
		currentEdge := startEdge
		lastVertex := start
		for !used[currentEdge] {
			used[currentEdge] = true
			edge := edges[currentEdge]
			builder.WriteString(vertices[edge[1]][len(vertices[edge[1]])-1:])
			lastVertex = edge[1]
			next := -1
			for _, candidate := range adjacency[edge[1]] {
				if !used[candidate] {
					next = candidate
					break
				}
			}
			if next < 0 {
				break
			}
			currentEdge = next
		}
		contig := builder.String()
		if lastVertex == start && len(contig) > len(vertices[start]) {
			contig = contig[:len(contig)-len(vertices[start])]
		}
		contigs = append(contigs, contig)
	}
	sort.Strings(contigs)
	return contigs
}
