// Package sat provides a small DPLL solver used to verify CNF reductions in
// tests.
package sat

// Solve reports whether the CNF formula is satisfiable and, if so, returns one
// satisfying assignment. Variables are numbered from 1 through variables.
func Solve(clauses [][]int, variables int) (bool, []bool) {
	assignment := make([]int, variables+1) // 0 unknown, 1 true, -1 false
	if solve(clauses, assignment, variables) {
		result := make([]bool, variables+1)
		for i := 1; i <= variables; i++ {
			result[i] = assignment[i] == 1
		}
		return true, result
	}
	return false, nil
}

func solve(clauses [][]int, assignment []int, variables int) bool {
	for {
		changed, ok := propagate(clauses, assignment)
		if !ok {
			return false
		}
		if !changed {
			break
		}
	}

	best := -1
	bestLength := 1 << 30
	for _, clause := range clauses {
		satisfied := false
		unassigned := 0
		for _, literal := range clause {
			value := assignment[abs(literal)]
			if value == 0 {
				unassigned++
			} else if (value == 1) == (literal > 0) {
				satisfied = true
				break
			}
		}
		if satisfied {
			continue
		}
		if unassigned == 0 {
			return false
		}
		if unassigned < bestLength {
			bestLength = unassigned
			for _, literal := range clause {
				if assignment[abs(literal)] == 0 {
					best = abs(literal)
					break
				}
			}
		}
	}
	if best == -1 {
		for i := 1; i <= variables; i++ {
			if assignment[i] == 0 {
				assignment[i] = 1
			}
		}
		return true
	}

	for _, value := range []int{1, -1} {
		snapshot := append([]int(nil), assignment...)
		assignment[best] = value
		if solve(clauses, assignment, variables) {
			return true
		}
		copy(assignment, snapshot)
	}
	return false
}

func propagate(clauses [][]int, assignment []int) (bool, bool) {
	changed := false
	for {
		progress := false
		for _, clause := range clauses {
			satisfied := false
			unassigned := 0
			last := 0
			for _, literal := range clause {
				value := assignment[abs(literal)]
				if value == 0 {
					unassigned++
					last = literal
				} else if (value == 1) == (literal > 0) {
					satisfied = true
					break
				}
			}
			if satisfied {
				continue
			}
			if unassigned == 0 {
				return changed, false
			}
			if unassigned == 1 {
				assignment[abs(last)] = 1
				if last < 0 {
					assignment[abs(last)] = -1
				}
				progress = true
				changed = true
			}
		}
		if !progress {
			return changed, true
		}
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
