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
	root := filepath.Join("..", "..", "..")
	data, err := os.ReadFile(filepath.Join(root, "dataset.txt"))
	if err != nil {
		t.Skip("visible dataset not present")
	}
	ref, err := os.ReadFile(filepath.Join(root, "genome.txt"))
	if err != nil {
		t.Skip("reference genome not present")
	}
	reference := strings.TrimSpace(string(ref))
	ds := parseAssemblerDataset(data)
	contigs := assembleDataset(ds)
	if got := referenceCoverage(contigs, reference); got < 0.95 {
		t.Fatalf("coverage = %v", got)
	}
	if got := nga50(contigs, reference); got < len(reference)*8/10 {
		t.Fatalf("NGA50 = %d, want at least %d", got, len(reference)*8/10)
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
	scaffolds := scaffoldContigs(contigs, ds.pairs, 21)
	if len(scaffolds) >= len(contigs) {
		t.Fatalf("scaffolding did not reduce contig count: %d -> %d", len(contigs), len(scaffolds))
	}
	if referenceCoverage(scaffolds, genome) < 0.6 {
		t.Fatalf("coverage = %v", referenceCoverage(scaffolds, genome))
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
	contigs := assembleDataset(&assemblerDataset{reads: double})
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

func TestAssembleDatasetSingle(t *testing.T) {
	genome := randomDNA(42, 200)
	reads := circularReads(genome, 40, 20)
	ds := &assemblerDataset{reads: reads}
	contigs := assembleDataset(ds)
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
