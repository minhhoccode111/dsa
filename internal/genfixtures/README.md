# genfixtures

Best-effort IO fixture generator for the Go solutions.

It treats each Go solution as the oracle, mutates the existing `test/` inputs,
and keeps mutations that increase code coverage. Expected outputs are the
solution's stdout. It does not parse the original problem statements.

## Run

```sh
go run ./internal/genfixtures -write
```

Dry run (no files written):

```sh
go run ./internal/genfixtures -only 'week_2_linear_programming'
```

Useful flags:

| Flag            | Default | Meaning                                                     |
| --------------- | ------: | ----------------------------------------------------------- |
| `-write`        | `false` | Write new `test/N` and `test/N.a` files                     |
| `-only`         |   empty | Regexp matched against each problem path                    |
| `-rounds`       |     `6` | Mutation rounds per problem                                 |
| `-batch`        |    `64` | Candidates per round                                        |
| `-max-new`      |    `10` | Maximum new fixtures per problem                            |
| `-jobs`         |     `4` | Parallel candidate runs                                     |
| `-timeout`      |    `2s` | Timeout for one solution run                                |
| `-memory-mb`    |  `3072` | Address-space limit per run via `prlimit`                   |
| `-python-check` |  `true` | Run the sibling Python solution as a secondary check        |
| `-python-veto`  | `false` | Skip candidates when the Python solution fails or disagrees |

Without `-write`, the tool still runs the full pipeline and prints a summary.
With `-write`, it also updates `docs/generated-fixtures.md`.

## Pipeline

1. Discover each problem directory with a single non-test Go file and a `test/`
   directory. Seeds are inputs that have a matching `.a` file.
2. Build an instrumented binary with `go build -cover` and a normal binary.
3. Parse the solution AST and keep only functions reachable from `main`.
   `if err != nil { return }` input-error guards are excluded from the target.
4. Run the seeds with `GOCOVERDIR` and record per-seed block coverage.
5. Mutate seeds (numbers, lines, tokens, bytes, splices, truncation) and run
   candidates. A candidate is kept only when it adds target coverage.
6. Before writing, require deterministic output, stable shape (token/line count,
   numeric domain and line kinds derived from the seeds), and no ignored trailing
   token. The Python solution is run too; its failure or disagreement is recorded
   as a warning unless `-python-veto` is set.
7. Write the selected inputs after the highest existing numeric fixture name.
   Existing files are never overwritten.

## Why the gates exist

A raw mutator quickly produces inputs that still exit successfully but are
meaningless: a count changed without its payload, extra tokens ignored by
`fmt.Fscan`, or a DNA string replaced by digits. The shape gates and the
ignored-suffix check remove many of those. The Python cross-check is not
authoritative: a mismatch can be a Python bug, a non-unique answer, or a real
difference. By default it only records a warning because the Go solutions are
the trusted oracle. Use `-python-veto` if you want Python to block candidates.

## Limitations

- No knowledge of Coursera constraints; semantic validity is not guaranteed.
- Go exposes block coverage, not condition/branch coverage. A `&&`/`||`
  sub-condition is not counted separately.
- Dead inline helper functions are excluded by the reachability filter, so
  "100%" means every reachable non-error block found by the instrumented run.
- Problems with several valid outputs or floating-point answers may produce
  implementation-specific `.a` files.
- Solutions that consume input line-by-line may need per-problem generation;
  the generic mutator is intentionally conservative.
