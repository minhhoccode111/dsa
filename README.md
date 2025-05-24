# Data Structures and Algorithms Specialization Coursera

**Disclaimer**: The code is terrible and buggy, it has been a few months since I worked on it. Use at your own risk.

[Data Structures and Algorithms Specialization](https://www.coursera.org/specializations/data-structures-algorithms)

This repository contains both Python and Go implementations of the assignments. The Go solutions live next to the Python ones, include tests, and can be run from the repository root.

## Languages

- Python 3
- Go 1.21+ (`go.mod` uses the `dsa` module)

## Commands

Run `make help` for the live list. Targets are documented in the `Makefile` with
`##` comments.

| Target | Description |
|---|---|
| `make test` | Run Go unit tests |
| `make fixtures` | Run Go IO fixtures (`test/<name>` against `test/<name>.a`) |
| `make fixtures-python` | Run the same fixtures with Python solutions |

`make fixtures` builds each Go solution once, checks every fixture, prints
failures, and exits non-zero on failure. It writes no files.

## Course 1 - [Algorithmic Toolbox](https://www.coursera.org/learn/algorithmic-toolbox)

## Course 2 - [Data Structures](https://www.coursera.org/learn/data-structures)

## Course 3 - [Algorithms on Graphs](https://www.coursera.org/learn/algorithms-on-graphs)

## Course 4 - [Algorithms on Strings](https://www.coursera.org/learn/algorithms-on-strings)

## Course 5 - [Advanced Algorithms and Complexity](https://www.coursera.org/learn/advanced-algorithms-and-complexity)

## Course 6 - [Genome Assembly Programming Challenge](https://www.coursera.org/learn/assembling-genomes)

### Final project: Implementing an Assembler

`course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler.go`
is a single-file Go assembler for the final project.

Input is either:

- a read count followed by one read per line, or
- a read-pair count followed by `READ1|READ2|d` lines.

Output is FASTA. Run it with:

```bash
go run course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/implementing_an_assembler.go \
  < course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler/test/dataset.txt
```

The implementation uses solid k-mers, implicit de Bruijn unitigs, low-coverage
tip removal, Bloom-based read correction for error-prone data, paired-read
orientation normalization, and chain scaffolding. For paired data, the second
mate is reverse-complemented when the library is forward-reverse so assembly
stays on one strand; forward-forward libraries are left unchanged. Scaffolds
are built from unique best links and chained through multiple contigs. On local
E. coli proxies it produced NGA50 47,576 for 1.4M error-prone single reads,
47,120 for 1.4M error-free single reads, 73,696 for 700,000 error-prone
forward-forward read-pairs, 96,871 for 700,000 error-prone forward-reverse
read-pairs, 132,568 for 700,000 error-free forward-forward read-pairs, and
139,619 for 700,000 error-free forward-reverse read-pairs. The slowest
correction proxy took 176 s and 1.8 GB, inside the grader's observed 450 s and
4 GB limits. Hidden dataset quality is still not guaranteed.

`test/dataset.txt` and `test/genome.txt` are optional local fixtures. They are
not committed. The visible-dataset test skips when they are absent:

```bash
go test -run TestVisibleDatasetAssembly ./course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler -v
```

## References

- Algorithms (1st Edition). Sanjoy Dasgupta, Christos Papadimitriou, and Umesh Vazirani {1}
- Introduction to Algorithms (3rd Edition). Thomas H. Cormen, Charles E. Leiserson, Ronald L. Rivest, Clifford Stein {2}
- Bioinformatics Algorithms: An Active Learning Approach. Phillip Compeau, Pavel Pevzner {3}
- [Notes on Data Structures and Programming Techniques (CPSC 223, Spring 2022)](https://www.cs.yale.edu/homes/aspnes/classes/223/notes.html) {4}
- Advanced Dynamic Programming Lecture Notes. Jeff Erickson {5}
- Algorithms on Strings, Trees and Sequences: Computer Science and Computational Biology (1st Edition). Dan Gusfield. {6}
- Algorithm design (1st Edition). Addison-Wesley, Jon M. Kleinberg and Eva Tardos. {7}
- [Learn Python In Y Minutes](https://learnxinyminutes.com/docs/python/) {8}
- [Discrete Mathematics Lecture notes by László Lovász, Yale University](https://cims.nyu.edu/~regev/teaching/discrete_math_fall_2005/dmbook.pdf) {9}
- [Learn Go In Y Minutes](https://learnxinyminutes.com/go/) {10}

<details>
<summary><strong>Reading List Chapters</strong></summary>

- Course 0: Prerequisites
  - {8}: Basic Python Programming
  - {9}: Basic Discrete Mathematics
  - {10}: Basic Go Programming
- Course 1: Algorithmic Toolbox
  - Week 1: Intro
  - Week 2: Warm-up
    - {1}: Section 0.2
    - {1}: Section 1.2.3
    - {2}: Section 31.1
    - {1}: Section 0.3
  - Week 3: Greedy Algorithms
  - Week 4: Divide and Conquer
    - {1}: Section 2.1
    - {1}: Section 2.2
    - {1}: Section 2.3
    - {2}: Chapter 7
  - Week 5: Dynamic Programming 1
    - {3}: Chapter 5: Intro to Dynamic Programming (P.236)
    - {1}: Section 6.3
    - {3}: Chapter 5: How Do We Compare Biological Sequences (optional)
    - {5}: Advanced Dynamic Programming Notes (optional)
  - Week 6: Dynamic Programming 2
    - {1}: Section 6.4
- Course 2: Data Structures
  - Week 1: Basic Data Structures
    - {2}: Section 10.2
    - {2}: Section 10.1
    - {2}: Section 10.4
  - Week 2: Dynamic Arrays and Amortized Analysis
    - {2}: Chapter 17
  - Week 3: Priority Queues and Disjoint Sets
    - {2}: Chapter 6
    - {2}: Section 6.4
    - {2}: Section 21.1
    - {2}: Section 21.2
    - {1}: Section 5.1.4
  - Week 4: Hash Tables
    - {1}: Section 1.5.1
    - {2}: Section 11.1
    - {2}: Section 11.2
    - {1}: Section 1.5
    - {2}: Section 11.3
    - {2}: Section 32.1
    - {2}: Section 32.2
  - Week 5: Binary Search Trees
    - {2}: Chapter 12
    - {4}: Section 5.11.1
    - {4}: Section 5.11.2
    - {o}: [AVL Trees](https://en.wikipedia.org/wiki/AVL_tree)
  - Week 6: Binary Search Trees 2
    - {2}: Section 14.1
    - {2}: Section 14.2
    - {4}: Section 5.11.6
    - {o}: [Self-Adjusting Binary Search Trees](https://www.cs.cmu.edu/~sleator/papers/self-adjusting.pdf)
- Course 3: Algorithms on Graphs
  - Week 1: Decomposition of Graphs 1
    - {1}: Section 3.1
    - {1}: Section 3.2
  - Week 2: Decomposition of Graphs 2
    - {1}: Section 3.3
    - {1}: Section 3.4
  - Week 3: Paths in Graphs 1
    - {1}: Section 4.1
    - {1}: Section 4.2
  - Week 4: Paths in Graphs 2
    - {1}: Section 4.3
    - {1}: Section 4.4
    - {1}: Section 4.6
  - Week 5: Minimum Spanning Trees
    - {1}: Section 5.1
  - Week 6: Advanced Shortest Paths Project (optional)
- Course 4: Algorithms on Strings
  - Week 1: Suffix Trees
    - {3}: Chapter 3: How Do We Assemble Genomes
    - {3}: Chapter 9: How Do We Locate Disease-Causing Mutations
  - Week 2: Burrows Wheeler Transform and Suffix Arrays
    - {3}: Chapter 3: How Do We Assemble Genomes
    - {3}: Chapter 9: How Do We Locate Disease-Causing Mutations
  - Week 3: Knuth-Morris-Pratt Algorithm
    - {6}: Chapter 1
    - {6}: Section 2.3
    - {6}: Section 3.3
  - Week 4: Constructing Suffix Arrays and Suffix Trees
    - {6}: Chapter 4
    - {3}: ...
- Course 5: Advanced Algorithms and Complexity
  - Week 1: Flows in Networks
    - {1}: Section 7.1
    - {1}: Section 7.2
    - {1}: Section 7.3
    - {7}: Chapter 7
    - {2}: Chapter 26
  - Week 2: Linear Programming
    - {1}: Chapter 7
    - {2}: Chapter 29
  - Week 3: NP-Complete Problems
    - {1}: Chapter 8
    - {7}: Chapter 8
    - {2}: Chapter 34
  - Week 4: Coping with NP-Completeness
    - {1}: Chapter 9
    - {7}: Chapter 10
    - {3}: Chapter 35
  - Week 5: Streaming Algorithms (optional)
- Course 6: Genome Assembly Programming Challenge
  - Week 1: The 2011 European E. coli Outbreak
  - Week 2: Assembling Genomes Using de Bruijn Graphs
  - Week 3: Genome Assembly Faces Real Sequencing Data

</details>

## Problems

Weekly assignments can be found at `/resources`\
Fixture files are `test/<name>` plus `test/<name>.a`. A single runner checks them
without writing intermediate files:

```sh
make help          # list targets and descriptions
make test          # Go unit tests
make fixtures      # Go IO fixtures
make fixtures-python

# One problem, verbose
go run ./internal/checkfixtures -only 'week_2_linear_programming' -v

# C/C++: compile a binary named solution in each problem directory
go run ./internal/checkfixtures -bin ./solution
```

Each Go solution has a matching `_test.go` file with fixture checks and edge cases. Every Go solution file is self-contained: shared helpers used by several assignments are inlined so that the file can be submitted to Coursera as a single source file. The `internal/` packages are kept for tests and are not required by a submission.

### Generated fixtures

Some problems have extra `test/<number>` and `test/<number>.a` pairs produced by
`internal/genfixtures`. They are best-effort coverage fixtures, not official
Coursera test data. The expected output is the Go solution's stdout; inputs are
mutated from the existing fixtures and kept only when they reach new reachable
code. The generator applies shape checks and records Python disagreements as
warnings; Go remains the oracle. An input can still fall outside the original
problem constraints, and problems with several valid answers may have
implementation-specific `.a` files.

See [`docs/generated-fixtures.md`](docs/generated-fixtures.md) for per-problem
coverage and caveats. Regenerate with:

```sh
go run ./internal/genfixtures -write
```

Use `-only <regexp>` to target one problem and `-max-new`, `-rounds`, `-batch`
to change the effort. `internal/genfixtures/README.md` documents the pipeline.

The runner compares whitespace-insensitively by default and allows a relative
tolerance for floating-point tokens. Problems with non-unique outputs are listed
explicitly in `internal/checkfixtures/modes.go`; they use rotation/set comparison
or a smoke check, and smoke results are reported as `SKIP`, not `PASS`.
`internal/checkfixtures/checker.go` adds semantic checks for topological order,
binary de Bruijn cycles, non-shared substrings, and longest common substrings.
The remaining smoke checks are `gsm_network`, `cleaning_apartment`, and
`suffix_tree_from_array`. `internal/checkfixtures/README.md` documents the modes,
checkers, and flags.
