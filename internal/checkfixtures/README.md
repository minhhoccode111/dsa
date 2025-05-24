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
  - `smoke`: run the solution and check exit status only
  - `modeSmoke` results are reported as `SKIP`, not `PASS`.

`checker.go` implements semantic checkers where output is not unique:

| Problem | Check |
|---|---|
| `2_order_of_courses` | Output is a permutation and every course edge points forward |
| `3_k_universal_circular_string` | Output is a binary de Bruijn cycle containing every k-mer |
| `5_non_shared_substring` | Answer has the expected shortest length, occurs in the first string, absent from the second |
| `5_longest_common_substring` | Each `i j length` triple has the expected length and equal substrings |

`SKIP` is now limited to `gsm_network`, `cleaning_apartment`, and
`suffix_tree_from_array`; those need a CNF/tree checker before they can be
treated as regression tests.

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
