# checkfixtures

Run every `test/<name>` input against `test/<name>.a`.

```sh
go run ./internal/checkfixtures
go run ./internal/checkfixtures -only 'week_4_divide_and_conquer' -v
```

The runner builds each Go solution once, runs all its fixtures, compares output
without writing files, prints failures, and exits non-zero on failure.

## Languages

```sh
# Python
go run ./internal/checkfixtures -lang python

# C/C++: compile a binary named `solution` in every problem directory, then
go run ./internal/checkfixtures -bin ./solution
```

`-bin` is resolved relative to each problem directory unless it is absolute.

## Comparison

- Default: whitespace-insensitive tokens, blank lines ignored, trailing newline
  irrelevant.
- `-float 1e-6`: allow relative error for decimal tokens (default).
- `-exact`: byte-for-byte comparison.
- A few problems have non-unique outputs (topological order, Eulerian cycle,
  CNF formula, suffix tree edge order, genome rotation). `modes.go` lists them:
  - `rotation-lines`: accept rotations of each line
  - `unordered-lines`: compare line sets
  - `unordered-tokens`: compare token sets
  - `smoke`: run the solution and check exit status only (available, currently
    unused; results are reported as `SKIP`, not `PASS`)

`checker.go`, `checker_more.go`, `checker_cnf.go`, and
`checker_suffix_tree.go` implement semantic checkers where output is not unique:

| Problem | Check |
|---|---|
| `5_collecting_signatures` | Every segment is covered; point count matches the expected minimum |
| `6_maximum_number_of_prizes` | `k` distinct positive summands sum to `n`; `k` matches the expected maximum |
| `2_primitive_calculator` | Sequence starts at 1, ends at `n`, uses allowed steps; length matches the expected minimum |
| `1_make_heap` | Swaps produce a valid min-heap |
| `2_order_of_courses` | Output is a permutation and every course edge points forward |
| `5_longest_common_substring` | Each `i j length` triple has the expected length and equal substrings |
| `5_non_shared_substring` | Answer has the expected shortest length, occurs in the first string, absent from the second |
| `4_suffix_tree_from_array` | Edge labels and order match the suffix tree built from the suffix and LCP arrays |
| `2_airline_crews` | Valid matching with the expected maximum cardinality |
| `2_diet`, `3_ad_allocation` | Solution is feasible and its objective matches the reference optimum |
| `1_gsm_network` | CNF is satisfiable exactly when the graph is 3-colorable |
| `2_cleaning_apartment` | CNF is satisfiable exactly when a Hamiltonian path exists (`n<=20`; reference fallback above) |
| `3_budget_allocation` | CNF is satisfiable exactly when the 0/1 constraints are feasible (`m<=20`; reference fallback above) |
| `1_circuit_design` | Assignment satisfies every 2-CNF clause |
| `3_school_bus` | Tour visits every vertex, uses real edges, and has the reference cost |
| `4_reschedule_exams` | Valid recoloring with no same-color edge and no unchanged vertex |
| `1_puzzle_assembly` | Every piece is used once and adjacent colors match |
| `1_network_circulation` | Flow respects lower/upper bounds and conservation |
| `2_optimal_kmer_size` | De Bruijn graph is a single cycle |
| `3_k_universal_circular_string` | Output is a binary de Bruijn cycle containing every k-mer |

All 107 problems now have a semantic checker or an explicit comparison mode;
`make fixtures` passes 462/462 with no `SKIP`. `smoke` remains available for a
future problem but is currently unused.

## Flags

| Flag | Default | Meaning |
|---|---:|---|
| `-root` | `.` | repository root |
| `-only` | empty | regexp filter on problem paths |
| `-lang` | `go` | `go` or `python` |
| `-bin` | empty | executable relative to each problem directory |
| `-jobs` | `4` | parallel test runs |
| `-timeout` | `10s` | per-run timeout |
| `-memory-mb` | `3072` | address-space limit via `prlimit` |
| `-exact` | `false` | byte-exact output |
| `-float` | `1e-6` | relative tolerance for decimal tokens |
| `-v` | `false` | print PASS/SKIP lines |
| `-list` | `false` | list tests without running |

## Why not `out` files

`out` files were committed command logs, one Python command per fixture. They
were stale, bash-only, Python-only, and had no exit status. The runner uses the
`test/<name>` + `test/<name>.a` data directly and supports Go, Python, C, and
C++.
