package main

import (
	"fmt"
	"strings"
	"testing"
)

func requireCheck(t *testing.T, ok bool, detail string, want bool) {
	t.Helper()
	if ok != want {
		t.Fatalf("got %v (%s), want %v", ok, detail, want)
	}
}

func TestCollectingSignaturesChecker(t *testing.T) {
	input := []byte("3\n1 3\n2 5\n3 6\n")
	expected := []byte("1\n3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "1\n3\n", true},
		{"same valid point", "1\n3\n", true},
		{"too many points", "2\n3 6\n", false},
		{"missing segment", "1\n10\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkCollectingSignatures(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestMaximumPrizesChecker(t *testing.T) {
	input := []byte("6\n")
	expected := []byte("3\n1 2 3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "3\n1 2 3\n", true},
		{"another valid split", "3\n1 2 3\n", true},
		{"not distinct", "3\n1 1 4\n", false},
		{"wrong sum", "3\n1 2 4\n", false},
		{"wrong count", "2\n1 5\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkMaximumPrizes(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestPrimitiveCalculatorChecker(t *testing.T) {
	input := []byte("5\n")
	expected := []byte("3\n1 2 4 5\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "3\n1 2 4 5\n", true},
		{"another valid", "3\n1 3 4 5\n", true},
		{"wrong steps", "3\n1 2 5 5\n", false},
		{"wrong length", "2\n1 5\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkPrimitiveCalculator(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestMakeHeapChecker(t *testing.T) {
	input := []byte("5\n5 4 3 2 1\n")
	expected := []byte("3\n1 4\n0 1\n1 3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "3\n1 4\n0 1\n1 3\n", true},
		{"wrong sequence", "2\n0 1\n1 3\n", false},
		{"too many swaps", "2\n1 4\n0 1\n", false},
		{"bad index", "1\n0 9\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkMakeHeap(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestAirlineCrewsChecker(t *testing.T) {
	input := []byte("3 3\n1 1 0\n1 1 0\n0 0 1\n")
	expected := []byte("1 2 3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "1 2 3\n", true},
		{"another valid", "2 1 3\n", true},
		{"duplicate crew", "1 1 3\n", false},
		{"ineligible crew", "1 2 2\n", false},
		{"too few assigned", "1 2 -1\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkAirlineCrews(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestSchoolBusChecker(t *testing.T) {
	input := []byte("4 6\n1 2 20\n1 3 42\n1 4 35\n2 3 30\n2 4 34\n3 4 12\n")
	expected := []byte("97\n1 4 3 2\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "97\n1 4 3 2\n", true},
		{"alternative tour", "97\n1 2 3 4\n", true},
		{"wrong cost", "96\n1 4 3 2\n", false},
		{"missing vertex", "97\n1 4 3 3\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkSchoolBus(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestRescheduleExamsChecker(t *testing.T) {
	input := []byte("4 5\nRRRG\n1 3\n1 4\n3 4\n2 4\n2 3\n")
	expected := []byte("GGBR\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "GGBR\n", true},
		{"same color", "RRGB\n", false},
		{"edge conflict", "GGRB\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkRescheduleExams(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestCircuitDesignChecker(t *testing.T) {
	input := []byte("3 3\n1 -3\n-1 2\n-2 -3\n")
	expected := []byte("SATISFIABLE\n1 2 -3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "SATISFIABLE\n1 2 -3\n", true},
		{"alternative assignment", "SATISFIABLE\n-1 2 -3\n", true},
		{"unsatisfied clause", "SATISFIABLE\n1 2 3\n", false},
		{"missing assignment", "SATISFIABLE\n1 2\n", false},
		{"wrong status", "UNSATISFIABLE\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkCircuitDesign(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestDietChecker(t *testing.T) {
	input := []byte("1 1\n1\n2\n3\n")
	expected := []byte("Bounded solution\n2.000000000000000000\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "Bounded solution\n2.000000000000000000\n", true},
		{"not optimal", "Bounded solution\n1.000000000000000000\n", false},
		{"infeasible", "Bounded solution\n3.000000000000000000\n", false},
		{"wrong status", "Infinity\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkDiet(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestAdAllocationChecker(t *testing.T) {
	input := []byte("1 1\n1\n2\n3\n")
	expected := []byte("Bounded solution\n2.000000000000000000\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "Bounded solution\n2.000000000000000000\n", true},
		{"not optimal", "Bounded solution\n1.000000000000000000\n", false},
		{"wrong status", "Infinity\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkAdAllocation(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestPuzzleAssemblyChecker(t *testing.T) {
	var pieces [5][5]string
	var inputLines []string
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			pieces[i][j] = fmt.Sprintf("(r%d,c%d,r%d,c%d)", i, j, i+1, j+1)
			inputLines = append(inputLines, pieces[i][j])
		}
	}
	input := []byte(strings.Join(inputLines, "\n") + "\n")
	var outLines []string
	for i := 0; i < 5; i++ {
		outLines = append(outLines, strings.Join(pieces[i][:], ";"))
	}
	valid := strings.Join(outLines, "\n") + "\n"
	swapped := append([]string(nil), outLines...)
	swapped[0] = strings.Join([]string{pieces[0][1], pieces[0][0], pieces[0][2], pieces[0][3], pieces[0][4]}, ";")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", valid, true},
		{"wrong adjacency", strings.Join(swapped, "\n") + "\n", false},
		{"missing line", strings.Join(outLines[:4], "\n") + "\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkPuzzleAssembly(input, nil, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestNetworkCirculationChecker(t *testing.T) {
	input := []byte("3 3\n1 2 0 1\n2 3 0 1\n3 1 0 1\n")
	expected := []byte("YES\n1\n1\n1\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "YES\n1\n1\n1\n", true},
		{"flow out of range", "YES\n2\n1\n1\n", false},
		{"not conserved", "YES\n1\n0\n1\n", false},
		{"wrong status", "NO\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkNetworkCirculation(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestOptimalKmerSizeChecker(t *testing.T) {
	input := []byte("AACG\nACGT\nCAAC\nGTTG\nTGCA\n")
	expected := []byte("3\n")
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", "3\n", true},
		{"unbalanced k", "4\n", false},
		{"too large", "5\n", false},
		{"empty graph", "10\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkOptimalKmerSize(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestSuffixTreeFromArrayChecker(t *testing.T) {
	input := []byte("GTAGT$\n5 2 3 0 4 1\n0 0 2 0 1\n")
	expected := []byte("5 6\n2 6\n3 5\n5 6\n2 6\n4 5\n5 6\n2 6\n")
	valid := "GTAGT$\n" + string(expected)
	tests := []struct {
		name   string
		output string
		ok     bool
	}{
		{"valid", valid, true},
		{"valid without text", string(expected), true},
		{"wrong order", "GTAGT$\n2 6\n5 6\n3 5\n5 6\n2 6\n4 5\n5 6\n2 6\n", false},
		{"bad range", "GTAGT$\n5 6\n2 9\n3 5\n5 6\n2 6\n4 5\n5 6\n2 6\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkSuffixTreeFromArray(input, expected, []byte(tt.output))
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestGSMNetworkChecker(t *testing.T) {
	k3 := []byte("3 3\n1 2\n2 3\n1 3\n")
	k4 := []byte("4 6\n1 2\n1 3\n1 4\n2 3\n2 4\n3 4\n")
	sat := []byte("1 1\n1 -1 0\n")
	unsat := []byte("2 1\n1 0\n-1 0\n")
	tests := []struct {
		name   string
		input  []byte
		output []byte
		ok     bool
	}{
		{"K3 sat", k3, sat, true},
		{"K4 unsat", k4, unsat, true},
		{"K4 but sat formula", k4, sat, false},
		{"empty formula", k3, nil, false},
		{"K3 but unsat formula", k3, unsat, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkGSMNetwork(tt.input, nil, tt.output)
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestCleaningApartmentChecker(t *testing.T) {
	path := []byte("5 4\n1 2\n2 3\n3 5\n4 5\n") // 1-2-3-5-4
	star := []byte("4 3\n1 2\n1 3\n1 4\n")      // no Hamiltonian path
	sat := []byte("1 1\n1 -1 0\n")
	unsat := []byte("2 1\n1 0\n-1 0\n")
	tests := []struct {
		name   string
		input  []byte
		output []byte
		ok     bool
	}{
		{"path sat", path, sat, true},
		{"star unsat", star, unsat, true},
		{"star but sat formula", star, sat, false},
		{"empty formula", path, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkCleaningApartment(tt.input, nil, tt.output)
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}

func TestBudgetAllocationChecker(t *testing.T) {
	input := []byte("1 2\n1 1\n1\n")
	good := []byte("1 2\n-1 -2 0\n")
	bad := []byte("2 1\n1 0\n-1 0\n")
	secondVariable := []byte("1 2\n1 -1\n-1\n")
	secondVariableCNF := []byte("2 2\n-1 0\n2 0\n")
	tests := []struct {
		name   string
		input  []byte
		output []byte
		ok     bool
	}{
		{"equivalent", input, good, true},
		{"not equivalent", input, bad, false},
		{"needs second variable", secondVariable, secondVariableCNF, true},
		{"empty formula for data", input, nil, false},
		{"empty legacy", nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, detail := checkBudgetAllocation(tt.input, nil, tt.output)
			requireCheck(t, ok, detail, tt.ok)
		})
	}
}
