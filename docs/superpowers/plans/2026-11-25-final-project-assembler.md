# Final Project Assembler Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a single-file Go genome assembler that passes the Coursera final project hidden datasets by maximizing NGA50.

**Architecture:** Parse reads or read-pairs, count solid k-mers with a Bloom-filtered two-pass counter, build unitigs on the implicit de Bruijn graph, remove low-coverage tips and bubbles, choose the best candidate `k`, then scaffold contigs with read-pair links.

**Tech Stack:** Go 1.21 standard library only; `go test`; no external assembler or alignment tool.

**Spec:** `docs/superpowers/specs/2026-11-25-final-project-assembler-design.md`

## Global Constraints

- Coursera compiles exactly one Go file with no `go.mod`.
- No `dsa/...` imports in the solution file.
- Read stdin; write FASTA to stdout.
- Final project limits: about 100 s and 512 MB.
- Visible data: `dataset.txt` has 33,609 100 bp reads; `genome.txt` is 112,031 bp.
- Double-stranded data uses canonical k-mers; single-stranded data uses forward k-mers.
- Tiny datasets use `minCount = 1` and skip coverage filtering when `estimatedCoverage < 4` or `len(reads) < 100`.
- `kMax = min(61, max(3, medianReadLength-1))`.
- Commit dates continue after `2025-04-14`; use `2025-04-15T10:00:00+0700` for this feature.

## Review Focus

- One very short real read must not collapse `k`; median read length controls the cap.
- A forward-only dataset must not be treated as double-stranded and duplicated.
- An error-prone 1.4M-read dataset must not store tens of millions of singleton k-mers.
- The PDF sample with 3 bp reads must still output the exact seven contigs.
- A read pair spanning a repeat must join the correct contig ends, not their reverse complements.

---

### Task 1: Parse Reads and Read-Pairs

**Files:**
- Modify: `course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: none.
- Produces:
  - `type readPair struct { first, second []byte; distance int }`
  - `type assemblerDataset struct { reads [][]byte; pairs []readPair; paired bool }`
  - `func parseAssemblerDataset(input []byte) *assemblerDataset`
  - `func parseAssemblerReads(input string) []string`
  - `func assembleContigs(reads []string) string`
  - `func formatFASTA(contigs []string) string`

- [ ] **Step 1: Write failing parser tests**

Add `TestParseAssemblerDatasetSingle` and `TestParseAssemblerDatasetPairs` to the test file.

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestParseAssemblerDataset' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `parseAssemblerDataset` is undefined.

- [ ] **Step 3: Implement parsing**

Read all input, split on newlines, trim `\r` and spaces, skip blank lines, drop the first non-blank line, and split pair lines on `|`. Store slices that reference the input buffer so reads are not copied. Uppercase bases in place. Keep `parseAssemblerReads` as a wrapper that returns all reads from the parsed dataset. Keep `formatFASTA` deterministic: sort contigs lexicographically before printing.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestParseAssemblerDataset|TestImplementingAnAssemblerFixture' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS. The existing sample test may fail until Task 4, so run only parser tests now.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler.go course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler_test.go docs/superpowers
GIT_AUTHOR_DATE="2025-04-15T10:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T10:00:00+0700" git commit -m "feat: parse assembler reads and read-pairs" --date="2025-04-15T10:00:00+0700"
```

---

### Task 2: Two-Bit K-mer Encoding

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: none.
- Produces:
  - `type kmer struct { lo, hi uint64 }`
  - `func encodeKmer(seq []byte, k int) (kmer, bool)`
  - `func decodeKmer(code kmer, k int) []byte`
  - `func (code kmer) prefix() kmer`
  - `func (code kmer) suffix(k int) kmer`
  - `func (code kmer) appendBase(base byte, k int) kmer`
  - `func (code kmer) reverseComplement(k int) kmer`
  - `func canonicalKmer(code kmer, k int) kmer`
  - `func kmerLess(a, b kmer) bool`

- [ ] **Step 1: Write failing encoding tests**

```go
func TestKmerRoundTrip(t *testing.T) {
	for _, k := range []int{3, 21, 31, 32, 37, 55, 61} {
		seq := []byte("ACGTACGTTGCAACGTACGTACGTACGTACGTACGTACGTACGTACGTACGTACGTACGTA"[:k])
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestKmer' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `kmer` is undefined.

- [ ] **Step 3: Implement encoding**

Use `.lo` for the low 64 bits and `.hi` for bits 64..127. The last base is in the low two bits. `prefix` shifts right by two. `suffix` shifts left by two and masks to `2*(k-1)` bits. `appendBase` shifts left and masks to `2*k` bits. Return `false` for any base outside `ACGT`. `canonicalKmer` returns `min(code, code.reverseComplement(k))` using `kmerLess`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestKmer' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T11:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T11:00:00+0700" git commit -m "feat: add two-bit k-mer encoding" --date="2025-04-15T11:00:00+0700"
```

---

### Task 3: Bloom-Filtered K-mer Counting

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: `kmer`, `encodeKmer`, `canonicalKmer`, `assemblerDataset`.
- Produces:
  - `type countedKMers struct { keys []kmer; counts []uint32; totalInstances uint64; estimatedCoverage float64 }`
  - `func countSolidKMers(ds *assemblerDataset, k int, canonical bool) *countedKMers`
  - `func (c *countedKMers) index(code kmer) int`
  - `func (c *countedKMers) count(code kmer) uint32`
  - `func estimateErrorFraction(ds *assemblerDataset, counts *countedKMers, k int) float64`

- [ ] **Step 1: Write failing counting tests**

```go
func TestCountSolidKMers(t *testing.T) {
	ds := &assemblerDataset{reads: [][]byte{[]byte("ACGTACGT"), []byte("ACGTACGT"), []byte("TTTTACGT")}}
	counts := countSolidKMers(ds, 4, false)
	if counts.count(mustKmer("ACGT", 4)) != 5 {
		t.Fatalf("ACGT count = %d", counts.count(mustKmer("ACGT", 4)))
	}
	if counts.count(mustKmer("TTTT", 4)) != 1 {
		t.Fatalf("TTTT count = %d", counts.count(mustKmer("TTTT", 4)))
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
```

Add helper `mustKmer(t *testing.T, s string, k int) kmer`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestCount|TestEstimateError' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `countSolidKMers` is undefined.

- [ ] **Step 3: Implement the counter**

For datasets with at most 100 reads or fewer than 100,000 k-mer instances, count directly into an open-addressed table. For larger datasets, use the two-pass Bloom path:

- Bloom bitset: `1 << 30` bits (128 MB) for large data, two independent 32-bit hashes from `lo` and `hi`.
- Pass 1 over all reads: test the two bits; if both are set, insert the k-mer into a candidate table; then set both bits.
- Pass 2: count only k-mers present in the candidate table.
- Extract keys with count at least 2, sort keys and counts together with a `sort.Interface` that swaps both slices.
- `estimatedCoverage = totalInstances / len(keys)`.

For small datasets, use `minCount = 1` and return all counted keys.

`estimateErrorFraction` samples up to 1,000 reads and returns the fraction of their k-mers absent from `counts`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestCount|TestEstimateError' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T12:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T12:00:00+0700" git commit -m "feat: add Bloom-filtered k-mer counting" --date="2025-04-15T12:00:00+0700"
```

---

### Task 4: Unitigs, Coverage Filtering, Tips, and Sample Compatibility

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: `countedKMers`, `kmer`, `assemblerDataset`.
- Produces:
  - `type assemblyConfig struct { k int; minCount uint32; filter bool }`
  - `func assembleAtK(counts *countedKMers, cfg assemblyConfig) []string`
  - `func (c *countedKMers) vertexIn(vertex kmer, k int) int`
  - `func (c *countedKMers) vertexOut(vertex kmer, k int) int`

- [ ] **Step 1: Write failing unitig tests**

```go
func TestAssembleAtKCircular(t *testing.T) {
	genome := "ACGTTGCAACGTACGATCGTAGCTAGCATCGATCGTACGATCGATGCTAGCTAGCATCGATCGTAGCTAGCTACG"
	reads := circularReads(genome, 30, 12)
	ds := &assemblerDataset{reads: reads}
	counts := countSolidKMers(ds, 11, false)
	contigs := assembleAtK(counts, assemblyConfig{k: 11, minCount: 1})
	if len(contigs) != 1 || len(contigs[0]) != len(genome) {
		t.Fatalf("contigs = %v", contigs)
	}
	if !isRotation(contigs[0], genome) {
		t.Fatalf("contig is not a rotation of the genome")
	}
}

func TestAssembleAtKRemovesErrorTips(t *testing.T) {
	genome := "ACGTTGCAACGTACGATCGTAGCTAGCATCGATCGTACGATCGATGCTAGCTAGCATCGATCGTAGCTAGCTACG"
	reads := circularReadsWithErrors(genome, 120, 12, 1)
	ds := &assemblerDataset{reads: reads}
	counts := countSolidKMers(ds, 11, false)
	contigs := assembleAtK(counts, assemblyConfig{k: 11, minCount: 2, filter: true})
	total := 0
	for _, c := range contigs {
		total += len(c)
	}
	if total < len(genome)*9/10 {
		t.Fatalf("total = %d, want at least %d; contigs = %v", total, len(genome)*9/10, contigs)
	}
}
```

Add test helpers `circularReads`, `circularReadsWithErrors`, `isRotation`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestAssembleAtK' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `assembleAtK` is undefined.

- [ ] **Step 3: Implement unitig construction**

Use the implicit graph:

- For a k-mer edge, source is `prefix()` and destination is `suffix(k)`.
- Start a unitig when `vertexIn(source) != 1 || vertexOut(source) != 1 || vertexIn(destination) != 1`.
- Walk while the destination vertex has indegree 1 and outdegree 1.
- Mark each k-mer edge used once.
- Handle pure cycles by starting from any unused edge.
- For cycles, trim the repeated `k-1` suffix so the contig length equals the cycle edge count.

When `cfg.filter` is true:

- Estimate `threshold = max(2, 0.20*counts.estimatedCoverage)`.
- Remove unitigs with average coverage below `threshold` and length at most `max(2*k, 500)`.
- Remove dead-end tips under the same length and coverage rule.
- Repeat coverage/tip removal until no k-mer is removed.

When `cfg.filter` is false, emit every unitig. This is the tiny-input path that preserves the PDF sample.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestAssembleAtK|TestImplementingAnAssemblerFixture' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS for the unitig tests. The fixture requires `assembleContigs` wiring, which is completed in Task 5.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T13:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T13:00:00+0700" git commit -m "feat: assemble unitigs and remove low-coverage tips" --date="2025-04-15T13:00:00+0700"
```

---

### Task 5: Candidate K Selection and Single-Read Pipeline

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: `assemblerDataset`, `countSolidKMers`, `estimateErrorFraction`, `assembleAtK`, `formatFASTA`.
- Produces:
  - `func assembleDataset(ds *assemblerDataset) []string`
  - `func assembleContigs(reads []string) string`
  - `func medianReadLength(ds *assemblerDataset) int`
  - `func assemblyScore(contigs []string, genomeEstimate float64) float64`

- [ ] **Step 1: Write failing pipeline tests**

```go
func TestAssembleDatasetSingle(t *testing.T) {
	genome := "ACGTTGCAACGTACGATCGTAGCTAGCATCGATCGTACGATCGATGCTAGCTAGCATCGATCGTAGCTAGCTACG"
	reads := circularReads(genome, 40, 12)
	ds := &assemblerDataset{reads: reads}
	contigs := assembleDataset(ds)
	if len(contigs) != 1 || !isRotation(contigs[0], genome) {
		t.Fatalf("contigs = %v", contigs)
	}
}

func TestAssembleContigsSample(t *testing.T) {
	input, _ := os.ReadFile(filepath.Join("test", "1"))
	want, _ := os.ReadFile(filepath.Join("test", "1.a"))
	got := assembleContigs(parseAssemblerReads(string(input)))
	if got != strings.TrimSpace(string(want)) {
		t.Fatalf("got\n%s\nwant\n%s", got, strings.TrimSpace(string(want)))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestAssembleDatasetSingle|TestAssembleContigsSample' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `assembleDataset` is undefined or the sample output is wrong.

- [ ] **Step 3: Implement selection**

- `medianReadLength` uses reads and pairs.
- `kMax = min(61, max(3, median-1))`.
- Tiny datasets (`len(reads) < 100` or estimated coverage < 4) use `k = min(kMax, median)`, `minCount = 1`, `filter = false`.
- Count at `k=31` or the largest odd `k <= min(31, kMax)` for error detection.
- If `estimateErrorFraction > 0.10`, use candidates `[31, 37]` capped by `kMax`.
- Otherwise use candidates `[45, 55, 61]` capped by `kMax`; count directly, not with Bloom, when possible.
- Always add `kMax` itself as a fallback candidate, so small genomes still have at least one valid `k`.
- For each candidate, assemble, compute `assemblyScore`, and keep the best. `assemblyScore` returns `n50 * min(1, totalLength/genomeEstimate)`; tie-break toward larger `k`.
- `assembleContigs` parses a single-read dataset and calls `assembleDataset`, then `formatFASTA`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestAssembleDatasetSingle|TestAssembleContigsSample|TestImplementingAnAssemblerFixture' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS, including the exact seven-contig PDF sample.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T14:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T14:00:00+0700" git commit -m "feat: select k and assemble single-read datasets" --date="2025-04-15T14:00:00+0700"
```

---

### Task 6: Double-Strand Detection and Canonical Output

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: `countedKMers`, `kmer`, `reverseComplement`, `canonicalKmer`.
- Produces:
  - `func isDoubleStranded(ds *assemblerDataset, counts *countedKMers, k int) bool`
  - `func canonicalContigs(contigs []string) []string`

- [ ] **Step 1: Write failing strand tests**

```go
func TestIsDoubleStranded(t *testing.T) {
	genome := "ACGTTGCAACGTACGATCGTAGCTAGCATCGATCGTACGATCGATGCTAGCTAGCATCGATCGTAGCTAGCTACG"
	forward := circularReads(genome, 40, 12)
	double := append([] []byte{}, forward...)
	for _, read := range forward {
		double = append(double, reverseComplementBytes(read))
	}
	fwdDS := &assemblerDataset{reads: forward}
	fwdCounts := countSolidKMers(fwdDS, 11, false)
	if isDoubleStranded(fwdDS, fwdCounts, 11) {
		t.Fatal("forward-only data detected as double-stranded")
	}
	doubleDS := &assemblerDataset{reads: double}
	doubleCounts := countSolidKMers(doubleDS, 11, false)
	if !isDoubleStranded(doubleDS, doubleCounts, 11) {
		t.Fatal("double-stranded data not detected")
	}
}

func TestCanonicalContigs(t *testing.T) {
	got := canonicalContigs([]string{"ACGT", "ACGT", "TTTT"})
	if len(got) != 2 {
		t.Fatalf("contigs = %v", got)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run 'TestIsDoubleStranded|TestCanonicalContigs' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because the functions are undefined.

- [ ] **Step 3: Implement detection and canonicalization**

`isDoubleStranded` samples up to 1,000 reads. For each read, count k-mers of its reverse complement that are present in `counts`. Return true when the present fraction is above 0.20.

`canonicalContigs` keeps `min(contig, reverseComplement(contig))`, removes duplicates, and sorts the result.

Wire detection into `assembleDataset`: the first k=31 count is non-canonical. If double-stranded, re-count with `canonical = true` before choosing candidates, and canonicalize contigs before returning.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run 'TestIsDoubleStranded|TestCanonicalContigs|TestAssembleDatasetSingle' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T15:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T15:00:00+0700" git commit -m "feat: handle double-stranded reads with canonical k-mers" --date="2025-04-15T15:00:00+0700"
```

---

### Task 7: Read-Pair Scaffolding

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`

**Interfaces:**
- Consumes: `readPair`, contigs from `assembleDataset`.
- Produces:
  - `func scaffoldContigs(contigs []string, pairs []readPair, seedK int) []string`
  - `type contigHit struct { contig, start int; strand int8 }`
  - `func mapRead(read []byte, index map[uint64]contigHit, k int) (contigHit, bool)`

- [ ] **Step 1: Write a failing scaffolding test**

Build a circular genome with three unique flanks separated by two copies of a 220 bp repeat:

```go
func TestScaffoldContigsSpansRepeat(t *testing.T) {
	uniqueA := randomDNA(1, 300)
	uniqueB := randomDNA(2, 300)
	uniqueC := randomDNA(3, 300)
	repeat := strings.Repeat("ACGT", 55) // 220 bp
	genome := uniqueA + repeat + uniqueB + repeat + uniqueC
	// Generate 30x 100 bp read-pairs with insert 350.
	ds := pairedDatasetFor(genome, 30, 100, 350, 20)
	contigs := assembleDataset(ds)
	scaffolds := scaffoldContigs(contigs, ds.pairs, 21)
	if len(scaffolds) >= len(contigs) {
		t.Fatalf("scaffolding did not reduce contig count: %d -> %d", len(contigs), len(scaffolds))
	}
	if referenceCoverage(scaffolds, genome) < 0.95 {
		t.Fatalf("coverage = %v", referenceCoverage(scaffolds, genome))
	}
}
```

Add helpers `randomDNA`, `pairedDatasetFor`, and `referenceCoverage`. `referenceCoverage` must split contigs on `N` before exact placement, because scaffolds may contain gap runs.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run 'TestScaffoldContigsSpansRepeat' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: FAIL because `scaffoldContigs` is undefined.

- [ ] **Step 3: Implement scaffolding**

- Build a unique seed index over contig k-mers (21-mers, or 15 when reads are short). Do not store k-mers that occur more than once.
- For each pair, sample seeds across both mates and both orientations. A mate maps when several seeds agree on the same contig and nearby position.
- Convert each mapped pair into a candidate link between two contig ends. Record orientation and estimated gap from `d`.
- Count support by `(fromContig, fromEnd, toContig, toEnd, orientation)`. Keep the best unique link per contig end when support is at least 3 and at least twice the runner-up.
- Greedily merge contig chains. Reverse-complement a contig when the link orientation requires it.
- If adjacent sequences overlap by at least `seedK`, merge without a gap. Otherwise insert `N`s for the estimated gap.
- Return unchanged contigs when no pair maps or no link is unique.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -run 'TestScaffoldContigsSpansRepeat' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T16:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T16:00:00+0700" git commit -m "feat: scaffold contigs with read-pair links" --date="2025-04-15T16:00:00+0700"
```

---

### Task 8: Visible Dataset Test and NGA50 Helper

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler_test.go`
- No production code expected unless the visible dataset exposes a bug.

**Interfaces:**
- Consumes: `parseAssemblerDataset`, `assembleDataset`.
- Produces:
  - `func nga50(contigs []string, reference string) int`
  - `func referenceCoverage(contigs []string, reference string) float64`
  - `func TestVisibleDatasetAssembly(t *testing.T)`

- [ ] **Step 1: Add the visible dataset test**

```go
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
```

- [ ] **Step 2: Implement NGA50 and coverage in the test file**

Use exact circular placement:

- Concatenate `reference + reference`.
- Split each contig on `N` into aligned blocks before placement.
- For each block, find all exact matches in the doubled reference, choose the placement that adds the most uncovered reference bases.
- `referenceCoverage` counts covered reference positions divided by reference length.
- `nga50` sorts placed blocks by length descending, accumulates newly covered bases, and returns the block length when cumulative coverage reaches 50% of the reference.

- [ ] **Step 3: Run the visible dataset test**

Run: `go test -run 'TestVisibleDatasetAssembly' ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v`

Expected: PASS with coverage at least 0.95 and NGA50 at least 89,624.

If it fails, inspect the selected `k`, the estimated error fraction, and the number of singleton k-mers. Fix the pipeline, not the test threshold.

- [ ] **Step 4: Run the whole package and repository tests**

Run: `go test -count=1 ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler`

Expected: PASS.

Run: `go test -count=1 ./...`

Expected: all packages pass.

- [ ] **Step 5: Commit**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T17:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T17:00:00+0700" git commit -m "test: verify assembler on the visible dataset" --date="2025-04-15T17:00:00+0700"
```

---

### Task 9: Large-Data Performance and Memory

**Files:**
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go`
- Test: manual benchmark, no committed 1.4M-read fixture.

**Interfaces:**
- Consumes: the finished assembler.
- Produces: a binary that fits 512 MB and 100 s on 1.4M 100 bp reads.

- [ ] **Step 1: Build the final binary**

Run:

```bash
go build -o /tmp/opencode/assembler ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
```

Expected: binary created.

- [ ] **Step 2: Generate 1.4M synthetic E. coli reads**

Use the downloaded reference at `/tmp/opencode/ecoli.txt` if present. Generate 1.4M reads of length 100 with a fixed seed, one error per read, in FASTA-like input format.

- [ ] **Step 3: Measure time and memory**

Run with `/usr/bin/time -v`:

```bash
/usr/bin/time -v /tmp/opencode/assembler < /tmp/opencode/ecoli_1.4m.txt > /tmp/opencode/ecoli_contigs.fa
```

Expected: elapsed under 100 s and maximum resident set size under 512 MB. If memory is over, reduce the Bloom filter to 64 MB, shrink the candidate table, or drop to `uint64` keys for `k <= 32`.

- [ ] **Step 4: Verify assembly sanity**

Run:

```bash
grep -c '^>' /tmp/opencode/ecoli_contigs.fa
awk '/^>/{next}{n+=length($0)} END{print n}' /tmp/opencode/ecoli_contigs.fa
```

Expected: total contig length between 0.9x and 1.2x the reference length.

- [ ] **Step 5: Commit any performance fixes**

```bash
git add course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T18:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T18:00:00+0700" git commit -m "perf: keep assembler within Coursera memory limits" --date="2025-04-15T18:00:00+0700"
```

---

### Task 10: README, Final Verification, and Submission Check

**Files:**
- Modify: `README.md`
- Modify: `.../6_implementing_an_assembler/implementing_an_assembler.go` only if verification finds a defect.

**Interfaces:**
- Consumes: all previous tasks.
- Produces: documented usage and a submission-ready file.

- [ ] **Step 1: Update README**

Add a short section under Course 6:

- the final project accepts single reads or `READ1|READ2|d` pairs;
- output is FASTA;
- local run command;
- visible dataset test command;
- note that `dataset.txt` and `genome.txt` are optional local fixtures and are not committed.

- [ ] **Step 2: Run formatting and vet**

Run:

```bash
gofmt -w course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
go vet ./...
gofmt -l .
```

Expected: `gofmt -l .` prints nothing; `go vet` exits 0.

- [ ] **Step 3: Run the full suite**

Run: `go test -count=1 ./...`

Expected: all packages pass.

- [ ] **Step 4: Verify the single-file build**

Run:

```bash
GO111MODULE=off go build -o /tmp/opencode/final_assembler course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler.go
```

Expected: build succeeds without `go.mod`.

- [ ] **Step 5: Run the visible dataset once more**

Run:

```bash
/tmp/opencode/final_assembler < dataset.txt > /tmp/opencode/visible_contigs.fa
grep -c '^>' /tmp/opencode/visible_contigs.fa
```

Expected: one or a few contigs.

- [ ] **Step 6: Commit**

```bash
git add README.md course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler
GIT_AUTHOR_DATE="2025-04-15T19:00:00+0700" GIT_COMMITTER_DATE="2025-04-15T19:00:00+0700" git commit -m "docs: document final project assembler" --date="2025-04-15T19:00:00+0700"
```
