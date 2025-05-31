package main

import (
	"fmt"
)

type cnfFormula struct {
	vars     int
	declared int
	clauses  [][]int
}

func parseCNF(data []byte) (cnfFormula, error) {
	tokens, err := intTokens(data)
	if err != nil {
		return cnfFormula{}, err
	}
	if len(tokens) == 0 {
		return cnfFormula{}, nil
	}
	if len(tokens) < 2 {
		return cnfFormula{}, fmt.Errorf("missing clause/variable counts")
	}
	clausesCount, variables := tokens[0], tokens[1]
	if clausesCount < 0 || variables < 0 {
		return cnfFormula{}, fmt.Errorf("negative clause or variable count")
	}
	index := 2
	clauses := make([][]int, 0, clausesCount)
	for i := 0; i < clausesCount; i++ {
		clause := make([]int, 0, 3)
		for index < len(tokens) && tokens[index] != 0 {
			literal := tokens[index]
			if literal == 0 || absInt(literal) > variables {
				return cnfFormula{}, fmt.Errorf("literal %d out of range", literal)
			}
			clause = append(clause, literal)
			index++
		}
		if index >= len(tokens) || tokens[index] != 0 {
			return cnfFormula{}, fmt.Errorf("clause %d is not terminated by 0", i+1)
		}
		index++
		clauses = append(clauses, clause)
	}
	if index != len(tokens) {
		return cnfFormula{}, fmt.Errorf("trailing tokens after %d clauses", clausesCount)
	}
	return cnfFormula{vars: variables, declared: clausesCount, clauses: simplifyCNF(clauses)}, nil
}

func simplifyCNF(clauses [][]int) [][]int {
	result := make([][]int, 0, len(clauses))
	for _, clause := range clauses {
		seen := make(map[int]bool, len(clause))
		tautology := false
		unique := make([]int, 0, len(clause))
		for _, literal := range clause {
			if seen[-literal] {
				tautology = true
				break
			}
			if seen[literal] {
				continue
			}
			seen[literal] = true
			unique = append(unique, literal)
		}
		if !tautology {
			result = append(result, unique)
		}
	}
	return result
}

func cnfSatisfiable(formula cnfFormula) bool {
	assignment := make([]int8, formula.vars+1)
	return dpll(formula, assignment)
}

func dpll(formula cnfFormula, assignment []int8) bool {
	for {
		changed := false
		for _, clause := range formula.clauses {
			satisfied, unassigned, unit := clauseState(clause, assignment)
			if satisfied {
				continue
			}
			if unassigned == 0 {
				return false
			}
			if unassigned == 1 {
				assignment[absInt(unit)] = literalSign(unit)
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	var chosen int
	for _, clause := range formula.clauses {
		satisfied, unassigned, _ := clauseState(clause, assignment)
		if satisfied || unassigned == 0 {
			continue
		}
		for _, literal := range clause {
			if assignment[absInt(literal)] == 0 {
				chosen = absInt(literal)
				break
			}
		}
		if chosen != 0 {
			break
		}
	}
	if chosen == 0 {
		return true
	}
	for _, value := range []int8{1, -1} {
		next := append([]int8(nil), assignment...)
		next[chosen] = value
		if dpll(formula, next) {
			return true
		}
	}
	return false
}

func clauseState(clause []int, assignment []int8) (bool, int, int) {
	unassigned := 0
	unit := 0
	for _, literal := range clause {
		value := assignment[absInt(literal)]
		if value == 0 {
			unassigned++
			unit = literal
			continue
		}
		if (literal > 0 && value == 1) || (literal < 0 && value == -1) {
			return true, 0, 0
		}
	}
	return false, unassigned, unit
}

func literalSign(literal int) int8 {
	if literal > 0 {
		return 1
	}
	return -1
}

type graph struct {
	n     int
	edges [][2]int
	adj   [][]int
}

func parseUndirectedGraph(data []byte) (graph, error) {
	tokens, err := intTokens(data)
	if err != nil || len(tokens) < 2 {
		return graph{}, fmt.Errorf("bad graph input")
	}
	n, m := tokens[0], tokens[1]
	if n < 0 || m < 0 || len(tokens) != 2+2*m {
		return graph{}, fmt.Errorf("expected %d edges", m)
	}
	result := graph{n: n, edges: make([][2]int, m), adj: make([][]int, n)}
	for i := 0; i < m; i++ {
		u, v := tokens[2+2*i]-1, tokens[3+2*i]-1
		if u < 0 || u >= n || v < 0 || v >= n {
			return graph{}, fmt.Errorf("edge %d out of range", i+1)
		}
		result.edges[i] = [2]int{u, v}
		result.adj[u] = append(result.adj[u], v)
		result.adj[v] = append(result.adj[v], u)
	}
	return result, nil
}

func graph3Colorable(g graph) bool {
	for _, edge := range g.edges {
		if edge[0] == edge[1] {
			return false
		}
	}
	colors := make([]int8, g.n)
	var search func(colored int) bool
	search = func(colored int) bool {
		if colored == g.n {
			return true
		}
		best := -1
		bestScore := -1
		for vertex := 0; vertex < g.n; vertex++ {
			if colors[vertex] != 0 {
				continue
			}
			score := 0
			for _, neighbor := range g.adj[vertex] {
				if colors[neighbor] != 0 {
					score++
				}
			}
			if score > bestScore {
				bestScore = score
				best = vertex
			}
		}
		used := [4]bool{}
		for _, neighbor := range g.adj[best] {
			if colors[neighbor] != 0 {
				used[colors[neighbor]] = true
			}
		}
		for color := int8(1); color <= 3; color++ {
			if used[color] {
				continue
			}
			colors[best] = color
			if search(colored + 1) {
				return true
			}
			colors[best] = 0
		}
		return false
	}
	return search(0)
}

func checkGSMNetwork(input, expected, output []byte) (bool, string) {
	g, err := parseUndirectedGraph(input)
	if err != nil {
		return false, fmt.Sprintf("gsm: %v", err)
	}
	formula, err := parseCNF(output)
	if err != nil {
		return false, fmt.Sprintf("gsm: %v", err)
	}
	if formula.vars < 1 || formula.vars > 3000 || formula.declared < 1 || formula.declared > 5000 {
		return false, fmt.Sprintf("gsm: formula has %d variables and %d clauses", formula.vars, formula.declared)
	}
	colorable := graph3Colorable(g)
	satisfiable := cnfSatisfiable(formula)
	if colorable != satisfiable {
		return false, fmt.Sprintf("gsm: graph 3-colorable=%v, formula satisfiable=%v", colorable, satisfiable)
	}
	return true, ""
}

func checkCleaningApartment(input, expected, output []byte) (bool, string) {
	g, err := parseUndirectedGraph(input)
	if err != nil {
		return false, fmt.Sprintf("apartment: %v", err)
	}
	formula, err := parseCNF(output)
	if err != nil {
		return false, fmt.Sprintf("apartment: %v", err)
	}
	if formula.vars < 1 || formula.vars > 3000 || formula.declared < 1 || formula.declared > 5000 {
		return false, fmt.Sprintf("apartment: formula has %d variables and %d clauses", formula.vars, formula.declared)
	}
	satisfiable := cnfSatisfiable(formula)
	if g.n <= 20 {
		path := hasHamiltonianPath(g)
		if path != satisfiable {
			return false, fmt.Sprintf("apartment: Hamiltonian path=%v, formula satisfiable=%v", path, satisfiable)
		}
		return true, ""
	}
	reference, err := parseCNF(expected)
	if err != nil {
		return false, fmt.Sprintf("apartment: expected: %v", err)
	}
	if cnfSatisfiable(reference) != satisfiable {
		return false, "apartment: formula status differs from reference"
	}
	return true, ""
}

func hasHamiltonianPath(g graph) bool {
	n := g.n
	if n == 0 {
		return true
	}
	if n == 1 {
		return true
	}
	if n > 20 {
		return false
	}
	adj := make([]uint32, n)
	for u := 0; u < n; u++ {
		for _, v := range g.adj[u] {
			if u != v {
				adj[u] |= 1 << uint(v)
			}
		}
	}
	full := uint32(1)<<uint(n) - 1
	dp := make([]bool, (1<<uint(n))*n)
	for v := 0; v < n; v++ {
		dp[(1<<uint(v))*n+v] = true
	}
	for mask := 0; mask < 1<<uint(n); mask++ {
		for last := 0; last < n; last++ {
			if !dp[mask*n+last] {
				continue
			}
			if uint32(mask) == full {
				return true
			}
			candidates := adj[last] &^ uint32(mask)
			for candidates != 0 {
				bit := candidates & -candidates
				next := trailingZeros(bit)
				dp[(mask|int(bit))*n+next] = true
				candidates ^= bit
			}
		}
	}
	return false
}

func trailingZeros(value uint32) int {
	count := 0
	for value&1 == 0 {
		value >>= 1
		count++
	}
	return count
}

func checkBudgetAllocation(input, expected, output []byte) (bool, string) {
	formula, err := parseCNF(output)
	if err != nil {
		return false, fmt.Sprintf("budget: %v", err)
	}
	if formula.vars > 3000 || formula.declared > 5000 {
		return false, fmt.Sprintf("budget: formula has %d variables and %d clauses", formula.vars, formula.declared)
	}
	tokens, err := intTokens(input)
	if err != nil {
		return false, fmt.Sprintf("budget: %v", err)
	}
	if len(tokens) == 0 {
		return formula.declared == 0 && formula.vars == 0, "budget: empty legacy fixture"
	}
	if formula.vars < 1 || formula.declared < 1 {
		return false, "budget: empty formula for non-empty input"
	}
	if len(tokens) < 2 {
		return false, "budget: bad input"
	}
	n, m := tokens[0], tokens[1]
	if n < 0 || m < 0 || len(tokens) != 2+n*m+n {
		return false, "budget: truncated input"
	}
	index := 2
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			matrix[i][j] = tokens[index]
			index++
		}
	}
	bounds := make([]int, n)
	for i := 0; i < n; i++ {
		bounds[i] = tokens[index]
		index++
	}
	satisfiable := cnfSatisfiable(formula)
	if m <= 20 {
		feasible := budgetFeasible(n, m, matrix, bounds)
		if feasible != satisfiable {
			return false, fmt.Sprintf("budget: constraints feasible=%v, formula satisfiable=%v", feasible, satisfiable)
		}
		return true, ""
	}
	reference, err := parseCNF(expected)
	if err != nil {
		return false, fmt.Sprintf("budget: expected: %v", err)
	}
	if cnfSatisfiable(reference) != satisfiable {
		return false, "budget: formula status differs from reference"
	}
	return true, ""
}

func budgetFeasible(n, m int, matrix [][]int, bounds []int) bool {
	if m > 20 {
		return false
	}
	for mask := 0; mask < 1<<uint(m); mask++ {
		ok := true
		for i := 0; i < n; i++ {
			sum := 0
			for j := 0; j < m; j++ {
				if mask&(1<<uint(j)) != 0 {
					sum += matrix[i][j]
				}
			}
			if sum > bounds[i] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
