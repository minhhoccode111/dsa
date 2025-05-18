package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Options struct {
	Write       bool
	Jobs        int
	Rounds      int
	Batch       int
	MaxNew      int
	Timeout     time.Duration
	MaxInput    int
	MaxOutput   int
	MemoryBytes int64
	PythonCheck bool
	PythonVeto  bool
	Verbose     bool
}

type generatedFixture struct {
	Input  []byte
	Output []byte
	Name   string
}

type Result struct {
	Rel            string
	Status         string
	Existing       int
	Generated      int
	RawBefore      float64
	RawAfter       float64
	TargetBefore   float64
	TargetAfter    float64
	TargetBlocks   int
	Uncovered      []Block
	Warnings       []string
	GeneratedFiles []string
	BuildErr       string
}

type candidateResult struct {
	Input []byte
	Run   RunResult
	Set   CoverSet
}

type lineKind int

const (
	lineKindNone lineKind = iota
	lineKindNumeric
	lineKindString
	lineKindMixed
)

type lineRelation int

const (
	lineRelationNone lineRelation = iota
	lineRelationPlusOne
	lineRelationSame
)

type seedShape struct {
	MinTokens           int
	MaxTokens           int
	MinLines            int
	MaxLines            int
	AllowZero           bool
	AllowNegative       bool
	AllAnswersNonEmpty  bool
	AllowNegativeAnswer bool
	StringChars         map[byte]bool
	LineKinds           []lineKind
	LineTokenCounts     []int
	LineStringLen       []int
	LineMax             map[int]int64
	LineRelation        lineRelation
	EmptyAnswerFirst    map[int64]bool
	NonEmptyAnswerFirst map[int64]bool
}

var numericToken = regexp.MustCompile(`^-?[0-9]+$`)
var tokenPattern = regexp.MustCompile(`\S+`)

func shapeFromSeeds(seeds []Seed) (seedShape, bool) {
	shape := seedShape{
		MinTokens:           -1,
		MinLines:            -1,
		AllAnswersNonEmpty:  true,
		StringChars:         make(map[byte]bool),
		LineMax:             make(map[int]int64),
		EmptyAnswerFirst:    make(map[int64]bool),
		NonEmptyAnswerFirst: make(map[int64]bool),
	}
	found := false
	relationKnown := true
	relationSeen := false
	relationPlusOne := true
	relationSame := true
	for _, seed := range seeds {
		tokens := strings.Fields(string(seed.Input))
		if len(tokens) == 0 {
			continue
		}
		found = true
		for _, token := range tokens {
			if numericToken.MatchString(token) {
				continue
			}
			for i := 0; i < len(token); i++ {
				shape.StringChars[token[i]] = true
			}
		}
		lines := lineCount(seed.Input)
		inputLines := splitInputLines(seed.Input)
		var seedFirst int64
		seedHasFirst := false
		if len(inputLines) > 0 {
			firstFields := strings.Fields(inputLines[0])
			if len(firstFields) > 0 && numericToken.MatchString(firstFields[0]) {
				firstValue, err := strconv.ParseInt(firstFields[0], 10, 64)
				if err == nil {
					seedFirst = firstValue
					seedHasFirst = true
					relationSeen = true
					if int64(lines) != firstValue+1 {
						relationPlusOne = false
					}
					if int64(lines) != firstValue {
						relationSame = false
					}
				} else {
					relationKnown = false
				}
			} else {
				relationKnown = false
			}
		}
		if shape.MinTokens < 0 || len(tokens) < shape.MinTokens {
			shape.MinTokens = len(tokens)
		}
		if len(tokens) > shape.MaxTokens {
			shape.MaxTokens = len(tokens)
		}
		if shape.MinLines < 0 || lines < shape.MinLines {
			shape.MinLines = lines
		}
		if lines > shape.MaxLines {
			shape.MaxLines = lines
		}
		answerEmpty := len(strings.TrimSpace(string(seed.Answer))) == 0
		if answerEmpty {
			shape.AllAnswersNonEmpty = false
			if seedHasFirst {
				shape.EmptyAnswerFirst[seedFirst] = true
			}
		} else if seedHasFirst {
			shape.NonEmptyAnswerFirst[seedFirst] = true
		}
		for _, match := range numberPattern.FindAll(seed.Answer, -1) {
			if len(match) > 0 && match[0] == '-' {
				shape.AllowNegativeAnswer = true
			}
		}
		for index, line := range splitInputLines(seed.Input) {
			kind := classifyLine(line)
			if kind == lineKindNone {
				continue
			}
			for len(shape.LineKinds) <= index {
				shape.LineKinds = append(shape.LineKinds, lineKindNone)
			}
			shape.LineKinds[index] = mergeLineKind(shape.LineKinds[index], kind)
			for len(shape.LineTokenCounts) <= index {
				shape.LineTokenCounts = append(shape.LineTokenCounts, -1)
			}
			count := len(strings.Fields(line))
			switch {
			case shape.LineTokenCounts[index] == -1:
				shape.LineTokenCounts[index] = count
			case shape.LineTokenCounts[index] > 0 && shape.LineTokenCounts[index] != count:
				shape.LineTokenCounts[index] = 0
			}
			for len(shape.LineStringLen) <= index {
				shape.LineStringLen = append(shape.LineStringLen, -1)
			}
			if kind == lineKindString {
				lineTokens := strings.Fields(line)
				fixedLength := 0
				if len(lineTokens) > 0 {
					fixedLength = len(lineTokens[0])
					for _, token := range lineTokens[1:] {
						if len(token) != fixedLength {
							fixedLength = 0
							break
						}
					}
				}
				switch {
				case fixedLength == 0:
					shape.LineStringLen[index] = 0
				case shape.LineStringLen[index] == -1:
					shape.LineStringLen[index] = fixedLength
				case shape.LineStringLen[index] > 0 && shape.LineStringLen[index] != fixedLength:
					shape.LineStringLen[index] = 0
				}
			} else {
				shape.LineStringLen[index] = 0
			}
			for _, match := range numberPattern.FindAllString(line, -1) {
				value, err := strconv.ParseInt(match, 10, 64)
				if err != nil {
					continue
				}
				current, exists := shape.LineMax[index]
				if !exists || value > current {
					shape.LineMax[index] = value
				}
			}
		}
		for _, match := range numberPattern.FindAll(seed.Input, -1) {
			value, err := strconv.ParseInt(string(match), 10, 64)
			if err != nil {
				continue
			}
			if value == 0 {
				shape.AllowZero = true
			}
			if value < 0 {
				shape.AllowNegative = true
			}
		}
	}
	if relationKnown && relationSeen {
		switch {
		case relationPlusOne:
			shape.LineRelation = lineRelationPlusOne
		case relationSame:
			shape.LineRelation = lineRelationSame
		}
	}
	return shape, found
}

func shapeAllowsInput(data []byte, shape seedShape) bool {
	tokens := strings.Fields(string(data))
	if len(tokens) < shape.MinTokens || len(tokens) > shape.MaxTokens {
		return false
	}
	lines := lineCount(data)
	if lines < shape.MinLines || lines > shape.MaxLines {
		return false
	}
	if shape.LineRelation != lineRelationNone {
		inputLines := splitInputLines(data)
		if len(inputLines) == 0 {
			return false
		}
		firstFields := strings.Fields(inputLines[0])
		if len(firstFields) == 0 || !numericToken.MatchString(firstFields[0]) {
			return false
		}
		firstValue, err := strconv.ParseInt(firstFields[0], 10, 64)
		if err != nil {
			return false
		}
		expected := firstValue
		if shape.LineRelation == lineRelationPlusOne {
			expected++
		}
		if int64(lines) != expected {
			return false
		}
	}
	for _, token := range tokens {
		if numericToken.MatchString(token) {
			continue
		}
		for i := 0; i < len(token); i++ {
			if !shape.StringChars[token[i]] {
				return false
			}
		}
	}
	for index, line := range splitInputLines(data) {
		kind := lineKindNone
		if index < len(shape.LineKinds) {
			kind = shape.LineKinds[index]
		}
		lineTokens := strings.Fields(line)
		if index < len(shape.LineTokenCounts) {
			fixedCount := shape.LineTokenCounts[index]
			if fixedCount > 0 && len(lineTokens) != fixedCount {
				return false
			}
		}
		maxValue, hasMax := shape.LineMax[index]
		fixedStringLen := 0
		if index < len(shape.LineStringLen) {
			fixedStringLen = shape.LineStringLen[index]
		}
		for _, token := range lineTokens {
			if kind == lineKindNumeric && !numericToken.MatchString(token) {
				return false
			}
			if kind == lineKindString && numericToken.MatchString(token) {
				return false
			}
			if fixedStringLen > 0 && len(token) != fixedStringLen {
				return false
			}
			if !hasMax || maxValue <= 0 || !numericToken.MatchString(token) {
				continue
			}
			value, err := strconv.ParseInt(token, 10, 64)
			if err == nil && value > maxValue*4+16 {
				return false
			}
		}
	}
	for _, match := range numberPattern.FindAll(data, -1) {
		value, err := strconv.ParseInt(string(match), 10, 64)
		if err != nil {
			continue
		}
		if value == 0 && !shape.AllowZero {
			return false
		}
		if value < 0 && !shape.AllowNegative {
			return false
		}
	}
	return true
}

func lineCount(data []byte) int {
	value := strings.TrimRight(string(data), "\r\n")
	if value == "" {
		return 0
	}
	return strings.Count(value, "\n") + 1
}

func splitInputLines(data []byte) []string {
	value := strings.TrimRight(string(data), "\r\n")
	if value == "" {
		return nil
	}
	lines := strings.Split(value, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	return lines
}

func classifyLine(line string) lineKind {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return lineKindNone
	}
	numeric := 0
	textual := 0
	for _, token := range tokens {
		if numericToken.MatchString(token) {
			numeric++
		} else {
			textual++
		}
	}
	switch {
	case textual == 0:
		return lineKindNumeric
	case numeric == 0:
		return lineKindString
	default:
		return lineKindMixed
	}
}

func mergeLineKind(current, next lineKind) lineKind {
	if current == lineKindNone {
		return next
	}
	if current == next {
		return current
	}
	return lineKindMixed
}

func firstTokenValue(data []byte) (int64, bool) {
	lines := splitInputLines(data)
	if len(lines) == 0 {
		return 0, false
	}
	fields := strings.Fields(lines[0])
	if len(fields) == 0 || !numericToken.MatchString(fields[0]) {
		return 0, false
	}
	value, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func inputSignature(data []byte) string {
	value := strings.TrimRight(string(data), "\r\n")
	lines := strings.Split(value, "\n")
	for i := range lines {
		lines[i] = strings.Join(strings.Fields(lines[i]), " ")
	}
	return strings.Join(lines, "\n")
}

func processProblem(ctx context.Context, problem *Problem, opts Options, rng *rand.Rand) Result {
	result := Result{
		Rel:      problem.Rel,
		Existing: len(problem.Seeds),
		Warnings: append([]string(nil), problem.Warnings...),
	}
	if len(problem.Seeds) == 0 {
		result.Status = "no-seed"
		return result
	}
	shape, meaningful := shapeFromSeeds(problem.Seeds)
	if !meaningful {
		result.Status = "weak-seed"
		result.Warnings = append(result.Warnings, "all seed inputs are empty; refusing to mutate them")
		return result
	}

	reach, err := analyzeReach(problem.Solution)
	if err != nil {
		result.Status = "analysis-fail"
		result.Warnings = append(result.Warnings, err.Error())
		return result
	}

	tmpDir, err := os.MkdirTemp("", "dsa-genfixtures-")
	if err != nil {
		result.Status = "temp-fail"
		result.Warnings = append(result.Warnings, err.Error())
		return result
	}
	defer os.RemoveAll(tmpDir)

	binCov := filepath.Join(tmpDir, "solution-cover")
	binPlain := filepath.Join(tmpDir, "solution")
	if err := buildSolution(problem, binCov, true); err != nil {
		result.Status = "build-fail"
		result.BuildErr = err.Error()
		return result
	}
	if err := buildSolution(problem, binPlain, false); err != nil {
		result.Status = "build-fail"
		result.BuildErr = err.Error()
		return result
	}

	infos := make(map[Block]BlockInfo)
	baseline := make(CoverSet)
	corpus := make([][]byte, 0, len(problem.Seeds))
	for i, seed := range problem.Seeds {
		covDir := filepath.Join(tmpDir, fmt.Sprintf("seed-%d", i))
		_ = os.MkdirAll(covDir, 0o755)
		run := runBinary(ctx, binCov, seed.Input, problem.Dir, instrumentedEnv(covDir), runOptions(opts))
		if !runOK(run) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("seed %s failed: %s", filepath.Base(seed.InputPath), failureText(run)))
			continue
		}
		set, profile, err := readCoverage(covDir)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("seed %s: %v", filepath.Base(seed.InputPath), err))
			continue
		}
		mergeInfos(infos, profile)
		baseline = unionCoverage(baseline, set)
		if len(strings.Fields(string(seed.Input))) > 0 {
			corpus = append(corpus, seed.Input)
		}
		if !outputsEquivalent(run.Stdout, seed.Answer) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("existing fixture %s does not match Go output", filepath.Base(seed.InputPath)))
		}
	}

	emptyCov := filepath.Join(tmpDir, "empty")
	_ = os.MkdirAll(emptyCov, 0o755)
	emptyRun := runBinary(ctx, binCov, nil, problem.Dir, instrumentedEnv(emptyCov), runOptions(opts))
	if !emptyRun.TimedOut && !emptyRun.OutputLimited {
		if _, profile, err := readCoverage(emptyCov); err == nil {
			mergeInfos(infos, profile)
		}
	}

	if len(infos) == 0 {
		result.Status = "no-coverage"
		return result
	}
	rawBefore := statsFor(baseline, infos, reach.IsReachable)
	targetBefore := statsFor(baseline, infos, reach.IsTarget)
	result.RawBefore = rawBefore.Percent()
	result.TargetBefore = targetBefore.Percent()
	result.TargetBlocks = targetBefore.Total
	if targetBefore.Total == 0 {
		result.Status = "no-target"
		return result
	}
	if targetBefore.CoveredStmts >= targetBefore.TotalStmts {
		result.Status = "already-covered"
		result.RawAfter = rawBefore.Percent()
		result.TargetAfter = targetBefore.Percent()
		return result
	}

	covered := make(CoverSet)
	addCoverage(covered, baseline)
	seen := make(map[string]bool)
	for _, data := range corpus {
		seen[inputSignature(data)] = true
	}

	var accepted []generatedFixture
	stagnant := 0
	for round := 0; round < opts.Rounds && len(accepted) < opts.MaxNew; round++ {
		candidates := generateCandidates(rng, corpus, opts.Batch, opts.MaxInput, seen, shape)
		if len(candidates) == 0 {
			break
		}
		results := runCandidates(ctx, binCov, problem.Dir, tmpDir, candidates, opts)

		type scored struct {
			index int
			set   CoverSet
			new   int
			run   RunResult
		}
		var scoredResults []scored
		for i := range results {
			if results[i].Set == nil {
				continue
			}
			newSet := newCoverage(results[i].Set, covered)
			targetNew := statsFor(newSet, infos, reach.IsTarget)
			if targetNew.CoveredStmts == 0 {
				continue
			}
			scoredResults = append(scoredResults, scored{
				index: i,
				set:   results[i].Set,
				new:   targetNew.CoveredStmts,
				run:   results[i].Run,
			})
		}
		sort.Slice(scoredResults, func(i, j int) bool {
			if scoredResults[i].new != scoredResults[j].new {
				return scoredResults[i].new > scoredResults[j].new
			}
			return len(results[scoredResults[i].index].Input) < len(results[scoredResults[j].index].Input)
		})

		addedThisRound := 0
		for _, candidate := range scoredResults {
			if len(accepted) >= opts.MaxNew {
				break
			}
			input := results[candidate.index].Input
			remaining := newCoverage(candidate.set, covered)
			remainingTarget := statsFor(remaining, infos, reach.IsTarget)
			if remainingTarget.CoveredStmts == 0 {
				continue
			}
			output, ok, warning := verifyDeterministic(ctx, binPlain, problem.Dir, input, opts)
			if warning != "" {
				result.Warnings = append(result.Warnings, warning)
			}
			if !ok {
				continue
			}
			if len(strings.TrimSpace(string(output))) == 0 {
				if shape.AllAnswersNonEmpty {
					continue
				}
				if first, ok := firstTokenValue(input); ok &&
					shape.NonEmptyAnswerFirst[first] && !shape.EmptyAnswerFirst[first] {
					continue
				}
			}
			if !shape.AllowNegativeAnswer && hasNegativeNumber(output) {
				continue
			}
			if opts.PythonCheck {
				mismatch, failed := pythonValidate(ctx, problem, input, output, opts)
				warning := mismatch
				if failed {
					warning = "Python solution failed on a generated input"
				}
				if warning != "" {
					if opts.PythonVeto {
						warning = "skipped candidate: " + warning
					} else {
						warning += " (accepted; Python is warning-only)"
					}
					result.Warnings = append(result.Warnings, warning)
					if opts.PythonVeto {
						continue
					}
				}
			}
			if hasIgnoredSuffix(ctx, binCov, problem.Dir, tmpDir, input, output, candidate.set, reach, infos, opts) {
				continue
			}
			accepted = append(accepted, generatedFixture{Input: input, Output: output})
			for block := range candidate.set {
				if reach.IsTarget(block) {
					covered[block] = true
				}
			}
			corpus = append(corpus, input)
			seen[inputSignature(input)] = true
			addedThisRound++
		}
		if addedThisRound == 0 {
			stagnant++
		} else {
			stagnant = 0
		}
		if stagnant >= 2 {
			break
		}
	}

	finalRaw := statsFor(covered, infos, reach.IsReachable)
	finalTarget := statsFor(covered, infos, reach.IsTarget)
	result.RawAfter = finalRaw.Percent()
	result.TargetAfter = finalTarget.Percent()
	result.Uncovered = finalTarget.Uncovered
	result.Generated = len(accepted)

	if opts.Write {
		next := problem.NextNumber
		for i := range accepted {
			name, inputPath, answerPath, err := nextFreeName(problem.TestDir, next)
			if err != nil {
				result.Warnings = append(result.Warnings, err.Error())
				break
			}
			next, _ = strconv.Atoi(name)
			next++
			if err := os.WriteFile(inputPath, accepted[i].Input, 0o644); err != nil {
				result.Warnings = append(result.Warnings, err.Error())
				break
			}
			if err := os.WriteFile(answerPath, accepted[i].Output, 0o644); err != nil {
				result.Warnings = append(result.Warnings, err.Error())
				break
			}
			accepted[i].Name = name
			result.GeneratedFiles = append(result.GeneratedFiles, name)
		}
	}

	result.Warnings = dedupeWarnings(result.Warnings)
	if result.TargetAfter >= 99.999 {
		result.Status = "ok"
	} else if result.Generated > 0 {
		result.Status = "partial"
	} else {
		result.Status = "no-new"
	}
	return result
}

func buildSolution(problem *Problem, output string, cover bool) error {
	args := []string{"build"}
	if cover {
		args = append(args, "-cover")
	}
	args = append(args, "-o", output, filepath.Base(problem.Solution))
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goToolPath(), args...)
	cmd.Dir = problem.Dir
	cmd.Env = os.Environ()
	data, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("build timeout")
	}
	if err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(data)))
	}
	return nil
}

func runOptions(opts Options) RunOptions {
	return RunOptions{
		Timeout:     opts.Timeout,
		MaxOutput:   opts.MaxOutput,
		MemoryBytes: opts.MemoryBytes,
	}
}

func runOK(run RunResult) bool {
	return !run.TimedOut && !run.OutputLimited && run.Err == nil
}

func failureText(run RunResult) string {
	if run.TimedOut {
		return "timeout"
	}
	if run.OutputLimited {
		return "output limit"
	}
	if run.Err != nil {
		return run.Err.Error()
	}
	return "unknown"
}

func generateCandidates(rng *rand.Rand, corpus [][]byte, count, maxInput int, seen map[string]bool, shape seedShape) [][]byte {
	if len(corpus) == 0 {
		return nil
	}
	candidates := make([][]byte, 0, count)
	attempts := 0
	for len(candidates) < count && attempts < count*60 {
		attempts++
		base := corpus[rng.Intn(len(corpus))]
		candidate := mutate(rng, base, corpus)
		if len(candidate) == 0 || len(candidate) > maxInput {
			continue
		}
		if !shapeAllowsInput(candidate, shape) {
			continue
		}
		key := inputSignature(candidate)
		if seen[key] {
			continue
		}
		seen[key] = true
		candidates = append(candidates, candidate)
	}
	return candidates
}

func runCandidates(ctx context.Context, binary, workDir, tmpDir string, candidates [][]byte, opts Options) []candidateResult {
	results := make([]candidateResult, len(candidates))
	jobs := opts.Jobs
	if jobs < 1 {
		jobs = 1
	}
	sem := make(chan struct{}, jobs)
	var wg sync.WaitGroup
	for i := range candidates {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			covDir := filepath.Join(tmpDir, fmt.Sprintf("cand-%d", index))
			_ = os.MkdirAll(covDir, 0o755)
			defer os.RemoveAll(covDir)
			run := runBinary(ctx, binary, candidates[index], workDir, instrumentedEnv(covDir), runOptions(opts))
			results[index] = candidateResult{Input: candidates[index], Run: run}
			if !runOK(run) {
				return
			}
			set, _, err := readCoverage(covDir)
			if err != nil {
				return
			}
			results[index].Set = set
		}(i)
	}
	wg.Wait()
	return results
}

func verifyDeterministic(ctx context.Context, binary, workDir string, input []byte, opts Options) ([]byte, bool, string) {
	first := runBinary(ctx, binary, input, workDir, nil, runOptions(opts))
	if !runOK(first) {
		return nil, false, ""
	}
	second := runBinary(ctx, binary, input, workDir, nil, runOptions(opts))
	if !runOK(second) {
		return nil, false, "generated input is not stable on a second run"
	}
	if normalizeOutput(first.Stdout) != normalizeOutput(second.Stdout) {
		return nil, false, "skipped nondeterministic output"
	}
	return first.Stdout, true, ""
}

func nextFreeName(testDir string, start int) (string, string, string, error) {
	for n := start; ; n++ {
		name := strconv.Itoa(n)
		inputPath := filepath.Join(testDir, name)
		answerPath := inputPath + ".a"
		_, errInput := os.Stat(inputPath)
		_, errAnswer := os.Stat(answerPath)
		if os.IsNotExist(errInput) && os.IsNotExist(errAnswer) {
			return name, inputPath, answerPath, nil
		}
	}
}

func hasIgnoredSuffix(ctx context.Context, binary, workDir, tmpDir string, input, output []byte, candidate CoverSet, reach *reachInfo, infos map[Block]BlockInfo, opts Options) bool {
	matches := tokenPattern.FindAllIndex(input, -1)
	if len(matches) < 2 {
		return false
	}
	last := matches[len(matches)-1]
	prefix := input[:last[0]]
	if len(strings.Fields(string(prefix))) == 0 {
		return false
	}

	covDir, err := os.MkdirTemp(tmpDir, "suffix-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(covDir)
	run := runBinary(ctx, binary, prefix, workDir, instrumentedEnv(covDir), runOptions(opts))
	if !runOK(run) {
		return false
	}
	if normalizeOutput(run.Stdout) != normalizeOutput(output) {
		return false
	}
	set, profile, err := readCoverage(covDir)
	if err != nil {
		return false
	}
	mergeInfos(infos, profile)
	remaining := newCoverage(candidate, set)
	remainingTarget := statsFor(remaining, infos, reach.IsTarget)
	return remainingTarget.CoveredStmts == 0
}

func pythonValidate(ctx context.Context, problem *Problem, input, output []byte, opts Options) (string, bool) {
	if problem.Python == "" {
		return "", false
	}
	py, err := exec.LookPath("python3")
	if err != nil {
		return "", false
	}
	runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	command := py
	args := []string{filepath.Base(problem.Python)}
	if prlimitPath != "" && opts.MemoryBytes > 0 {
		command = prlimitPath
		args = []string{"--as=" + itoa64(opts.MemoryBytes), "--core=0", "--cpu=10", "--", py, filepath.Base(problem.Python)}
	}
	cmd := exec.CommandContext(runCtx, command, args...)
	cmd.Dir = problem.Dir
	cmd.Stdin = strings.NewReader(string(input))
	stdout := newLimitedBuffer(opts.MaxOutput)
	cmd.Stdout = stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Run(); err != nil {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", true
	}
	if !outputsEquivalent(stdout.Bytes(), output) {
		return "Python and Go outputs differ (non-unique, invalid, or Python bug)", false
	}
	return "", false
}

func outputsEquivalent(a, b []byte) bool {
	left := strings.Fields(string(a))
	right := strings.Fields(string(b))
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] == right[i] {
			continue
		}
		if !isFloatToken(left[i]) || !isFloatToken(right[i]) {
			return false
		}
		first, errFirst := strconv.ParseFloat(left[i], 64)
		second, errSecond := strconv.ParseFloat(right[i], 64)
		if errFirst != nil || errSecond != nil {
			return false
		}
		scale := math.Max(1, math.Max(math.Abs(first), math.Abs(second)))
		if math.Abs(first-second) > 1e-6*scale {
			return false
		}
	}
	return true
}

func isFloatToken(value string) bool {
	return strings.ContainsAny(value, ".eE")
}

func hasNegativeNumber(data []byte) bool {
	for _, match := range numberPattern.FindAll(data, -1) {
		if len(match) > 0 && match[0] == '-' {
			return true
		}
	}
	return false
}

func dedupeWarnings(warnings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		if warning == "" || seen[warning] {
			continue
		}
		seen[warning] = true
		result = append(result, warning)
		if len(result) >= 12 {
			break
		}
	}
	return result
}
