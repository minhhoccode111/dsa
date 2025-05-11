package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestParseAssemblerDatasetSingle(t *testing.T) {
	ds := parseAssemblerDataset([]byte("2\nACGT\nTGCA\n"))
	if ds.paired || len(ds.reads) != 2 || len(ds.pairs) != 0 {
		t.Fatalf("dataset = %+v", ds)
	}
	if string(ds.reads[0]) != "ACGT" || string(ds.reads[1]) != "TGCA" {
		t.Fatalf("reads = %q", ds.reads)
	}
}

func TestParseAssemblerDatasetPairs(t *testing.T) {
	ds := parseAssemblerDataset([]byte("2\nACGT|TTTT|350\nTGCA|AAAA|300\n"))
	if !ds.paired || len(ds.pairs) != 2 {
		t.Fatalf("dataset = %+v", ds)
	}
	if string(ds.pairs[0].first) != "ACGT" || string(ds.pairs[0].second) != "TTTT" || ds.pairs[0].distance != 350 {
		t.Fatalf("pair = %+v", ds.pairs[0])
	}
	if len(ds.reads) != 4 {
		t.Fatalf("reads = %q", ds.reads)
	}
}

func TestKmerRoundTrip(t *testing.T) {
	for _, k := range []int{3, 21, 31, 32, 37, 55, 61} {
		seq := []byte(strings.Repeat("ACGT", 20)[:k])
		code, ok := encodeKmer(seq, k)
		if !ok {
			t.Fatalf("k=%d encode failed", k)
		}
		if got := string(decodeKmer(code, k)); got != string(seq) {
			t.Fatalf("k=%d decode = %q, want %q", k, got, seq)
		}
	}
}

func TestKmerReverseComplement(t *testing.T) {
	code, _ := encodeKmer([]byte("ACGTAC"), 6)
	rc := code.reverseComplement(6)
	if got := string(decodeKmer(rc, 6)); got != "GTACGT" {
		t.Fatalf("rc = %q", got)
	}
	if canonicalKmer(code, 6) != canonicalKmer(rc, 6) {
		t.Fatal("canonical k-mers differ")
	}
}

func TestKmerPrefixSuffix(t *testing.T) {
	code, _ := encodeKmer([]byte("ACGTA"), 5)
	if got := string(decodeKmer(code.prefix(), 4)); got != "ACGT" {
		t.Fatalf("prefix = %q", got)
	}
	if got := string(decodeKmer(code.suffix(5), 4)); got != "CGTA" {
		t.Fatalf("suffix = %q", got)
	}
	if got := string(decodeKmer(code.appendBase('C', 5), 5)); got != "CGTAC" {
		t.Fatalf("append = %q", got)
	}
}

func TestCountSolidKMers(t *testing.T) {
	ds := &assemblerDataset{reads: [][]byte{[]byte("ACGTACGT"), []byte("ACGTACGT"), []byte("TTTTACGT")}}
	counts := countSolidKMers(ds, 4, false)
	if counts.count(mustKmer(t, "ACGT", 4)) != 5 {
		t.Fatalf("ACGT count = %d", counts.count(mustKmer(t, "ACGT", 4)))
	}
	if counts.count(mustKmer(t, "TTTT", 4)) != 1 {
		t.Fatalf("TTTT count = %d", counts.count(mustKmer(t, "TTTT", 4)))
	}
	if counts.estimatedCoverage <= 0 {
		t.Fatal("coverage not estimated")
	}
}

func TestEstimateErrorFraction(t *testing.T) {
	ds := &assemblerDataset{reads: [][]byte{[]byte("ACGTACGTACGT"), []byte("ACGTACGTACGT")}}
	counts := countSolidKMers(ds, 4, false)
	if got := estimateErrorFraction(ds, counts, 4); got > 0.01 {
		t.Fatalf("error fraction = %v", got)
	}
}

func TestCountBloomPath(t *testing.T) {
	read := strings.Repeat("ACGT", 25)
	reads := make([][]byte, 1500)
	for i := range reads {
		reads[i] = []byte(read)
	}
	ds := &assemblerDataset{reads: reads}
	counts := countSolidKMers(ds, 31, false)
	want := []uint32{27000, 27000, 25500, 25500}
	for start := 0; start < 4; start++ {
		code := mustKmer(t, read[start:start+31], 31)
		if got := counts.count(code); got != want[start] {
			t.Fatalf("start %d count = %d, want %d", start, got, want[start])
		}
	}
}

func TestVisibleDatasetAssembly(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("test", "dataset.txt"))
	if err != nil {
		t.Skip("visible dataset not present")
	}
	ref, err := os.ReadFile(filepath.Join("test", "genome.txt"))
	if err != nil {
		t.Skip("reference genome not present")
	}
	reference := strings.TrimSpace(string(ref))
	ds := parseAssemblerDataset(data)
	contigs, _, _ := assembleDataset(ds)
	if got := referenceCoverage(contigs, reference); got < 0.95 {
		t.Fatalf("coverage = %v", got)
	}
	if got := nga50(contigs, reference); got < len(reference)*8/10 {
		t.Fatalf("NGA50 = %d, want at least %d", got, len(reference)*8/10)
	}
}

func scaffoldCounts(t *testing.T, reads [][]byte, k int) *countedKMers {
	t.Helper()
	if len(reads) == 0 {
		return nil
	}
	lengths := make([]int, 0, len(reads))
	for _, read := range reads {
		if len(read) >= k {
			lengths = append(lengths, len(read))
		}
	}
	if len(lengths) < 100 {
		scaled := make([][]byte, 0, 200)
		for len(scaled) < 200 {
			scaled = append(scaled, reads...)
		}
		reads = scaled
	}
	return countSolidKMers(&assemblerDataset{reads: reads}, k, false)
}

func TestScaffoldContigsFROrientation(t *testing.T) {
	a := randomDNA(1, 600)
	b := randomDNA(2, 600)
	genome := a + b
	ds := frDatasetFor(genome, 300, 100, 350)
	counts := scaffoldCounts(t, ds.reads, 61)
	scaffolds := scaffoldContigs([]string{a, b}, ds.pairs, 21, counts, 61)
	if len(scaffolds) != 1 {
		t.Fatalf("FR scaffolds = %d", len(scaffolds))
	}
	if !strings.Contains(scaffolds[0], a) || !strings.Contains(scaffolds[0], b) {
		t.Fatalf("FR scaffold is not A+B")
	}
}

func TestShortMateDoesNotDisableScaffolding(t *testing.T) {
	uniqueA := randomDNA(1, 300)
	uniqueB := randomDNA(2, 300)
	uniqueC := randomDNA(3, 300)
	repeat := randomDNA(99, 220)
	genome := uniqueA + repeat + uniqueB + repeat + uniqueC
	ds := pairedDatasetFor(genome, 30, 100, 350, 20)
	counts := scaffoldCounts(t, ds.reads, 61)
	base := scaffoldContigs([]string{uniqueA, uniqueB, uniqueC}, ds.pairs, 21, counts, 61)
	ds.pairs = append(ds.pairs, readPair{first: []byte("ACGTACGT"), second: []byte("ACGT"), distance: 350})
	withShort := scaffoldContigs([]string{uniqueA, uniqueB, uniqueC}, ds.pairs, 21, counts, 61)
	if len(withShort) != len(base) {
		t.Fatalf("short mate changed scaffolds: %d -> %d", len(base), len(withShort))
	}
}

func TestScaffoldContigsSpansRepeat(t *testing.T) {
	uniqueA := randomDNA(1, 300)
	uniqueB := randomDNA(2, 300)
	uniqueC := randomDNA(3, 300)
	repeat := randomDNA(99, 220)
	genome := uniqueA + repeat + uniqueB + repeat + uniqueC
	ds := pairedDatasetFor(genome, 30, 100, 350, 20)
	contigs := []string{uniqueA, uniqueB, uniqueC}
	counts := scaffoldCounts(t, ds.reads, 61)
	scaffolds := scaffoldContigs(contigs, ds.pairs, 21, counts, 61)
	if len(scaffolds) >= len(contigs) {
		t.Fatalf("scaffolding did not reduce contig count: %d -> %d", len(contigs), len(scaffolds))
	}
	if referenceCoverage(scaffolds, genome) < 0.6 {
		t.Fatalf("coverage = %v", referenceCoverage(scaffolds, genome))
	}
}

func TestOrientPairedReadsFR(t *testing.T) {
	genome := randomDNA(71, 6000)
	comp := func(read []byte) []byte { return reverseComplementSequence(read) }
	rng := rand.New(rand.NewSource(12))
	ds := &assemblerDataset{paired: true}
	for len(ds.pairs) < 1000 {
		start := rng.Intn(len(genome) - 200)
		distance := 300 + rng.Intn(101)
		if start+distance+100 > len(genome) {
			continue
		}
		first := []byte(genome[start : start+100])
		second := comp([]byte(genome[start+distance : start+distance+100]))
		ds.pairs = append(ds.pairs, readPair{first: first, second: second, distance: 350})
		ds.reads = append(ds.reads, first, second)
	}
	original := append([]byte(nil), ds.pairs[0].second...)
	want := comp(original)
	orientPairedReads(ds)
	if string(ds.pairs[0].second) != string(want) {
		t.Fatal("paired read was not oriented to the forward strand")
	}
	if string(ds.reads[1]) != string(want) {
		t.Fatal("read list was not updated with the oriented mate")
	}
}

func TestMergeScaffoldLinksChain(t *testing.T) {
	a := randomDNA(1, 300)
	b := randomDNA(2, 300)
	c := randomDNA(3, 300)
	genome := a + b + c
	reads := circularReads(genome, 40, 100)
	counts := scaffoldCounts(t, reads, 61)
	links := []contigLink{
		{from: 0, fromEnd: 1, to: 1, toEnd: 0, support: 10, gaps: []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, gapTotal: 50},
		{from: 1, fromEnd: 1, to: 2, toEnd: 0, support: 10, gaps: []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, gapTotal: 50},
	}
	computeLinkFills([]string{a, b, c}, links, counts, 61)
	scaffolds := mergeScaffoldLinksChain([]string{a, b, c}, links, 21, counts, 61)
	if len(scaffolds) != 1 {
		t.Fatalf("chain scaffolds = %d", len(scaffolds))
	}
	for _, contig := range []string{a, b, c} {
		if !strings.Contains(scaffolds[0], contig) {
			t.Fatal("chain scaffold does not contain a contig")
		}
	}
}

func TestCorrectReadsDoubleStranded(t *testing.T) {
	genome := randomDNA(81, 4000)
	reads := linearReadsBothStrands(genome, 40, 100, 5)
	clean := append([]byte(nil), reads[1]...)
	mutated := append([]byte(nil), clean...)
	position := len(mutated) / 2
	mutated[position] = "ACGT"[(strings.IndexByte("ACGT", mutated[position])+1)%4]
	ds := &assemblerDataset{reads: [][]byte{mutated}}
	solid := solidKMersForCorrection(&assemblerDataset{reads: reads}, 31, 3)
	bloom := buildKmerBloom(solid.keys)
	correctReads(ds, 31, bloom)
	if string(ds.reads[0]) != string(clean) {
		t.Fatal("double-stranded read was not corrected")
	}
}

func TestIsDoubleStranded(t *testing.T) {
	genome := randomDNA(42, 200)
	forward := circularReads(genome, 40, 20)
	double := append([][]byte{}, forward...)
	for _, read := range forward {
		double = append(double, reverseComplementBytes(read))
	}
	fwdDS := &assemblerDataset{reads: forward}
	fwdCounts := countSolidKMers(fwdDS, 19, false)
	if isDoubleStranded(fwdDS, fwdCounts, 19) {
		t.Fatal("forward-only data detected as double-stranded")
	}
	doubleDS := &assemblerDataset{reads: double}
	doubleCounts := countSolidKMers(doubleDS, 19, false)
	if !isDoubleStranded(doubleDS, doubleCounts, 19) {
		t.Fatal("double-stranded data not detected")
	}
}

func TestAssembleDatasetDoubleStranded(t *testing.T) {
	genome := randomDNA(42, 200)
	forward := circularReads(genome, 40, 20)
	double := append([][]byte{}, forward...)
	for _, read := range forward {
		double = append(double, reverseComplementBytes(read))
	}
	contigs, _, _ := assembleDataset(&assemblerDataset{reads: double})
	if len(contigs) != 1 || len(contigs[0]) != len(genome) {
		t.Fatalf("contigs = %v", contigs)
	}
}

func TestCanonicalContigs(t *testing.T) {
	got := canonicalContigs([]string{"ACGT", "ACGT", "TTTT"})
	if len(got) != 2 {
		t.Fatalf("contigs = %v", got)
	}
}

func TestCanonicalContigsRotationReverse(t *testing.T) {
	a := "ACGTTGCA"
	rotated := "CGTTGCAA"
	b := string(reverseComplementBytes([]byte(rotated)))
	got := canonicalContigs([]string{a, b})
	if len(got) != 1 {
		t.Fatalf("rotation/reverse contigs = %v", got)
	}
}

func TestAssembleDatasetDoubleStrandedLinear(t *testing.T) {
	genome := randomDNA(77, 1000)
	reads := linearReadsBothStrands(genome, 30, 100, 5)
	contigs, _, _ := assembleDataset(&assemblerDataset{reads: reads})
	if len(contigs) != 1 {
		t.Fatalf("contigs = %d lengths = %v", len(contigs), contigLengths(contigs))
	}
}

func contigLengths(contigs []string) []int {
	lengths := make([]int, len(contigs))
	for i, contig := range contigs {
		lengths[i] = len(contig)
	}
	return lengths
}

func TestCorrectReadSingleError(t *testing.T) {
	genome := randomDNA(11, 2000)
	cleanReads := circularReads(genome, 30, 100)
	clean := append([]byte(nil), cleanReads[0]...)
	solid := solidKMersForCorrection(&assemblerDataset{reads: cleanReads}, 31, 3)
	bloom := buildKmerBloom(solid.keys)
	mutated := append([]byte(nil), clean...)
	mutated[50] = "ACGT"[(strings.IndexByte("ACGT", mutated[50])+1)%4]
	correctReads(&assemblerDataset{reads: [][]byte{mutated}}, 31, bloom)
	if string(mutated) != string(clean) {
		t.Fatalf("corrected read differs from clean read")
	}
}

func TestDirectLargeKeepsSingletons(t *testing.T) {
	read := []byte(randomDNA(31, 100))
	dataset := &assemblerDataset{reads: [][]byte{read}}
	counts := countDirectLarge(dataset, 61, false, 40, false)
	code := mustKmer(t, string(read[:61]), 61)
	if got := counts.count(code); got != 1 {
		t.Fatalf("singleton count = %d, want 1", got)
	}
}

func TestKmerTableGrowsAtCapacity(t *testing.T) {
	table := newKmerTable(16)
	for i := 0; i < 16; i++ {
		table.add(kmer{lo: uint64(i)})
	}
	if len(table.keys) <= 16 {
		t.Fatalf("table capacity = %d, want growth before full", len(table.keys))
	}
}

func TestKmerTableGrows(t *testing.T) {
	table := newKmerTable(16)
	for i := 0; i < 200; i++ {
		table.add(kmer{lo: uint64(i)})
	}
	if len(table.keys) <= 16 {
		t.Fatalf("table capacity = %d, want growth", len(table.keys))
	}
	for i := 0; i < 200; i++ {
		if _, ok := table.find(kmer{lo: uint64(i)}); !ok {
			t.Fatalf("missing key %d", i)
		}
	}
}

func TestDetectErrorFractionCleanVsErrors(t *testing.T) {
	genome := randomDNA(7, 2000)
	clean := detectErrorFraction(&assemblerDataset{reads: circularReads(genome, 30, 100)}, 31)
	if clean > 0.10 {
		t.Fatalf("clean error fraction = %v", clean)
	}
	erroneous := detectErrorFraction(&assemblerDataset{reads: circularReadsWithErrors(genome, 600, 100, 1)}, 31)
	if erroneous < 0.25 {
		t.Fatalf("error-prone fraction = %v", erroneous)
	}
}

func TestAssembleDatasetSingle(t *testing.T) {
	genome := randomDNA(42, 200)
	reads := circularReads(genome, 40, 20)
	ds := &assemblerDataset{reads: reads}
	contigs, _, _ := assembleDataset(ds)
	if len(contigs) != 1 || !sameSequenceEitherStrand(contigs[0], genome) {
		t.Fatalf("contigs = %v", contigs)
	}
}

func TestAssembleContigsSample(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	got := assembleContigs(parseAssemblerReads(string(input)))
	if got != strings.TrimSpace(string(want)) {
		t.Fatalf("got\n%s\nwant\n%s", got, strings.TrimSpace(string(want)))
	}
}

func TestAssembleAtKCircular(t *testing.T) {
	genome := randomDNA(42, 200)
	reads := circularReads(genome, 30, 20)
	ds := &assemblerDataset{reads: reads}
	counts := countSolidKMers(ds, 19, false)
	contigs := assembleAtK(counts, assemblyConfig{k: 19, minCount: 1})
	if len(contigs) != 1 || len(contigs[0]) != len(genome) {
		t.Fatalf("contigs = %v", contigs)
	}
	if !sameSequenceEitherStrand(contigs[0], genome) {
		t.Fatalf("contig is not a rotation of either genome strand")
	}
}

func TestAssembleAtKRemovesErrorTips(t *testing.T) {
	genome := randomDNA(42, 200)
	reads := circularReadsWithErrors(genome, 600, 20, 1)
	ds := &assemblerDataset{reads: reads}
	counts := countSolidKMers(ds, 19, false)
	contigs := assembleAtK(counts, assemblyConfig{k: 19, minCount: 2, filter: true})
	total := 0
	for _, c := range contigs {
		total += len(c)
	}
	if total < len(genome)*9/10 {
		t.Fatalf("total = %d, want at least %d; contigs = %v", total, len(genome)*9/10, contigs)
	}
}

func frDatasetFor(genome string, count, readLen, insert int) *assemblerDataset {
	rng := rand.New(rand.NewSource(11))
	ds := &assemblerDataset{paired: true}
	maxStart := len(genome) - readLen
	for len(ds.pairs) < count {
		firstStart := rng.Intn(maxStart + 1)
		secondStart := firstStart + insert
		if secondStart+readLen > len(genome) {
			continue
		}
		first := []byte(genome[firstStart : firstStart+readLen])
		second := reverseComplementBytes([]byte(genome[secondStart : secondStart+readLen]))
		ds.pairs = append(ds.pairs, readPair{first: first, second: second, distance: insert})
		ds.reads = append(ds.reads, first, second)
	}
	return ds
}

func pairedDatasetFor(genome string, coverage, readLen, insert, spread int) *assemblerDataset {
	length := len(genome)
	count := coverage * length / (2 * readLen)
	if count < 1 {
		count = 1
	}
	rng := rand.New(rand.NewSource(7))
	dataset := &assemblerDataset{paired: true}
	maxStart := length - readLen
	for len(dataset.pairs) < count {
		start := rng.Intn(maxStart + 1)
		distance := insert + rng.Intn(2*spread+1) - spread
		secondStart := start + distance
		if secondStart < 0 || secondStart > maxStart {
			continue
		}
		first := []byte(genome[start : start+readLen])
		second := []byte(genome[secondStart : secondStart+readLen])
		dataset.pairs = append(dataset.pairs, readPair{first: first, second: second, distance: insert})
		dataset.reads = append(dataset.reads, first, second)
	}
	return dataset
}

func nga50(contigs []string, reference string) int {
	length := len(reference)
	if length == 0 {
		return 0
	}
	type block struct {
		sequence string
		length   int
	}
	blocks := make([]block, 0)
	for _, contig := range contigs {
		for _, piece := range strings.Split(contig, "N") {
			if len(piece) > 0 && len(piece) <= length {
				blocks = append(blocks, block{sequence: piece, length: len(piece)})
			}
		}
	}
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].length > blocks[j].length })
	doubled := reference + reference
	covered := make([]bool, length)
	cumulative := 0
	for _, item := range blocks {
		position, added, ok := bestPlacement(item.sequence, doubled, length, covered)
		reverse := string(reverseComplementBytes([]byte(item.sequence)))
		reversePosition, reverseAdded, reverseOK := bestPlacement(reverse, doubled, length, covered)
		if reverseOK && (!ok || reverseAdded > added) {
			position, added, ok = reversePosition, reverseAdded, true
		}
		if !ok {
			continue
		}
		for i := 0; i < item.length; i++ {
			covered[(position+i)%length] = true
		}
		cumulative += added
		if cumulative >= length/2 {
			return item.length
		}
	}
	return 0
}

func bestPlacement(block, doubled string, length int, covered []bool) (int, int, bool) {
	bestPosition := -1
	bestAdded := -1
	searchStart := 0
	for {
		relative := strings.Index(doubled[searchStart:], block)
		if relative < 0 {
			break
		}
		position := searchStart + relative
		added := 0
		for i := 0; i < len(block); i++ {
			if !covered[(position+i)%length] {
				added++
			}
		}
		if added > bestAdded {
			bestAdded = added
			bestPosition = position
		}
		searchStart = position + 1
	}
	if bestPosition < 0 {
		return 0, 0, false
	}
	return bestPosition, bestAdded, true
}

func referenceCoverage(contigs []string, reference string) float64 {
	if len(reference) == 0 {
		return 0
	}
	doubled := reference + reference
	covered := make([]bool, len(reference))
	for _, contig := range contigs {
		for _, block := range strings.Split(contig, "N") {
			if len(block) == 0 || len(block) > len(reference) {
				continue
			}
			position := strings.Index(doubled, block)
			if position < 0 {
				block = string(reverseComplementBytes([]byte(block)))
				position = strings.Index(doubled, block)
			}
			if position < 0 {
				continue
			}
			for i := 0; i < len(block); i++ {
				covered[(position+i)%len(reference)] = true
			}
		}
	}
	count := 0
	for _, value := range covered {
		if value {
			count++
		}
	}
	return float64(count) / float64(len(reference))
}

func linearReadsBothStrands(genome string, coverage, readLen int, seed int64) [][]byte {
	rng := rand.New(rand.NewSource(seed))
	maxStart := len(genome) - readLen
	count := coverage * len(genome) / (2 * readLen)
	reads := make([][]byte, 0, 2*count)
	for i := 0; i < count; i++ {
		start := rng.Intn(maxStart + 1)
		first := []byte(genome[start : start+readLen])
		reads = append(reads, first, reverseComplementBytes(append([]byte(nil), first...)))
	}
	return reads
}

func randomDNA(seed int64, length int) string {
	rng := rand.New(rand.NewSource(seed))
	bases := make([]byte, length)
	for i := range bases {
		bases[i] = "ACGT"[rng.Intn(4)]
	}
	return string(bases)
}

func circularReads(genome string, coverage, readLen int) [][]byte {
	if readLen > len(genome) {
		readLen = len(genome)
	}
	count := coverage * len(genome) / readLen
	if count < len(genome) {
		count = len(genome)
	}
	doubled := genome + genome
	reads := make([][]byte, count)
	for i := range reads {
		start := i % len(genome)
		reads[i] = []byte(doubled[start : start+readLen])
	}
	return reads
}

func circularReadsWithErrors(genome string, count, readLen, errors int) [][]byte {
	doubled := genome + genome
	rng := rand.New(rand.NewSource(1))
	reads := make([][]byte, count)
	for i := range reads {
		start := rng.Intn(len(genome))
		read := []byte(doubled[start : start+readLen])
		for e := 0; e < errors; e++ {
			position := rng.Intn(len(read))
			replacement := byte("ACGT"[rng.Intn(4)])
			for replacement == read[position] {
				replacement = byte("ACGT"[rng.Intn(4)])
			}
			read[position] = replacement
		}
		reads[i] = read
	}
	return reads
}

func isRotation(a, b string) bool {
	return len(a) == len(b) && strings.Contains(b+b, a)
}

func sameSequenceEitherStrand(a, b string) bool {
	return isRotation(a, b) || isRotation(a, string(reverseComplementBytes([]byte(b))))
}

func mustKmer(t *testing.T, s string, k int) kmer {
	t.Helper()
	code, ok := encodeKmer([]byte(s), k)
	if !ok {
		t.Fatalf("invalid k-mer %q", s)
	}
	return code
}

func TestImplementingAnAssemblerFixture(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("test", "1"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("test", "1.a"))
	if err != nil {
		t.Fatal(err)
	}
	reads := parseAssemblerReads(string(input))
	got := assembleContigs(reads)
	want := strings.TrimSpace(string(expected))
	if got != want {
		t.Fatalf("assembleContigs =\n%s\nwant\n%s", got, want)
	}
}

func reverseComplementBytes(read []byte) []byte {
	return reverseComplementSequence(read)
}

func TestAssembleDatasetSmallErrorProneDoubleStranded(t *testing.T) {
	genome := randomDNA(101, 6000)
	doubled := genome + genome
	rng := rand.New(rand.NewSource(19))
	reads := make([][]byte, 2400)
	for i := range reads {
		pos := rng.Intn(len(genome))
		seq := []byte(doubled[pos : pos+100])
		if i%2 == 1 {
			seq = reverseComplementBytes(seq)
		}
		ep := rng.Intn(len(seq))
		seq[ep] = "ACGT"[(strings.IndexByte("ACGT", seq[ep])+1+rng.Intn(3))%4]
		reads[i] = seq
	}
	contigs, _, _ := assembleDataset(&assemblerDataset{reads: reads})
	if len(contigs) != 1 || !sameSequenceEitherStrand(contigs[0], genome) {
		t.Fatalf("contigs = %d (lengths=%v), want 1 circular contig", len(contigs), contigLengths(contigs))
	}
}

func TestBuildUnitigsMergeJunctionSymmetry(t *testing.T) {
	a := randomDNA(201, 200)
	b := randomDNA(202, 200)
	c := randomDNA(203, 200)
	d := randomDNA(204, 200)
	rep := randomDNA(205, 80)
	g1 := a + rep + b
	g2 := c + rep + d
	ds := &assemblerDataset{reads: [][]byte{[]byte(g1), []byte(g2)}}
	counts := countKMers(ds, 31, false, true, false)
	contigs := assembleAtK(counts, assemblyConfig{k: 31, minCount: 1, bothStrands: true})
	canon := canonicalContigs(contigs)
	if len(canon) != 5 {
		t.Fatalf("canonical unitigs = %d (lengths=%v), want 5", len(canon), contigLengths(canon))
	}
}

func TestScaffoldContigsSkipsAmbiguousRepeatHub(t *testing.T) {
	a := randomDNA(301, 600)
	b := randomDNA(302, 600)
	c := randomDNA(303, 600)
	d := randomDNA(304, 600)
	rep := randomDNA(305, 500)
	g1 := a + rep + b
	g2 := c + rep + d
	ds1 := pairedDatasetFor(g1, 30, 100, 350, 20)
	ds2 := pairedDatasetFor(g2, 30, 100, 350, 20)
	pairs := append(append([]readPair(nil), ds1.pairs...), ds2.pairs...)
	contigs := []string{a, b, c, d, rep}
	scaffolds := scaffoldContigs(contigs, pairs, 21, nil, 51)
	for _, s := range scaffolds {
		hasA := strings.Contains(s, a) || strings.Contains(s, string(reverseComplementBytes([]byte(a))))
		hasB := strings.Contains(s, b) || strings.Contains(s, string(reverseComplementBytes([]byte(b))))
		hasC := strings.Contains(s, c) || strings.Contains(s, string(reverseComplementBytes([]byte(c))))
		hasD := strings.Contains(s, d) || strings.Contains(s, string(reverseComplementBytes([]byte(d))))
		if (hasA && hasD) || (hasC && hasB) {
			t.Fatalf("chimeric scaffold joined unrelated contigs across repeat hub")
		}
	}
}

func TestPolishAlignRespectsWindowBounds(t *testing.T) {
	prefix := randomDNA(401, 40)
	target := randomDNA(402, 80)
	suffix := randomDNA(403, 2000)
	contig := prefix + target + suffix
	votes := make([]uint32, len(contig)*4)
	jStart, jEnd, matches, _ := polishAlign([]byte(target), contig, 20, 140, votes, nil)
	if jStart != 40 || jEnd != 120 || matches != 80 {
		t.Fatalf("polishAlign = (%d, %d, %d), want (40, 120, 80)", jStart, jEnd, matches)
	}
}

func TestBuildScaffoldLinksMixedStrandChain(t *testing.T) {
	a := randomDNA(501, 400)
	b := randomDNA(502, 400)
	c := randomDNA(503, 400)
	genome := a + b + c
	pairs := make([]readPair, 0, 20)
	for i := 0; i < 6; i++ {
		off := 120 + i*10
		pairs = append(pairs, readPair{
			first:    []byte(genome[off : off+100]),
			second:   reverseComplementBytes([]byte(genome[off+320 : off+420])),
			distance: 320,
		})
	}
	for i := 0; i < 6; i++ {
		off := 520 + i*10
		pairs = append(pairs, readPair{
			first:    reverseComplementBytes([]byte(genome[off+320 : off+420])),
			second:   []byte(genome[off : off+100]),
			distance: 320,
		})
	}
	reads := make([][]byte, 0, 1600)
	for len(reads) < 1600 {
		reads = append(reads, []byte(genome))
	}
	counts := scaffoldCounts(t, reads, 61)
	scaffolds := scaffoldContigs([]string{a, b, c}, pairs, 21, counts, 61)
	if len(scaffolds) != 1 || len(scaffolds[0]) != len(genome) {
		t.Fatalf("scaffolds = %d (lengths=%v), want 1 contig of len %d", len(scaffolds), contigLengths(scaffolds), len(genome))
	}
}

func TestJoinContigPairOverlapWithJitter(t *testing.T) {
	a := randomDNA(601, 200)
	ov := randomDNA(602, 28)
	b := randomDNA(603, 200)
	from := a + ov
	to := ov + b
	merged, ok := joinContigPair(from, to, -18, 21)
	if !ok || merged != a+ov+b {
		t.Fatalf("joinContigPair with jitter failed: ok=%v len=%d", ok, len(merged))
	}
	if _, badOK := joinContigPair(from, to, 35, 21); badOK {
		t.Fatal("joinContigPair accepted overlap inconsistent with predicted gap")
	}
}

func TestPolishAlignSpanPreservesOverhang(t *testing.T) {
	contig := randomDNA(701, 100)
	overhang := randomDNA(702, 40)
	read := contig[40:] + overhang
	votes := make([]uint32, len(contig)*4)
	rStart, rEnd, qStart, qEnd, matches, _ := polishAlignSpan([]byte(read), contig, 20, len(contig), votes, nil)
	if rStart != 40 || rEnd != 100 || qStart != 0 || qEnd != 60 || matches != 60 {
		t.Fatalf("polishAlignSpan = (%d, %d, %d, %d, %d), want (40, 100, 0, 60, 60)", rStart, rEnd, qStart, qEnd, matches)
	}
}

func TestExtractCountsSkipsErrorTail(t *testing.T) {
	table := newKmerTable(4096)
	for i := 0; i < 600; i++ {
		idx := table.add(kmer{lo: uint64(i + 1)})
		table.counts[idx] = 2
	}
	for i := 0; i < 350; i++ {
		idx := table.add(kmer{lo: uint64(i + 1000)})
		table.counts[idx] = 3
	}
	for i := 0; i < 150; i++ {
		idx := table.add(kmer{lo: uint64(i + 2000)})
		table.counts[idx] = 4
	}
	for i := 0; i < 200; i++ {
		idx := table.add(kmer{lo: uint64(i + 3000)})
		table.counts[idx] = 14
	}
	counts := extractCounts(table, 2, 10000, 21)
	if counts.estimatedCoverage < 10 {
		t.Fatalf("estimatedCoverage = %v, want >= 10 past error tail", counts.estimatedCoverage)
	}
}

func TestMapReadToleratesSingleIndel(t *testing.T) {
	contig := randomDNA(801, 400)
	index := buildContigIndex([]string{contig}, 21)
	sub := contig[100:220]
	// Introduce a 1-bp deletion in the middle of the read so the two halves
	// shift by 1 diagonal on the same contig.
	read := []byte(sub[:60] + sub[61:])
	hit, ok := mapRead(read, index, 21)
	if !ok || hit.contig != 0 || hit.strand != 1 {
		t.Fatalf("mapRead with single deletion = (%+v, %v), want contig 0 strand +1", hit, ok)
	}
}

func TestBuildUnitigsStopsAtSelfRCPalindrome(t *testing.T) {
	// k=5 -> vertex (k-1)=4; "ACGT" is its own reverse complement.
	// Pool flank + "ACGT" and its reverse complement symmetrically.
	flank := randomDNA(901, 20)
	// Ensure no accidental ACGT palindrome in flank by appending non-palindromic bases.
	seq := flank + "AAAT" + "ACGT"
	dataset := &assemblerDataset{reads: [][]byte{[]byte(seq)}}
	counts := countDirectLargeMode(dataset, 5, false, true, uint64(len(seq)-4), true)
	removed := make([]bool, len(counts.keys))
	units := buildUnitigs(counts, 5, removed, true)
	for _, u := range units {
		s := unitigString(counts, u, 5)
		if len(s) >= 20 && s[:6] == string(reverseComplementSequence([]byte(s[len(s)-6:]))) {
			t.Fatalf("buildUnitigs folded through self-RC vertex into hairpin %q", s)
		}
	}
}

func TestBuildScaffoldLinksIgnoresContainedAndTinyStubs(t *testing.T) {
	// Contigs: 0 (large 400bp), 1 (large 400bp), 2 (tiny 120bp stub), 3 (large 400bp deeply overlapping).
	contigs := []string{
		randomDNA(910, 400),
		randomDNA(911, 400),
		randomDNA(912, 120),
		randomDNA(913, 400),
	}
	var mapped []mappedReadPair
	// True adjacency 0 -> 1 with gap -20 (support 10).
	for i := 0; i < 10; i++ {
		mapped = append(mapped, mappedReadPair{
			first:  contigHit{contig: 0, start: 300, strand: 1},
			second: contigHit{contig: 1, start: 180, strand: 1},
			ok:     true,
		})
	}
	// True incoming link 3 -> 0 (gap -20, support 5) plus an incoming competing stub (4)
	// would trigger ambIn && ambOut (blacklisting contig 0) unless tiny stubs are ignored.
	for i := 0; i < 5; i++ {
		mapped = append(mapped, mappedReadPair{
			first:  contigHit{contig: 3, start: 300, strand: 1},
			second: contigHit{contig: 0, start: 180, strand: 1},
			ok:     true,
		})
	}
	// Tiny stub (<250bp) linking to/from contig 0 with gap 10 (support 4) must not mark contig 0 as repeat.
	for i := 0; i < 4; i++ {
		mapped = append(mapped, mappedReadPair{
			first:  contigHit{contig: 0, start: 300, strand: 1},
			second: contigHit{contig: 2, start: 100, strand: 1},
			ok:     true,
		})
		mapped = append(mapped, mappedReadPair{
			first:  contigHit{contig: 2, start: 10, strand: 1},
			second: contigHit{contig: 0, start: 150, strand: 1},
			ok:     true,
		})
	}
	// Deep contained overlap 0 -> 3 with gap -300 (< -250) must be rejected by validLink.
	for i := 0; i < 12; i++ {
		mapped = append(mapped, mappedReadPair{
			first:  contigHit{contig: 0, start: 200, strand: 1},
			second: contigHit{contig: 3, start: 360, strand: 1},
			ok:     true,
		})
	}
	links := buildScaffoldLinks(mapped, contigs, 260, false)
	foundTrue := false
	for _, l := range links {
		if l.from == 0 && !l.fromRev && l.to == 1 && !l.toRev {
			foundTrue = true
		}
		if (l.from == 0 && l.to == 3) || (l.from == 3 && l.to == 0) {
			t.Fatalf("unexpected deep-overlap link kept: %+v", l)
		}
	}
	if !foundTrue {
		t.Fatalf("expected 0->1 scaffold link to survive tiny stub & contained overlap; got %+v", links)
	}
}

func TestPolishReadsAgainstIndelConsensus(t *testing.T) {
	// True sequence has 8 'A's + flank + 6 'T's + flank.
	left := randomDNA(920, 60)
	mid := randomDNA(921, 50)
	right := randomDNA(922, 60)
	truth := left + "AAAAAAAA" + mid + "TTTTTT" + right
	// Draft has a 1-bp insertion in the 'A' run (9 'A's) and 1-bp deletion in 'T' run (5 'T's).
	draft := left + "AAAAAAAAA" + mid + "TTTTT" + right
	reads := make([][]byte, 8)
	for i := range reads {
		reads[i] = []byte(truth)
	}
	polished := polishReadsAgainst(reads, nil, false, []string{draft})
	if len(polished.reads) != len(reads) {
		t.Fatalf("polished read count = %d, want %d", len(polished.reads), len(reads))
	}
	for i, r := range polished.reads {
		if string(r) != truth {
			t.Fatalf("polished read %d mismatch:\n got: %s\nwant: %s", i, string(r), truth)
		}
	}
}
