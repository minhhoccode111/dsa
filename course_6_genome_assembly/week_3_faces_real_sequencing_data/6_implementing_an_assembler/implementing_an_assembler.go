// go run implementing_an_assembler.go

// Assemble contigs from single reads or read-pairs using de Bruijn graphs.

package main

import (
	"bytes"
	"fmt"
	"io"
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
	reads     [][]byte
	pairs     []readPair
	paired    bool
	corrected bool
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
	if (t.size+1)*20 >= len(t.keys)*19 {
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

func (t *kmerTable) find(key kmer) (int, bool) {
	index, found := t.slot(key)
	if !found {
		return 0, false
	}
	return index, true
}

func countSolidKMers(dataset *assemblerDataset, k int, canonical bool) *countedKMers {
	return countKMers(dataset, k, canonical, false)
}

func countSolidKMersRelease(dataset *assemblerDataset, k int, canonical bool) *countedKMers {
	return countKMers(dataset, k, canonical, true)
}

func countKMers(dataset *assemblerDataset, k int, canonical bool, releaseReads bool) *countedKMers {
	instances := uint64(0)
	for _, read := range dataset.reads {
		if len(read) >= k {
			instances += uint64(len(read) - k + 1)
		}
	}
	if len(dataset.reads) <= 100 || instances < 500000 {
		return countDirect(dataset, k, canonical)
	}
	return countBloom(dataset, k, canonical, instances, releaseReads)
}

func countDirect(dataset *assemblerDataset, k int, canonical bool) *countedKMers {
	table := newKmerTable(1024)
	total := uint64(0)
	for _, read := range dataset.reads {
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			if canonical {
				code = canonicalKmer(code, k)
			}
			index := table.add(code)
			table.counts[index]++
			total++
		}
	}
	return extractCounts(table, 1, total)
}

func countDirectLarge(dataset *assemblerDataset, k int, canonical bool, total uint64, releaseReads bool) *countedKMers {
	initial := 1 << 20
	if total > 20_000_000 {
		initial = 1 << 23
	} else if total > 8_000_000 {
		initial = 1 << 22
	}
	table := newKmerTable(initial)
	eachKmer(dataset, k, canonical, func(code kmer) {
		index := table.add(code)
		table.counts[index]++
	})
	if releaseReads {
		dataset.reads = nil
	}
	runtime.GC()
	debug.FreeOSMemory()
	return extractCounts(table, 1, total)
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

func countBloom(dataset *assemblerDataset, k int, canonical bool, total uint64, releaseReads bool) *countedKMers {
	const bloomBits = uint64(1) << 29
	bloom := make([]uint64, bloomBits/64)
	initial := 1 << 20
	if total > 20_000_000 {
		initial = 1 << 23
	} else if total > 8_000_000 {
		initial = 1 << 22
	}
	candidate := newKmerTable(initial)
	test := func(hash uint64) bool {
		index := hash & (bloomBits - 1)
		return bloom[index>>6]&(uint64(1)<<uint(index&63)) != 0
	}
	set := func(hash uint64) {
		index := hash & (bloomBits - 1)
		bloom[index>>6] |= uint64(1) << uint(index&63)
	}
	eachKmer(dataset, k, canonical, func(code kmer) {
		first := hashKmer(code)
		second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
		if test(first) && test(second) {
			candidate.add(code)
		}
		set(first)
		set(second)
	})
	eachKmer(dataset, k, canonical, func(code kmer) {
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
	return extractCounts(candidate, 2, total)
}

func eachKmer(dataset *assemblerDataset, k int, canonical bool, visit func(kmer)) {
	for _, read := range dataset.reads {
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			if canonical {
				code = canonicalKmer(code, k)
			}
			visit(code)
		}
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

func extractCounts(table *kmerTable, minCount uint32, total uint64) *countedKMers {
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
	}
	return result
}

func (c *countedKMers) index(code kmer) int {
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
	k        int
	minCount uint32
	filter   bool
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

func (c *countedKMers) vertexIn(vertex kmer, k int) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		if c.index(prependCode(vertex, base, k)) >= 0 {
			degree++
		}
	}
	return degree
}

func (c *countedKMers) vertexOut(vertex kmer, k int) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		if c.index(vertex.appendCode(base, k)) >= 0 {
			degree++
		}
	}
	return degree
}

func (c *countedKMers) uniqueOut(vertex kmer, k int) (kmer, bool) {
	var next kmer
	found := false
	for base := byte(0); base < 4; base++ {
		candidate := vertex.appendCode(base, k)
		if c.index(candidate) >= 0 {
			if found {
				return kmer{}, false
			}
			next = candidate
			found = true
		}
	}
	return next, found
}

func (c *countedKMers) vertexInLive(vertex kmer, k int, removed []bool) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		index := c.index(prependCode(vertex, base, k))
		if index >= 0 && !removed[index] {
			degree++
		}
	}
	return degree
}

func (c *countedKMers) vertexOutLive(vertex kmer, k int, removed []bool) int {
	degree := 0
	for base := byte(0); base < 4; base++ {
		index := c.index(vertex.appendCode(base, k))
		if index >= 0 && !removed[index] {
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
		if index >= 0 && !removed[index] {
			if found {
				return kmer{}, false
			}
			next = candidate
			found = true
		}
	}
	return next, found
}

func buildUnitigs(counts *countedKMers, k int, removed []bool) []unitig {
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
		destination := key.suffix(k)
		if counts.vertexInLive(source, k, removed) == 1 && counts.vertexOutLive(source, k, removed) == 1 && counts.vertexInLive(destination, k, removed) == 1 {
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
	removed := make([]bool, len(counts.keys))
	filter := cfg.filter && counts.estimatedCoverage >= 4
	maxLength := 2 * cfg.k
	if maxLength < 500 {
		maxLength = 500
	}
	threshold := 0.30 * counts.estimatedCoverage
	if threshold < 3 {
		threshold = 3
	}
	for pass := 0; pass < 20 && filter; pass++ {
		units := buildUnitigs(counts, cfg.k, removed)
		changed := false
		for _, unit := range units {
			if len(unit.edges) == 0 || len(unit.edges) > maxLength {
				continue
			}
			if unit.coverage < threshold {
				for _, edge := range unit.edges {
					removed[edge] = true
				}
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	units := buildUnitigs(counts, cfg.k, removed)
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

func assembleDataset(dataset *assemblerDataset) []string {
	if len(dataset.reads) == 0 {
		return nil
	}
	median := medianReadLength(dataset)
	if median < 2 {
		return nil
	}
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
	if kMax < 2 {
		return nil
	}
	if len(dataset.reads) < 100 {
		reads := make([]string, len(dataset.reads))
		for i, read := range dataset.reads {
			reads[i] = string(read)
		}
		return UnitigContigs(collectKMers(reads, kMax, true))
	}
	detectionK := 31
	if detectionK > kMax {
		detectionK = kMax
	}
	if detectionK > median {
		detectionK = median
	}
	errorFraction := 0.0
	if detectionK >= 2 {
		errorFraction = detectErrorFraction(dataset, detectionK)
	}
	if errorFraction > 0.10 && len(dataset.reads) >= 50000 {
		correctionK := 31
		if correctionK > kMax {
			correctionK = kMax
		}
		if correctionK > median {
			correctionK = median
		}
		if correctionK >= 2 {
			solid := solidKMersForCorrection(dataset, correctionK, 3)
			bloom := buildKmerBloom(solid.keys)
			corrected, attempted := correctReads(dataset, correctionK, bloom)
			solid = nil
			bloom = nil
			runtime.GC()
			debug.FreeOSMemory()
			if attempted == 0 || corrected*2 >= attempted {
				errorFraction = 0
				dataset.corrected = true
			}
		}
	}
	runtime.GC()
	debug.FreeOSMemory()
	sample := sampleDataset(dataset, 1000)
	chosenK := kMax
	if errorFraction > 0.10 {
		chosenK = 31
		if kMax >= 37 && median >= 37 {
			chosenK = 37
		}
	}
	if chosenK > median {
		chosenK = median
	}
	if chosenK < 3 {
		chosenK = kMax
	}
	if chosenK < 2 {
		return nil
	}
	readLimit := 1000000
	if errorFraction > 0.10 && !dataset.paired && len(dataset.reads) > readLimit {
		dataset.reads = dataset.reads[:readLimit]
	}
	var counts *countedKMers
	if errorFraction <= 0.10 {
		counts = countDirectLarge(dataset, chosenK, false, kmerInstances(dataset, chosenK), true)
	} else {
		counts = countSolidKMersRelease(dataset, chosenK, false)
	}
	canonical := isDoubleStranded(sample, counts, chosenK)
	contigs := assembleAtK(counts, assemblyConfig{k: chosenK, minCount: 2, filter: true})
	if len(contigs) == 0 {
		contigs = assembleAtK(counts, assemblyConfig{k: chosenK, minCount: 1})
	}
	if canonical {
		return canonicalContigs(contigs)
	}
	return contigs
}

func detectErrorFraction(dataset *assemblerDataset, k int) float64 {
	const bloomBits = uint64(1) << 28
	seen := make([]uint64, bloomBits/64)
	twice := make([]uint64, bloomBits/64)
	thrice := make([]uint64, bloomBits/64)
	test := func(bits []uint64, hash uint64) bool {
		index := hash & (bloomBits - 1)
		return bits[index>>6]&(uint64(1)<<uint(index&63)) != 0
	}
	set := func(bits []uint64, hash uint64) {
		index := hash & (bloomBits - 1)
		bits[index>>6] |= uint64(1) << uint(index&63)
	}
	limit := len(dataset.reads)
	if limit > 300000 {
		limit = 300000
	}
	for _, read := range dataset.reads[:limit] {
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			first := hashKmer(code)
			second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
			if test(twice, first) && test(twice, second) {
				set(thrice, first)
				set(thrice, second)
			}
			if test(seen, first) && test(seen, second) {
				set(twice, first)
				set(twice, second)
			}
			set(seen, first)
			set(seen, second)
		}
	}
	step := limit / 1000
	if step < 1 {
		step = 1
	}
	absent := 0
	total := 0
	for index := 0; index < limit; index += step {
		read := dataset.reads[index]
		for start := 0; start+k <= len(read); start++ {
			code, ok := encodeKmer(read[start:start+k], k)
			if !ok {
				continue
			}
			first := hashKmer(code)
			second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
			total++
			if !test(thrice, first) || !test(thrice, second) {
				absent++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(absent) / float64(total)
}

func solidKMersForCorrection(dataset *assemblerDataset, k int, minCount uint32) *countedKMers {
	total := kmerInstances(dataset, k)
	initial := 1 << 20
	if total > 20_000_000 {
		initial = 1 << 26
	} else if total > 8_000_000 {
		initial = 1 << 24
	}
	table := newKmerTable(initial)
	eachKmer(dataset, k, false, func(code kmer) {
		index := table.add(code)
		table.counts[index]++
	})
	return extractCounts(table, minCount, total)
}

type kmerBloom struct {
	bits []uint64
}

func buildKmerBloom(keys []kmer) *kmerBloom {
	const bloomBits = uint64(1) << 30
	bloom := &kmerBloom{bits: make([]uint64, bloomBits/64)}
	for _, code := range keys {
		first := hashKmer(code)
		second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
		index := first & (bloomBits - 1)
		bloom.bits[index>>6] |= uint64(1) << uint(index&63)
		index = second & (bloomBits - 1)
		bloom.bits[index>>6] |= uint64(1) << uint(index&63)
	}
	return bloom
}

func (b *kmerBloom) contains(code kmer) bool {
	const bloomBits = uint64(1) << 30
	first := hashKmer(code)
	second := hashKmer(kmer{lo: code.hi + 0x9e3779b97f4a7c15, hi: code.lo})
	index := first & (bloomBits - 1)
	if b.bits[index>>6]&(uint64(1)<<uint(index&63)) == 0 {
		return false
	}
	index = second & (bloomBits - 1)
	return b.bits[index>>6]&(uint64(1)<<uint(index&63)) != 0
}

func correctReads(dataset *assemblerDataset, k int, bloom *kmerBloom) (int, int) {
	corrected := 0
	attempted := 0
	for _, read := range dataset.reads {
		if len(read) < k {
			continue
		}
		for round := 0; round < 3; round++ {
			bad := -1
			for start := 0; start+k <= len(read); start++ {
				code, ok := encodeKmer(read[start:start+k], k)
				if !ok {
					continue
				}
				if !bloom.contains(code) {
					bad = start
					break
				}
			}
			if bad < 0 {
				break
			}
			attempted++
			if !correctReadAt(read, bad, k, bloom) {
				break
			}
			corrected++
		}
	}
	return corrected, attempted
}

func correctReadAt(read []byte, bad, k int, bloom *kmerBloom) bool {
	for position := bad; position < bad+k && position < len(read); position++ {
		original := read[position]
		for _, base := range []byte("ACGT") {
			if base == original {
				continue
			}
			read[position] = base
			if localKmersSolid(read, position, k, bloom) {
				return true
			}
			read[position] = original
		}
	}
	return false
}

func localKmersSolid(read []byte, position, k int, bloom *kmerBloom) bool {
	low := position - k + 1
	if low < 0 {
		low = 0
	}
	high := position
	if high > len(read)-k {
		high = len(read) - k
	}
	for start := low; start <= high; start++ {
		code, ok := encodeKmer(read[start:start+k], k)
		if !ok || !bloom.contains(code) {
			return false
		}
	}
	return true
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
	secondCount := 0
	for hit, count := range votes {
		if count > bestCount {
			secondCount = bestCount
			best = hit
			bestCount = count
		} else if count > secondCount {
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
	gapTotal int
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

func buildScaffoldLinks(mapped []mappedReadPair, contigs []string, insert int, fr bool) []contigLink {
	links := make(map[contigLink]*contigLink)
	for _, item := range mapped {
		if !item.ok || item.first.contig == item.second.contig {
			continue
		}
		from := item.first
		to := item.second
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
		key := contigLink{
			from: from.contig, fromRev: fromRev, fromEnd: fromEnd,
			to: to.contig, toRev: toRev, toEnd: toEnd,
		}
		entry := links[key]
		if entry == nil {
			entry = &contigLink{
				from: from.contig, fromRev: fromRev, fromEnd: fromEnd,
				to: to.contig, toRev: toRev, toEnd: toEnd,
			}
			links[key] = entry
		}
		entry.support++
		entry.gapTotal += gap
	}

	outBest := make(map[linkEnd]int)
	outSecond := make(map[linkEnd]int)
	inBest := make(map[linkEnd]int)
	inSecond := make(map[linkEnd]int)
	update := func(best, second map[linkEnd]int, key linkEnd, support int) {
		if support > best[key] {
			second[key] = best[key]
			best[key] = support
		} else if support > second[key] {
			second[key] = support
		}
	}
	for _, link := range links {
		update(outBest, outSecond, linkEnd{link.from, link.fromEnd, link.fromRev}, link.support)
		update(inBest, inSecond, linkEnd{link.to, link.toEnd, link.toRev}, link.support)
	}

	selected := make([]contigLink, 0, len(links))
	for _, link := range links {
		fromKey := linkEnd{link.from, link.fromEnd, link.fromRev}
		toKey := linkEnd{link.to, link.toEnd, link.toRev}
		if link.support < 3 {
			continue
		}
		if link.support < 2*outSecond[fromKey] || link.support < 2*inSecond[toKey] {
			continue
		}
		if link.support != outBest[fromKey] || link.support != inBest[toKey] {
			continue
		}
		selected = append(selected, *link)
	}
	sort.Slice(selected, func(i, j int) bool {
		if selected[i].support != selected[j].support {
			return selected[i].support > selected[j].support
		}
		return selected[i].gapTotal/selected[i].support < selected[j].gapTotal/selected[j].support
	})
	return selected
}

func mergeScaffoldLinks(contigs []string, links []contigLink, seedK int) []string {
	active := make([]bool, len(contigs))
	sequences := append([]string(nil), contigs...)
	for i := range active {
		active[i] = true
	}
	usedEnds := make([][2]bool, len(contigs))
	for _, link := range links {
		if !active[link.from] || !active[link.to] {
			continue
		}
		if usedEnds[link.from][link.fromEnd] || usedEnds[link.to][link.toEnd] {
			continue
		}
		from := sequences[link.from]
		if link.fromRev {
			from = string(reverseComplementSequence([]byte(from)))
		}
		to := sequences[link.to]
		if link.toRev {
			to = string(reverseComplementSequence([]byte(to)))
		}
		gap := link.gapTotal / link.support
		overlap := MaxOverlap(from, to)
		if gap < 0 && -gap > overlap {
			overlap = -gap
		}
		if overlap > len(from) {
			overlap = len(from)
		}
		if overlap > len(to) {
			overlap = len(to)
		}
		if overlap < 0 {
			overlap = 0
		}
		var merged string
		if overlap >= seedK && strings.HasSuffix(from, to[:overlap]) {
			merged = from + to[overlap:]
		} else {
			gapSize := gap
			if gapSize < 0 {
				gapSize = 0
			}
			if gapSize > 1000 {
				gapSize = 1000
			}
			merged = from + strings.Repeat("N", gapSize) + to
		}
		sequences[link.from] = merged
		usedEnds[link.from][0] = true
		usedEnds[link.from][1] = true
		active[link.to] = false
	}
	result := make([]string, 0, len(contigs))
	for i, sequence := range sequences {
		if active[i] {
			result = append(result, sequence)
		}
	}
	sort.Strings(result)
	return result
}

func scaffoldContigs(contigs []string, pairs []readPair, seedK int) []string {
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
		if !okFirst || !okSecond || first.contig != second.contig {
			continue
		}
		distance := second.start - first.start
		if distance < 0 {
			distance = -distance
		}
		sameDistances = append(sameDistances, distance)
		if first.strand != second.strand {
			oppositeCount++
		} else {
			sameCount++
		}
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
	if oppositeCount+sameCount >= 20 {
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
	return mergeScaffoldLinks(contigs, bestLinks, seedK)
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
	return formatFASTA(assembleDataset(dataset))
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

func main() {
	debug.SetGCPercent(50)
	debug.SetMemoryLimit(420 << 20)
	input := readAllStdin()
	dataset := parseAssemblerDataset(input)
	input = nil
	runtime.GC()
	contigs := assembleDataset(dataset)
	if dataset.paired && !dataset.corrected {
		contigs = scaffoldContigs(contigs, dataset.pairs, 21)
	}
	output := formatFASTA(contigs)
	if output != "" {
		fmt.Println(output)
	}
}

// The helpers below are inlined from dsa/course_6_genome_assembly/internal/assembly
// so that this file can be submitted to Coursera as a single Go file.

// FindEulerianCycle returns an Eulerian cycle as a sequence of vertices. The
// returned path repeats the starting vertex at the end.
func FindEulerianCycle(n int, edges [][2]int) ([]int, bool) {
	if len(edges) == 0 {
		return nil, true
	}

	adjacency := make([][]int, n)
	for id, edge := range edges {
		if edge[0] < 0 || edge[0] >= n || edge[1] < 0 || edge[1] >= n {
			return nil, false
		}
		adjacency[edge[0]] = append(adjacency[edge[0]], id)
	}
	if !hasEulerianCycle(n, edges) {
		return nil, false
	}

	start := -1
	for vertex := 0; vertex < n; vertex++ {
		if len(adjacency[vertex]) > 0 {
			start = vertex
			break
		}
	}
	if start < 0 {
		return nil, false
	}

	stack := []int{start}
	path := make([]int, 0, len(edges)+1)
	for len(stack) > 0 {
		vertex := stack[len(stack)-1]
		if len(adjacency[vertex]) > 0 {
			edgeID := adjacency[vertex][len(adjacency[vertex])-1]
			adjacency[vertex] = adjacency[vertex][:len(adjacency[vertex])-1]
			stack = append(stack, edges[edgeID][1])
		} else {
			path = append(path, vertex)
			stack = stack[:len(stack)-1]
		}
	}
	if len(path) != len(edges)+1 {
		return nil, false
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, true
}

func hasEulerianCycle(n int, edges [][2]int) bool {
	inDegree := make([]int, n)
	outDegree := make([]int, n)
	for _, edge := range edges {
		outDegree[edge[0]]++
		inDegree[edge[1]]++
	}
	for i := 0; i < n; i++ {
		if inDegree[i] != outDegree[i] {
			return false
		}
	}
	return true
}

// MaxOverlap returns the longest suffix of first that is a prefix of second.
func MaxOverlap(first, second string) int {
	limit := len(first)
	if len(second) < limit {
		limit = len(second)
	}
	for size := limit; size > 0; size-- {
		if first[len(first)-size:] == second[:size] {
			return size
		}
	}
	return 0
}

// CircularSuperstring returns the shortest circular superstring of a small set
// of reads, or an empty string when the read count is too large.
func CircularSuperstring(reads []string) string {
	n := len(reads)
	if n == 0 {
		return ""
	}
	if n > 9 {
		return ""
	}
	used := make([]bool, n)
	order := make([]int, 0, n)
	bestLength := 1 << 30
	best := ""
	var search func()
	search = func() {
		if len(order) == n {
			total := 0
			length := 0
			for i, read := range order {
				total += len(reads[read])
				if i > 0 {
					total -= MaxOverlap(reads[order[i-1]], reads[read])
				}
				length = len(reads[read])
			}
			_ = length
			closing := MaxOverlap(reads[order[n-1]], reads[order[0]])
			total -= closing
			if total > bestLength {
				return
			}
			var builder strings.Builder
			builder.WriteString(reads[order[0]])
			for i := 1; i < n; i++ {
				overlap := MaxOverlap(reads[order[i-1]], reads[order[i]])
				builder.WriteString(reads[order[i]][overlap:])
			}
			genome := builder.String()
			if closing > 0 {
				genome = genome[:len(genome)-closing]
			}
			if total < bestLength || (total == bestLength && (best == "" || genome < best)) {
				bestLength = total
				best = genome
			}
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			used[i] = true
			order = append(order, i)
			search()
			order = order[:len(order)-1]
			used[i] = false
		}
	}
	search()
	return best
}

// UniqueKMers returns all distinct k-mers from the reads, sorted.
func UniqueKMers(reads []string, k int) []string {
	seen := make(map[string]bool)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			seen[read[i:i+k]] = true
		}
	}
	result := make([]string, 0, len(seen))
	for kmer := range seen {
		result = append(result, kmer)
	}
	sort.Strings(result)
	return result
}

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
