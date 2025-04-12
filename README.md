# Data Structures and Algorithms Specialization Coursera

**Disclaimer**: The code is terrible and buggy, it has been a few months since I worked on it. Use at your own risk.

[Data Structures and Algorithms Specialization](https://www.coursera.org/specializations/data-structures-algorithms)

This repository contains both Python and Go implementations of the assignments. The Go solutions live next to the Python ones, include tests, and can be run from the repository root.

## Languages

- Python 3
- Go 1.21+ (`go.mod` uses the `dsa` module)

## Course 1 - [Algorithmic Toolbox](https://www.coursera.org/learn/algorithmic-toolbox)

## Course 2 - [Data Structures](https://www.coursera.org/learn/data-structures)

## Course 3 - [Algorithms on Graphs](https://www.coursera.org/learn/algorithms-on-graphs)

## Course 4 - [Algorithms on Strings](https://www.coursera.org/learn/algorithms-on-strings)

## Course 5 - [Advanced Algorithms and Complexity](https://www.coursera.org/learn/advanced-algorithms-and-complexity)

## Course 6 - [Genome Assembly Programming Challenge](https://www.coursera.org/learn/assembling-genomes)

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

<details>
<summary><strong>Reading List Chapters</strong></summary>

- Course 0: Prerequisites
  - {8}: Basic Python Programming
  - {9}: Basic Discrete Mathematics
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
Solutions and test files can be tested using the following commands (the commands write their output to the `out` file in the problem directory):

```sh
# Python
diff -b <(python3 solution.py < test/1) <(cat test/1.a)

# Go (from the problem directory)
diff -b <(go run solution.go < test/2) <(cat test/2.a)

# Run every Go test from the repository root
go test ./...
```

Each Go solution has a matching `_test.go` file with fixture checks and edge cases. Every Go solution file is self-contained: shared helpers used by several assignments are inlined so that the file can be submitted to Coursera as a single source file. The `internal/` packages are kept for tests and are not required by a submission.

<details>
<summary><strong>Command explain</strong></summary>

This command compares the output of the Python script (`solution.py`) when run with the input file (`test/1`) against the expected output stored in `test/1.a`. The `-b` flag ensures that differences in whitespace are ignored during the comparison.

- `diff`: A command-line utility used to compare two files line by line and output the differences between them. If there are no differences, it produces no output.
- `-b`: An option for `diff` that ignores differences in the amount of whitespace. For example, differences in the number of spaces or tabs will be overlooked.
- `<(python3 solution.py < test/1)`: A process substitution that runs `python3 solution.py` with the input redirected from the file `test/1`. This executes the Python script `solution.py` using `test/1` as its input and provides its output as a temporary file-like stream to `diff`.
- `<(cat test/1.a)`: Another process substitution that runs `cat test/1.a`, effectively reading the contents of the file `test/1.a` and providing it as another temporary file-like stream to `diff`.
- `solution.py`: your solution file
- `test/1`: provided input file
- `test/1.a`: expected output file

The Go variant is the same: `solution.go` is the program and `go run solution.go < test/2` compiles and runs it. From the repository root, `go test ./...` runs the full Go test suite, including fixture comparisons and edge cases for every solution.

</details>
