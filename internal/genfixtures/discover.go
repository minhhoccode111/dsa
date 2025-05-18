package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Seed struct {
	InputPath  string
	AnswerPath string
	Input      []byte
	Answer     []byte
}

type Problem struct {
	Rel        string
	Dir        string
	Solution   string
	Python     string
	TestDir    string
	Seeds      []Seed
	NextNumber int
	Warnings   []string
}

var numericName = regexp.MustCompile(`^[0-9]+$`)

func discover(root string) ([]*Problem, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var problems []*Problem
	err = filepath.WalkDir(absRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		if entry.Name() != "test" {
			return nil
		}
		if path == absRoot {
			return nil
		}

		dir := filepath.Dir(path)
		problem, ok, err := discoverProblem(absRoot, dir, path)
		if err != nil {
			return err
		}
		if ok {
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

func discoverProblem(absRoot, dir, testDir string) (*Problem, bool, error) {
	goFiles, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, false, err
	}
	var solutions []string
	for _, file := range goFiles {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		solutions = append(solutions, file)
	}
	if len(solutions) == 0 {
		return nil, false, nil
	}
	sort.Strings(solutions)
	solution := solutions[0]

	pyFiles, _ := filepath.Glob(filepath.Join(dir, "*.py"))
	var python string
	for _, file := range pyFiles {
		if filepath.Base(file) == "starter.py" {
			continue
		}
		python = file
		break
	}

	rel, err := filepath.Rel(absRoot, dir)
	if err != nil {
		return nil, false, err
	}
	problem := &Problem{
		Rel:        rel,
		Dir:        dir,
		Solution:   solution,
		Python:     python,
		TestDir:    testDir,
		NextNumber: 1,
	}
	if len(solutions) > 1 {
		problem.Warnings = append(problem.Warnings,
			fmt.Sprintf("%d non-test Go files; using %s", len(solutions), filepath.Base(solution)))
	}

	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, false, err
	}
	maxNumber := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if numericName.MatchString(name) {
			if n, err := strconv.Atoi(name); err == nil && n > maxNumber {
				maxNumber = n
			}
		}
		if strings.HasSuffix(name, ".a") || strings.HasPrefix(name, ".") {
			continue
		}
		inputPath := filepath.Join(testDir, name)
		answerPath := inputPath + ".a"
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if !info.Mode().IsRegular() || info.Size() > 1<<20 {
			continue
		}
		if _, err := os.Stat(answerPath); err != nil {
			continue
		}
		input, err := os.ReadFile(inputPath)
		if err != nil {
			continue
		}
		answer, err := os.ReadFile(answerPath)
		if err != nil {
			continue
		}
		problem.Seeds = append(problem.Seeds, Seed{
			InputPath:  inputPath,
			AnswerPath: answerPath,
			Input:      input,
			Answer:     answer,
		})
	}
	sort.Slice(problem.Seeds, func(i, j int) bool {
		return problem.Seeds[i].InputPath < problem.Seeds[j].InputPath
	})
	problem.NextNumber = maxNumber + 1
	return problem, true, nil
}
