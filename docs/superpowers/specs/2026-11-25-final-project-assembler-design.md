# Final Project Assembler Design

Date: 2026-11-25
Problem: `course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler`

## Intent

Replace the current sample-only assembler with a general de Bruijn assembler that
can be submitted to the Coursera final project. The grader compares assembled
contigs against SPAdes with NGA50 on hidden datasets: error-free and error-prone
single reads, error-free and error-prone read-pairs, and real read-pairs.

The visible dataset is `dataset.txt` (33,609 error-free 100 bp reads from
N. deltocephalinicola) with reference `genome.txt` (112,031 bp). It is one
dataset, not the whole hidden benchmark.

## Goals

- Read both input formats:
  - single reads: `t` followed by `t` reads;
  - read-pairs: `t` followed by `t` lines of `READ1|READ2|d`.
- Assemble near-complete contigs from the visible N. deltocephalinicola dataset.
- Handle 1% substitution errors with coverage filtering and graph cleanup.
- Use read-pair links to scaffold contigs across repeats.
- Stay within the Coursera limits: about 100 s and 512 MB for the final project.
- Remain a single self-contained Go file with no `dsa/...` imports.
- Keep tests deterministic and runnable without the untracked datasets.
- Update `README.md` with the final-project behavior and usage.

## Non-goals

- Reimplement SPAdes exactly.
- Handle insertions/deletions or high error rates beyond the course simulation.
- Commit the 3.3 MB `dataset.txt` or `genome.txt` without explicit approval.
- Push or open a pull request.
- Preserve the old sample-only assembly algorithm.

## Constraints

- Coursera compiles one Go file with no `go.mod`.
- Input is stdin; output is FASTA on stdout.
- Hidden datasets include up to about 1,400,000 single reads or 700,000
  read-pairs. Memory must stay below 512 MB.
- Reads may come from one strand or both strands. Real data comes from both.
- Read-pair distance `d` is only a center of a distribution, not exact.
- The visible dataset is error-free and forward-strand only.
- The repository convention is `go test` per problem, one solution file, and
  inlined helpers.

## Decisions made here

These are builder calls, not user decisions.

### Language and location

Rewrite
`6_implementing_an_assembler/implementing_an_assembler.go` in Go and keep the
existing sample test. Add tests in `implementing_an_assembler_test.go`.

### K-mer representation

Use a two-bit encoding in up to 128 bits. For `k <= 32` use one `uint64`; for
larger `k` use a `{lo, hi uint64}` key. Store counts in a custom open-addressing
hash table, not a Go map, so 5 million distinct k-mers fit with room for reads.
The table uses power-of-two capacity, linear probing, and count zero as empty.

After counting, extract solid keys into a sorted slice with a parallel count
array. Delete the hash table before graph traversal. Graph membership and
degree checks use binary search over the sorted keys, which avoids a second
large hash map.

### Error handling and assembly

Use the following pipeline for a candidate `k`:

1. Count all k-mers from reads. For double-stranded data, count canonical
   k-mers (`min(kmer, reverseComplement(kmer))`) instead of adding a second
   copy of every read.
2. Keep k-mers with count at least 2. A true k-mer with count 1 is negligible
   at the course coverage; singleton k-mers are almost always errors.
3. Estimate coverage as `total k-mer instances / number of solid k-mers`.
4. Build unitigs on the implicit de Bruijn graph. A k-mer edge is internal when
   its source vertex has outdegree 1 and its destination vertex has indegree 1.
   Start a unitig at every branch or dead end, walk until the next branch, and
   handle pure cycles separately.
5. Remove unitigs whose average k-mer coverage is below
   `max(2, 0.20 * estimatedCoverage)` and whose length is at most
   `max(2*k, 500)`. Remove their k-mers and rebuild unitigs. Repeat until no
   unitig is removed.
6. Remove short dead-end tips with the same coverage and length rule.
7. Remove short bubbles: for each branch, follow pairs of internally disjoint
   paths that reconverge within `max(2*k, 500)` edges; delete the path with the
   lower average k-mer coverage. Keep the higher-coverage path. Repeat with the
   tip pass until stable.
8. Emit the remaining unitigs as contigs.

The visible dataset at `k=35..61` is expected to produce one or a few contigs.
Simulated 1% error-prone N. deltocephalinicola data at `k=31..37` is expected to
produce near-complete contigs after unitig-coverage filtering.

### Counting error-prone large data

A direct count table cannot hold all singleton error k-mers from 1.4 million
error-prone E. coli reads. Use two passes over the stored reads.

Pass 1 uses a Bloom filter:

- For every k-mer, test the Bloom bits before setting them.
- If the k-mer was already seen, insert its key into a candidate hash set.
- Set the Bloom bits.

Pass 2 counts only candidate k-mers. This keeps the count table at true
coverage plus a small Bloom false-positive rate. Singletons never enter the
candidate set.

For 1.4 million reads, allocate about 64 MB for the Bloom filter and size the
candidate table for 8 million entries. Error-free large datasets may skip the
Bloom pass and count directly at the larger `k`, where the distinct count is
about the genome size.

### Detecting errors

After the k=31 candidate pass, sample up to 1,000 reads. Compute the fraction of
their k-mers that are absent from the candidate set. At 1% substitution errors
and k=31, about 30-45% of a read's k-mers cross an error; error-free data is
near zero. Treat the data as error-prone when more than 10% of sampled k-mers
are absent.

### Strand policy

- Single-stranded data: use forward k-mers only. The visible dataset is
  single-stranded.
- Double-stranded data: canonicalize every k-mer before counting, and
  canonicalize every output contig (keep `min(contig, reverseComplement)`).
  Remove duplicate reverse-complement contigs.
- Detection runs after the first k=31 count: sample reads, reverse-complement
  them, and measure how many of their k-mers are present in the k=31 table.
  More than 20% presence means double-stranded data.

This avoids storing both strands and avoids emitting each contig twice.

### The tiny-input exception

The PDF sample has 3 bp reads and coverage near 1. Applying `count >= 2` would
collapse it. Do not apply solid-k-mer filtering or coverage-based unitig removal
when `estimatedCoverage < 4` or the dataset has fewer than 100 reads. In that
case use `minCount = 1` and emit all unitigs. This preserves the exact 7-contig
sample output.

### Capping k on variable-length reads

Use the median read length, not the shortest read, to cap `k`:

`kMax = min(61, max(3, medianReadLength - 1))`.

Reads shorter than `k` are ignored for k-mer counting but remain available for
read-pair scaffolding when long enough. This prevents one truncated real read
from collapsing `k`.

### Choosing k

- Detect error rate from the k=31 count: if the fraction of distinct k-mers
  with count 1 is high, treat the data as error-prone.
- Error-prone data: candidates `k = 31, 37`.
- Error-free data: candidates `k = 45, 55, 61`, capped by the shortest read.
- Score each candidate assembly without a reference:
  - reject candidates whose total contig length is far below the estimated
    genome size;
  - prefer higher N50 and longer longest contig;
  - prefer fewer contigs at similar N50;
  - tie-break toward the larger `k`.
- Run only the candidates needed. The k=31 count is reused for error-prone
  assembly; error-free data pays for a second count at the larger `k`.

### Double-stranded detection

Sample up to 1,000 reads. If at least 20% of the sampled reads have their
reverse complement present in the read set, treat the data as double-stranded
and add reverse complements during counting. The visible dataset is
forward-only, so it should not trigger this path. Real data should.

### Read-pair scaffolding

After choosing contigs:

1. Index unique 21-mers of all contigs with their contig ID, position, and
   strand. Ambiguous 21-mers (repeats) are not used as seeds.
2. For each pair, sample seeds across both mates and their reverse complements.
   A mate maps when several seeds agree on the same contig, position, and
   strand. Use k=15 seeds when reads are shorter than 40 bp.
3. Convert each mapped pair into a link between contig ends, including
   orientation and estimated gap from `d` and the observed positions.
4. Count link support. Merge a contig end only when its best link is unique and
   has enough support. Build chains, orient contigs, and join them.
5. Join with an exact overlap when one exists; otherwise insert `N`s for the
   estimated gap. Correct orientation matters more than exact gap length.

Scaffolding is conservative: a wrong join can reduce NGA50 more than leaving
two contigs separate.

## Testing

Tests are written before implementation and observed failing.

- Existing sample from the PDF must still produce the documented FASTA contigs.
- Parser tests for single reads and `READ1|READ2|d` pairs.
- Two-bit key tests for prefix, suffix, reverse complement, and k around 32.
- Small circular genome assembled from exact reads: one contig, correct
  rotation, all reads present.
- Small circular genome with 1% substitution errors: one or a few contigs and
  high reference coverage.
- Repeat fixture where 100 bp reads produce two contigs at large `k`; synthetic
  read-pairs span the repeat and scaffolding joins them.
- Visible dataset integration test:
  - skip when `dataset.txt` or `genome.txt` is absent;
  - require at least 95% of the reference covered;
  - require NGA50 at least 80% of the reference length.
- NGA50 in tests is computed by exact circular placement against the reference,
  not by calling QUAST. A contig is placed at the reference position that adds
  the most newly covered bases. Aligned blocks are sorted by length, and NGA50
  is the block length at which cumulative newly covered bases reach 50% of the
  reference. This is valid because error-free and error-corrected contigs from
  these tests are exact substrings of the reference.
- A short-mode test on a deterministic synthetic dataset, so `go test ./...`
  does not depend on the 3.3 MB files.
- Manual performance check with 1,400,000 synthetic E. coli reads. The target
  is under 100 s and 512 MB.

## Success criteria

- `go test -count=1 ./...` passes with 113+ packages and no skipped required
  tests in a clean checkout.
- `go vet ./...` and `gofmt` are clean.
- The final file builds with Go 1.20.14 in module-off mode.
- The visible N. deltocephalinicola dataset assembles to at least 95% reference
  coverage and at least 80% reference-length NGA50.
- Synthetic error-prone and read-pair datasets improve over the current
  sample-only assembler.
- The submitted file remains self-contained.

## What is left to the builder

- Exact custom hash-table layout and probing details.
- Exact unitig coverage estimator and bubble traversal implementation.
- Exact scaffold link scoring and conflict handling.
- Test fixture sequences and deterministic random seeds.
- Whether to add the visible datasets to the repository after the user
  explicitly asks.
