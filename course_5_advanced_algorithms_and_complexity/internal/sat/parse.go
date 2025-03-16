package sat

import (
	"bufio"
	"fmt"
	"strings"
)

// ParseCNF parses a DIMACS-style CNF formula. The first line contains the
// number of clauses and variables, followed by one clause per line terminated
// by 0.
func ParseCNF(input string) ([][]int, int) {
	reader := bufio.NewReader(strings.NewReader(input))
	var clausesCount, variables int
	if _, err := fmt.Fscan(reader, &clausesCount, &variables); err != nil {
		return nil, 0
	}
	clauses := make([][]int, clausesCount)
	for i := range clauses {
		for {
			var literal int
			if _, err := fmt.Fscan(reader, &literal); err != nil {
				return clauses, variables
			}
			if literal == 0 {
				break
			}
			clauses[i] = append(clauses[i], literal)
		}
	}
	return clauses, variables
}
