package dsa_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Coursera compiles a single submitted Go file without this module. Every
// non-test solution must therefore avoid dsa/... imports and must not depend on
// helper declarations from another non-test file in the same directory.
func TestGoSolutionsAreSelfContained(t *testing.T) {
	fset := token.NewFileSet()
	solutions := make(map[string][]string)

	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// internal/ holds repository tooling, not Coursera submissions.
		internalDir := "internal" + string(filepath.Separator)
		if strings.HasPrefix(path, internalDir) ||
			strings.Contains(path, string(filepath.Separator)+internalDir) {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			if strings.HasPrefix(importPath, "dsa/") {
				t.Errorf("%s imports %s; a submitted file cannot use this module", path, importPath)
			}
		}

		dir := filepath.Dir(path)
		solutions[dir] = append(solutions[dir], path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for dir, files := range solutions {
		if len(files) > 1 {
			t.Errorf(
				"%s has %d non-test Go files; a Coursera submission must be one file: %v",
				dir,
				len(files),
				files,
			)
		}
	}
}
