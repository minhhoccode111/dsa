package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func main() {
	var opts Options
	var root string
	var only string
	var memoryMB int64
	flag.StringVar(&root, "root", ".", "repository root")
	flag.BoolVar(&opts.Write, "write", false, "write generated fixtures into test/ directories")
	flag.StringVar(&only, "only", "", "regular expression matched against problem paths")
	flag.IntVar(&opts.Jobs, "jobs", 4, "parallel candidate runs")
	flag.IntVar(&opts.Rounds, "rounds", 6, "mutation rounds per problem")
	flag.IntVar(&opts.Batch, "batch", 64, "candidate inputs per round")
	flag.IntVar(&opts.MaxNew, "max-new", 10, "maximum generated fixtures per problem")
	flag.DurationVar(&opts.Timeout, "timeout", 2*time.Second, "timeout for one solution run")
	flag.IntVar(&opts.MaxInput, "max-input", 1<<20, "maximum candidate input bytes")
	flag.IntVar(&opts.MaxOutput, "max-output", 1<<20, "maximum captured output bytes")
	flag.Int64Var(&memoryMB, "memory-mb", 3072, "address-space limit for one solution run")
	flag.BoolVar(
		&opts.PythonCheck,
		"python-check",
		true,
		"compare generated outputs with the sibling Python solution",
	)
	flag.BoolVar(
		&opts.PythonVeto,
		"python-veto",
		false,
		"skip candidates when the Python solution fails or disagrees",
	)
	flag.BoolVar(&opts.Verbose, "verbose", false, "print every problem result")
	flag.Parse()

	opts.MemoryBytes = memoryMB << 20
	if opts.Jobs < 1 {
		opts.Jobs = 1
	}
	if opts.Rounds < 0 || opts.Batch < 1 || opts.MaxNew < 1 {
		fmt.Fprintln(os.Stderr, "rounds must be >= 0, batch and max-new must be >= 1")
		os.Exit(2)
	}

	var filter *regexp.Regexp
	if only != "" {
		var err error
		filter, err = regexp.Compile(only)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -only pattern: %v\n", err)
			os.Exit(2)
		}
	}

	if !isGoToolAvailable() {
		fmt.Fprintln(os.Stderr, "go tool is not available")
		os.Exit(1)
	}

	problems, err := discover(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover: %v\n", err)
		os.Exit(1)
	}

	selected := make([]*Problem, 0, len(problems))
	for _, problem := range problems {
		if filter == nil || filter.MatchString(problem.Rel) {
			selected = append(selected, problem)
		}
	}
	if len(selected) == 0 {
		fmt.Fprintln(os.Stderr, "no problems selected")
		os.Exit(1)
	}

	ctx := context.Background()
	results := make([]Result, 0, len(selected))
	started := time.Now()
	for i, problem := range selected {
		rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)*7919))
		result := processProblem(ctx, problem, opts, rng)
		results = append(results, result)
		if opts.Verbose || opts.Write {
			fmt.Fprintf(os.Stderr, "[%d/%d] %-70s %-14s target %.1f%% -> %.1f%% (+%d)\n",
				i+1, len(selected), shortPath(problem.Rel), result.Status,
				result.TargetBefore, result.TargetAfter, result.Generated)
			if opts.Verbose {
				for j, block := range result.Uncovered {
					if j >= 10 {
						break
					}
					fmt.Fprintf(os.Stderr, "    uncovered: %s\n", block.String())
				}
			}
		}
	}

	summary := summarize(results)
	if opts.Write {
		reportPath := writeReport(root, results)
		fmt.Printf("report: %s\n", reportPath)
	}
	fmt.Printf("problems: %d, generated fixtures: %d, elapsed: %s\n",
		len(results), summary.Generated, time.Since(started).Round(time.Second))
	fmt.Printf(
		"status: ok=%d partial=%d already-covered=%d no-new=%d no-seed=%d weak-seed=%d other=%d\n",
		summary.OK,
		summary.Partial,
		summary.Already,
		summary.NoNew,
		summary.NoSeed,
		summary.Weak,
		summary.Other,
	)
}

func shortPath(path string) string {
	if len(path) <= 70 {
		return path
	}
	return "..." + path[len(path)-67:]
}

type summary struct {
	Generated int
	OK        int
	Partial   int
	Already   int
	NoNew     int
	NoSeed    int
	Weak      int
	Other     int
}

func summarize(results []Result) summary {
	var out summary
	for _, result := range results {
		out.Generated += result.Generated
		switch result.Status {
		case "ok":
			out.OK++
		case "partial":
			out.Partial++
		case "already-covered":
			out.Already++
		case "no-new":
			out.NoNew++
		case "no-seed":
			out.NoSeed++
		case "weak-seed":
			out.Weak++
		default:
			out.Other++
		}
	}
	return out
}
