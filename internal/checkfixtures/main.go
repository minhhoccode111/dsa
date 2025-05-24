package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type TestCase struct {
	Name       string
	InputPath  string
	AnswerPath string
	Input      []byte
	Answer     []byte
}

type Problem struct {
	Rel      string
	Dir      string
	Solution string
	Python   string
	Tests    []TestCase
}

type Options struct {
	Root     string
	Only     string
	Lang     string
	Bin      string
	Jobs     int
	Timeout  time.Duration
	MemoryMB int64
	Exact    bool
	Float    float64
	Verbose  bool
	List     bool
}

type commandRunner struct {
	Command string
	Args    []string
	Env     []string
}

func main() {
	var opts Options
	flag.StringVar(&opts.Root, "root", ".", "repository root")
	flag.StringVar(&opts.Only, "only", "", "regular expression matched against problem paths")
	flag.StringVar(&opts.Lang, "lang", "go", "solution language: go or python")
	flag.StringVar(&opts.Bin, "bin", "", "executable relative to each problem directory; overrides -lang")
	flag.IntVar(&opts.Jobs, "jobs", 4, "parallel test runs")
	flag.DurationVar(&opts.Timeout, "timeout", 10*time.Second, "timeout for one solution run")
	flag.Int64Var(&opts.MemoryMB, "memory-mb", 3072, "address-space limit for one solution run")
	flag.BoolVar(&opts.Exact, "exact", false, "require byte-exact output")
	flag.Float64Var(&opts.Float, "float", 1e-6, "relative tolerance for floating-point tokens; 0 disables")
	flag.BoolVar(&opts.Verbose, "v", false, "print passing tests")
	flag.BoolVar(&opts.List, "list", false, "list tests without running them")
	flag.Parse()

	if opts.Jobs < 1 {
		opts.Jobs = 1
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 10 * time.Second
	}

	var filter *regexp.Regexp
	if opts.Only != "" {
		var err error
		filter, err = regexp.Compile(opts.Only)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -only pattern: %v\n", err)
			os.Exit(2)
		}
	}

	problems, err := discover(opts.Root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover: %v\n", err)
		os.Exit(1)
	}

	selected := make([]Problem, 0, len(problems))
	totalTests := 0
	for _, problem := range problems {
		if filter != nil && !filter.MatchString(problem.Rel) {
			continue
		}
		selected = append(selected, problem)
		totalTests += len(problem.Tests)
	}
	if len(selected) == 0 {
		fmt.Fprintln(os.Stderr, "no problems selected")
		os.Exit(1)
	}

	if opts.List {
		for _, problem := range selected {
			fmt.Printf("%s (%d tests)\n", problem.Rel, len(problem.Tests))
			for _, test := range problem.Tests {
				fmt.Printf("  %s\n", test.Name)
			}
		}
		fmt.Printf("problems: %d, tests: %d\n", len(selected), totalTests)
		return
	}

	ctx := context.Background()
	started := time.Now()
	passed := 0
	failed := 0
	skipped := 0
	prepareFailed := 0
	for _, problem := range selected {
		result := checkProblem(ctx, problem, opts)
		if result.prepareErr != nil {
			prepareFailed++
			fmt.Printf("FAIL %s: %v\n", problem.Rel, result.prepareErr)
			continue
		}
		for _, outcome := range result.outcomes {
			if outcome.skipped {
				skipped++
				if opts.Verbose {
					fmt.Printf("SKIP %s/%s: %s\n", problem.Rel, outcome.test.Name, outcome.detail)
				}
				continue
			}
			if outcome.ok {
				passed++
				if opts.Verbose {
					fmt.Printf("PASS %s/%s\n", problem.Rel, outcome.test.Name)
				}
				continue
			}
			failed++
			fmt.Printf("FAIL %s/%s: %s\n", problem.Rel, outcome.test.Name, outcome.detail)
		}
	}

	elapsed := time.Since(started).Round(time.Millisecond)
	fmt.Printf("problems: %d, tests: %d, passed: %d, failed: %d, skipped: %d, prepare-failed: %d, elapsed: %s\n",
		len(selected), totalTests, passed, failed, skipped, prepareFailed, elapsed)
	if failed > 0 || prepareFailed > 0 {
		os.Exit(1)
	}
}

type problemResult struct {
	outcomes   []testOutcome
	prepareErr error
}

type testOutcome struct {
	test    TestCase
	ok      bool
	skipped bool
	detail  string
}

func checkProblem(ctx context.Context, problem Problem, opts Options) problemResult {
	if len(problem.Tests) == 0 {
		return problemResult{}
	}
	runner, cleanup, err := prepareRunner(problem, opts)
	if err != nil {
		return problemResult{prepareErr: err}
	}
	defer cleanup()

	mode := modeFor(problem.Rel)
	checker := checkerFor(problem.Rel)
	outcomes := make([]testOutcome, len(problem.Tests))
	sem := make(chan struct{}, opts.Jobs)
	var wg sync.WaitGroup
	for i := range problem.Tests {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			test := problem.Tests[index]
			run := runCommand(ctx, runner, test.Input, problem.Dir, opts)
			outcome := testOutcome{test: test}
			switch {
			case run.TimedOut:
				outcome.detail = "timeout"
			case run.OutputLimited:
				outcome.detail = "output limit exceeded"
			case run.Err != nil:
				outcome.detail = "run failed: " + run.Err.Error()
			case checker != nil:
				if ok, detail := checker(test.Input, test.Answer, run.Stdout); !ok {
					outcome.detail = "checker: " + detail
				} else {
					outcome.ok = true
				}
			case mode == modeSmoke:
				outcome.skipped = true
				outcome.detail = "non-unique output; no semantic checker"
			default:
				if ok, detail := compareOutput(run.Stdout, test.Answer, opts, mode); !ok {
					outcome.detail = detail
				} else {
					outcome.ok = true
				}
			}
			outcomes[index] = outcome
		}(i)
	}
	wg.Wait()
	return problemResult{outcomes: outcomes}
}

func prepareRunner(problem Problem, opts Options) (commandRunner, func(), error) {
	if opts.Bin != "" {
		binary := opts.Bin
		if !filepath.IsAbs(binary) {
			binary = filepath.Join(problem.Dir, binary)
		}
		if info, err := os.Stat(binary); err != nil {
			return commandRunner{}, func() {}, fmt.Errorf("binary %s: %w", binary, err)
		} else if info.IsDir() {
			return commandRunner{}, func() {}, fmt.Errorf("binary %s is a directory", binary)
		}
		return commandRunner{Command: binary}, func() {}, nil
	}

	switch opts.Lang {
	case "go":
		if problem.Solution == "" {
			return commandRunner{}, func() {}, fmt.Errorf("no Go solution")
		}
		tmpDir, err := os.MkdirTemp("", "dsa-checkfixtures-")
		if err != nil {
			return commandRunner{}, func() {}, err
		}
		binary := filepath.Join(tmpDir, "solution")
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		cmd := commandContext(ctx, goToolPath(), "build", "-o", binary, filepath.Base(problem.Solution))
		cmd.Dir = problem.Dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			os.RemoveAll(tmpDir)
			return commandRunner{}, func() {}, fmt.Errorf("build: %v: %s", err, strings.TrimSpace(string(output)))
		}
		return commandRunner{Command: binary}, func() { os.RemoveAll(tmpDir) }, nil
	case "python":
		if problem.Python == "" {
			return commandRunner{}, func() {}, fmt.Errorf("no Python solution")
		}
		py, err := findPython()
		if err != nil {
			return commandRunner{}, func() {}, err
		}
		return commandRunner{
			Command: py,
			Args:    []string{filepath.Base(problem.Python)},
			Env:     []string{"PYTHONDONTWRITEBYTECODE=1"},
		}, func() {}, nil
	default:
		return commandRunner{}, func() {}, fmt.Errorf("unknown -lang %q (use go, python, or -bin)", opts.Lang)
	}
}

func discover(root string) ([]Problem, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var problems []Problem
	err = filepath.WalkDir(absRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() || entry.Name() != "test" || path == absRoot {
			return nil
		}
		dir := filepath.Dir(path)
		rel, err := filepath.Rel(absRoot, dir)
		if err != nil {
			return err
		}
		problem := Problem{Rel: rel, Dir: dir}
		problem.Solution, problem.Python = findSolutions(dir)
		problem.Tests, err = findTests(path)
		if err != nil {
			return err
		}
		if len(problem.Tests) > 0 {
			problems = append(problems, problem)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(problems, func(i, j int) bool { return problems[i].Rel < problems[j].Rel })
	return problems, nil
}

func findSolutions(dir string) (string, string) {
	goFiles, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	var solution string
	for _, file := range goFiles {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		solution = file
		break
	}
	pyFiles, _ := filepath.Glob(filepath.Join(dir, "*.py"))
	var python string
	for _, file := range pyFiles {
		if filepath.Base(file) == "starter.py" {
			continue
		}
		python = file
		break
	}
	return solution, python
}

func findTests(testDir string) ([]TestCase, error) {
	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, err
	}
	var tests []TestCase
	for _, entry := range entries {
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".a") || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		inputPath := filepath.Join(testDir, entry.Name())
		answerPath := inputPath + ".a"
		if _, err := os.Stat(answerPath); err != nil {
			continue
		}
		input, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, err
		}
		answer, err := os.ReadFile(answerPath)
		if err != nil {
			return nil, err
		}
		tests = append(tests, TestCase{
			Name:       entry.Name(),
			InputPath:  inputPath,
			AnswerPath: answerPath,
			Input:      input,
			Answer:     answer,
		})
	}
	sort.Slice(tests, func(i, j int) bool { return tests[i].Name < tests[j].Name })
	return tests, nil
}
