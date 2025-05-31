package main

import (
	"fmt"
	"strconv"
	"strings"
)

type semanticChecker func(input, expected, output []byte) (bool, string)

var problemCheckers = map[string]semanticChecker{
	"course_1_algorithmic_toolbox/week_3_greedy_algorithms/5_collecting_signatures":                                checkCollectingSignatures,
	"course_1_algorithmic_toolbox/week_3_greedy_algorithms/6_maximum_number_of_prizes":                             checkMaximumPrizes,
	"course_1_algorithmic_toolbox/week_5_dynamic_programming_1/2_primitive_calculator":                             checkPrimitiveCalculator,
	"course_2_data_structures/week_2_priority_queues_and_disjoint_sets/1_make_heap":                                checkMakeHeap,
	"course_2_data_structures/week_3_hash_tables/5_longest_common_substring":                                       checkLongestCommonSubstring,
	"course_3_algorithms_on_graphs/week_2_graph_decomposition2/2_order_of_courses":                                 checkTopologicalOrder,
	"course_4_algorithms_on_strings/week_1_suffix_trees/5_non_shared_substring":                                    checkNonSharedSubstring,
	"course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/4_suffix_tree_from_array":   checkSuffixTreeFromArray,
	"course_5_advanced_algorithms_and_complexity/week_1_flows_in_networks/2_airline_crews":                         checkAirlineCrews,
	"course_5_advanced_algorithms_and_complexity/week_2_linear_programming/2_diet":                                 checkDiet,
	"course_5_advanced_algorithms_and_complexity/week_2_linear_programming/3_ad_allocation":                        checkAdAllocation,
	"course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/1_gsm_network":                        checkGSMNetwork,
	"course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/2_cleaning_apartment":                 checkCleaningApartment,
	"course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/3_budget_allocation":                  checkBudgetAllocation,
	"course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/1_circuit_design":              checkCircuitDesign,
	"course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/3_school_bus":                  checkSchoolBus,
	"course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/4_reschedule_exams":            checkRescheduleExams,
	"course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/1_puzzle_assembly":             checkPuzzleAssembly,
	"course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/3_k_universal_circular_string": checkKUniversal,
	"course_6_genome_assembly/week_3_faces_real_sequencing_data/1_network_circulation":                             checkNetworkCirculation,
	"course_6_genome_assembly/week_3_faces_real_sequencing_data/2_optimal_kmer_size":                               checkOptimalKmerSize,
}

func checkerFor(rel string) semanticChecker {
	return problemCheckers[rel]
}

func checkTopologicalOrder(input, expected, output []byte) (bool, string) {
	tokens, err := intTokens(input)
	if err != nil || len(tokens) < 2 {
		return false, fmt.Sprintf("toposort: bad input: %v", err)
	}
	n, m := tokens[0], tokens[1]
	if n < 0 || m < 0 || len(tokens) < 2+2*m {
		return false, "toposort: truncated input"
	}
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		edges[i] = [2]int{tokens[2+2*i], tokens[2+2*i+1]}
	}
	order, err := intTokens(output)
	if err != nil {
		return false, "toposort: output is not an integer sequence"
	}
	if n == 0 {
		if len(order) == 0 {
			return true, ""
		}
		return false, "toposort: output should be empty for n=0"
	}
	if len(order) != n {
		return false, fmt.Sprintf("toposort: got %d positions, want %d", len(order), n)
	}
	position := make(map[int]int, n)
	for i, vertex := range order {
		if vertex < 1 || vertex > n {
			return false, fmt.Sprintf("toposort: vertex %d out of range 1..%d", vertex, n)
		}
		if _, exists := position[vertex]; exists {
			return false, fmt.Sprintf("toposort: duplicate vertex %d", vertex)
		}
		position[vertex] = i
	}
	for _, edge := range edges {
		from, to := edge[0], edge[1]
		if from < 1 || from > n || to < 1 || to > n {
			return false, "toposort: edge out of range"
		}
		if position[from] >= position[to] {
			return false, fmt.Sprintf("toposort: edge %d -> %d is reversed", from, to)
		}
	}
	return true, ""
}

func checkKUniversal(input, expected, output []byte) (bool, string) {
	tokens := strings.Fields(string(input))
	if len(tokens) == 0 {
		return false, "k-universal: empty input"
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil || k < 1 || k > 20 {
		return false, fmt.Sprintf("k-universal: bad k %q", tokens[0])
	}
	value := strings.TrimSpace(string(output))
	wantLength := 1 << k
	if len(value) != wantLength {
		return false, fmt.Sprintf("k-universal: got length %d, want %d", len(value), wantLength)
	}
	for i := 0; i < len(value); i++ {
		if value[i] != '0' && value[i] != '1' {
			return false, fmt.Sprintf("k-universal: invalid character %q", value[i])
		}
	}
	seen := make(map[string]bool, wantLength)
	for i := 0; i < len(value); i++ {
		var builder strings.Builder
		builder.Grow(k)
		for j := 0; j < k; j++ {
			builder.WriteByte(value[(i+j)%len(value)])
		}
		seen[builder.String()] = true
	}
	if len(seen) != wantLength {
		return false, fmt.Sprintf("k-universal: saw %d/%d distinct k-mers", len(seen), wantLength)
	}
	return true, ""
}

func checkNonSharedSubstring(input, expected, output []byte) (bool, string) {
	fields := strings.Fields(string(input))
	if len(fields) < 2 {
		return false, "non-shared: input needs two strings"
	}
	p, q := fields[0], fields[1]
	answer := strings.TrimSpace(string(output))
	wantLength := len(strings.TrimSpace(string(expected)))
	if wantLength == 0 {
		return false, "non-shared: empty expected answer"
	}
	if len(answer) != wantLength {
		return false, fmt.Sprintf("non-shared: got length %d, want %d", len(answer), wantLength)
	}
	if !strings.Contains(p, answer) {
		return false, "non-shared: answer does not occur in the first string"
	}
	if strings.Contains(q, answer) {
		return false, "non-shared: answer also occurs in the second string"
	}
	return true, ""
}

func checkLongestCommonSubstring(input, expected, output []byte) (bool, string) {
	inputLines := nonEmptyLines(input)
	expectedLines := nonEmptyLines(expected)
	outputLines := nonEmptyLines(output)
	if len(inputLines) != len(expectedLines) || len(outputLines) != len(expectedLines) {
		return false, fmt.Sprintf("lcs: got %d output lines, want %d", len(outputLines), len(expectedLines))
	}
	for i := range inputLines {
		fields := strings.Fields(inputLines[i])
		if len(fields) < 2 {
			return false, fmt.Sprintf("lcs: input line %d needs two strings", i+1)
		}
		expectedFields := strings.Fields(expectedLines[i])
		if len(expectedFields) < 3 {
			return false, fmt.Sprintf("lcs: expected line %d needs three integers", i+1)
		}
		wantLength, err := strconv.Atoi(expectedFields[2])
		if err != nil || wantLength < 0 {
			return false, fmt.Sprintf("lcs: bad expected length on line %d", i+1)
		}
		gotFields := strings.Fields(outputLines[i])
		if len(gotFields) != 3 {
			return false, fmt.Sprintf("lcs: output line %d needs i j length", i+1)
		}
		startA, errA := strconv.Atoi(gotFields[0])
		startB, errB := strconv.Atoi(gotFields[1])
		length, errLength := strconv.Atoi(gotFields[2])
		if errA != nil || errB != nil || errLength != nil {
			return false, fmt.Sprintf("lcs: non-integer output on line %d", i+1)
		}
		a, b := fields[0], fields[1]
		if length != wantLength {
			return false, fmt.Sprintf("lcs: line %d length %d, want %d", i+1, length, wantLength)
		}
		if startA < 0 || startB < 0 || startA+length > len(a) || startB+length > len(b) {
			return false, fmt.Sprintf("lcs: line %d indices out of range", i+1)
		}
		if length > 0 && a[startA:startA+length] != b[startB:startB+length] {
			return false, fmt.Sprintf("lcs: line %d substrings differ", i+1)
		}
	}
	return true, ""
}

func intTokens(data []byte) ([]int, error) {
	fields := strings.Fields(string(data))
	result := make([]int, len(fields))
	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

func nonEmptyLines(data []byte) []string {
	raw := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	result := make([]string, 0, len(raw))
	for _, line := range raw {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
