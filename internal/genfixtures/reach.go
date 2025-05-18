package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
)

type lineRange struct {
	Start int
	End   int
}

type functionInfo struct {
	Name      string
	StartLine int
	EndLine   int
	Decl      *ast.FuncDecl
	Reachable bool
}

type reachInfo struct {
	SolutionBase string
	Functions    []functionInfo
	ErrorRanges  []lineRange
}

func analyzeReach(solution string) (*reachInfo, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, solution, nil, 0)
	if err != nil {
		return nil, err
	}

	info := &reachInfo{SolutionBase: filepath.Base(solution)}
	byName := make(map[string][]*ast.FuncDecl)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		start := fset.Position(fn.Pos()).Line
		end := fset.Position(fn.End()).Line
		entry := functionInfo{
			Name:      fn.Name.Name,
			StartLine: start,
			EndLine:   end,
			Decl:      fn,
		}
		info.Functions = append(info.Functions, entry)
		byName[fn.Name.Name] = append(byName[fn.Name.Name], fn)
	}
	sort.Slice(info.Functions, func(i, j int) bool {
		return info.Functions[i].StartLine < info.Functions[j].StartLine
	})

	visited := make(map[*ast.FuncDecl]bool)
	reachableNames := make(map[string]bool)
	queue := []string{"main"}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		if reachableNames[name] {
			continue
		}
		reachableNames[name] = true
		for _, decl := range byName[name] {
			if visited[decl] {
				continue
			}
			visited[decl] = true
			for callee := range callees(decl) {
				queue = append(queue, callee)
			}
		}
	}
	for i := range info.Functions {
		info.Functions[i].Reachable = reachableNames[info.Functions[i].Name]
	}

	ast.Inspect(file, func(node ast.Node) bool {
		ifStmt, ok := node.(*ast.IfStmt)
		if !ok {
			return true
		}
		if !isInputErrorGuard(ifStmt) {
			return true
		}
		start := fset.Position(ifStmt.Pos()).Line
		end := fset.Position(ifStmt.End()).Line
		info.ErrorRanges = append(info.ErrorRanges, lineRange{Start: start, End: end})
		return true
	})
	return info, nil
}

func callees(decl *ast.FuncDecl) map[string]bool {
	result := make(map[string]bool)
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if name := calledName(call.Fun); name != "" {
			result[name] = true
		}
		return true
	})
	return result
}

func calledName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return value.Sel.Name
	case *ast.ParenExpr:
		return calledName(value.X)
	case *ast.IndexExpr:
		return calledName(value.X)
	case *ast.IndexListExpr:
		return calledName(value.X)
	}
	return ""
}

func isInputErrorGuard(ifStmt *ast.IfStmt) bool {
	if len(ifStmt.Body.List) != 1 {
		return false
	}
	if _, ok := ifStmt.Body.List[0].(*ast.ReturnStmt); !ok {
		return false
	}
	found := false
	ast.Inspect(ifStmt.Cond, func(node ast.Node) bool {
		if ident, ok := node.(*ast.Ident); ok && ident.Name == "err" {
			found = true
		}
		return !found
	})
	return found
}

func (r *reachInfo) funcAt(line int) *functionInfo {
	var best *functionInfo
	for i := range r.Functions {
		fn := &r.Functions[i]
		if line < fn.StartLine || line > fn.EndLine {
			continue
		}
		if best == nil || (fn.EndLine-fn.StartLine) < (best.EndLine-best.StartLine) {
			best = fn
		}
	}
	return best
}

func (r *reachInfo) sameFile(block Block) bool {
	return filepath.Base(block.File) == r.SolutionBase
}

func (r *reachInfo) inErrorRange(line int) bool {
	for _, rr := range r.ErrorRanges {
		if line >= rr.Start && line <= rr.End {
			return true
		}
	}
	return false
}

func (r *reachInfo) IsReachable(block Block) bool {
	if !r.sameFile(block) {
		return false
	}
	fn := r.funcAt(block.StartLine)
	return fn != nil && fn.Reachable
}

func (r *reachInfo) IsTarget(block Block) bool {
	if !r.IsReachable(block) {
		return false
	}
	return !r.inErrorRange(block.StartLine)
}

func (r *reachInfo) FunctionAt(block Block) string {
	if fn := r.funcAt(block.StartLine); fn != nil {
		return fn.Name
	}
	return "?"
}
