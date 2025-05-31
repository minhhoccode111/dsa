package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func checkCollectingSignatures(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 1 {
		return false, "signatures: bad input"
	}
	n := in[0]
	if n < 0 || len(in) < 1+2*n {
		return false, "signatures: truncated input"
	}
	segments := make([][2]int, n)
	for i := 0; i < n; i++ {
		segments[i] = [2]int{in[1+2*i], in[2+2*i]}
	}
	exp, err := intTokens(expected)
	if err != nil || len(exp) < 1 {
		return false, "signatures: bad expected output"
	}
	want := exp[0]
	got, err := intTokens(output)
	if err != nil || len(got) != want+1 || got[0] != want {
		return false, fmt.Sprintf("signatures: got %d points, want %d", maxInt(0, len(got)-1), want)
	}
	points := got[1:]
	seen := make(map[int]bool, len(points))
	for _, point := range points {
		if seen[point] {
			return false, fmt.Sprintf("signatures: duplicate point %d", point)
		}
		seen[point] = true
	}
	for _, segment := range segments {
		covered := false
		for _, point := range points {
			if segment[0] <= point && point <= segment[1] {
				covered = true
				break
			}
		}
		if !covered {
			return false, fmt.Sprintf("signatures: segment [%d,%d] is not covered", segment[0], segment[1])
		}
	}
	return true, ""
}

func checkMaximumPrizes(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) != 1 {
		return false, "prizes: bad input"
	}
	n := in[0]
	exp, err := intTokens(expected)
	if err != nil || len(exp) < 1 {
		return false, "prizes: bad expected output"
	}
	want := exp[0]
	got, err := intTokens(output)
	if err != nil || len(got) != want+1 || got[0] != want {
		return false, fmt.Sprintf("prizes: got %d summands, want %d", maxInt(0, len(got)-1), want)
	}
	seen := make(map[int]bool, want)
	sum := 0
	for _, value := range got[1:] {
		if value <= 0 {
			return false, fmt.Sprintf("prizes: non-positive summand %d", value)
		}
		if seen[value] {
			return false, fmt.Sprintf("prizes: duplicate summand %d", value)
		}
		seen[value] = true
		sum += value
	}
	if sum != n {
		return false, fmt.Sprintf("prizes: sum is %d, want %d", sum, n)
	}
	return true, ""
}

func checkPrimitiveCalculator(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) != 1 {
		return false, "calculator: bad input"
	}
	n := in[0]
	exp, err := intTokens(expected)
	if err != nil || len(exp) < 1 {
		return false, "calculator: bad expected output"
	}
	want := exp[0]
	got, err := intTokens(output)
	if err != nil || len(got) != want+2 || got[0] != want {
		return false, fmt.Sprintf("calculator: got %d operations, want %d", maxInt(0, len(got)-2), want)
	}
	sequence := got[1:]
	if len(sequence) == 0 || sequence[0] != 1 {
		return false, "calculator: sequence must start at 1"
	}
	if sequence[len(sequence)-1] != n {
		return false, fmt.Sprintf("calculator: sequence ends at %d, want %d", sequence[len(sequence)-1], n)
	}
	for i := 1; i < len(sequence); i++ {
		previous := sequence[i-1]
		current := sequence[i]
		if current != previous+1 && current != 2*previous && current != 3*previous {
			return false, fmt.Sprintf("calculator: invalid step %d -> %d", previous, current)
		}
	}
	return true, ""
}

func checkMakeHeap(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 1 {
		return false, "heap: bad input"
	}
	n := in[0]
	if n < 0 || len(in) != n+1 {
		return false, "heap: truncated input"
	}
	values := append([]int(nil), in[1:]...)
	got, err := intTokens(output)
	if err != nil || len(got) < 1 {
		return false, "heap: bad output"
	}
	swaps := got[0]
	if swaps < 0 || swaps > 4*n || len(got) != 1+2*swaps {
		return false, fmt.Sprintf("heap: got %d swaps, want at most %d", swaps, 4*n)
	}
	for i := 0; i < swaps; i++ {
		a, b := got[1+2*i], got[2+2*i]
		if a < 0 || a >= n || b < 0 || b >= n {
			return false, fmt.Sprintf("heap: swap index out of range: %d %d", a, b)
		}
		values[a], values[b] = values[b], values[a]
	}
	for i := 0; i < n; i++ {
		left := 2*i + 1
		right := 2*i + 2
		if left < n && values[i] >= values[left] {
			return false, fmt.Sprintf("heap: parent %d >= left child %d", values[i], values[left])
		}
		if right < n && values[i] >= values[right] {
			return false, fmt.Sprintf("heap: parent %d >= right child %d", values[i], values[right])
		}
	}
	return true, ""
}

func checkAirlineCrews(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 2 {
		return false, "crews: bad input"
	}
	n, m := in[0], in[1]
	if n < 0 || m < 0 || len(in) < 2+n*m {
		return false, "crews: truncated input"
	}
	matrix := in[2 : 2+n*m]
	exp, err := intTokens(expected)
	if err != nil || len(exp) != n {
		return false, "crews: bad expected output"
	}
	wantAssigned := 0
	for _, crew := range exp {
		if crew > 0 {
			wantAssigned++
		}
	}
	got, err := intTokens(output)
	if err != nil || len(got) != n {
		return false, fmt.Sprintf("crews: got %d assignments, want %d", len(got), n)
	}
	used := make(map[int]bool, n)
	assigned := 0
	for flight, crew := range got {
		if crew == -1 {
			continue
		}
		if crew < 1 || crew > m {
			return false, fmt.Sprintf("crews: flight %d has invalid crew %d", flight+1, crew)
		}
		if used[crew] {
			return false, fmt.Sprintf("crews: crew %d assigned twice", crew)
		}
		if matrix[flight*m+crew-1] != 1 {
			return false, fmt.Sprintf("crews: crew %d cannot work flight %d", crew, flight+1)
		}
		used[crew] = true
		assigned++
	}
	if assigned != wantAssigned {
		return false, fmt.Sprintf("crews: assigned %d flights, want %d", assigned, wantAssigned)
	}
	return true, ""
}

func checkSchoolBus(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 2 {
		return false, "bus: bad input"
	}
	n, m := in[0], in[1]
	if n < 1 || m < 0 || len(in) < 2+3*m {
		return false, "bus: truncated input"
	}
	const busInf int64 = 1 << 60
	weight := make([][]int64, n)
	for i := range weight {
		weight[i] = make([]int64, n)
		for j := range weight[i] {
			weight[i][j] = busInf
		}
	}
	for i := 0; i < m; i++ {
		u, v := in[2+3*i]-1, in[3+3*i]-1
		w := int64(in[4+3*i])
		if u < 0 || u >= n || v < 0 || v >= n || w < 0 {
			return false, "bus: edge out of range"
		}
		if w < weight[u][v] {
			weight[u][v] = w
			weight[v][u] = w
		}
	}
	want, err := intTokens(expected)
	if err != nil || len(want) == 0 {
		return false, "bus: bad expected output"
	}
	if want[0] == -1 {
		got, err := intTokens(output)
		if err != nil || len(got) != 1 || got[0] != -1 {
			return false, "bus: expected -1"
		}
		return true, ""
	}
	got, err := intTokens(output)
	if err != nil || len(got) != n+1 {
		return false, fmt.Sprintf("bus: got %d tokens, want cost plus %d vertices", len(got), n)
	}
	cost := int64(got[0])
	path := got[1:]
	if len(path) != n {
		return false, fmt.Sprintf("bus: got path length %d, want %d", len(path), n)
	}
	seen := make([]bool, n)
	var sum int64
	for i, vertex := range path {
		vertex--
		if vertex < 0 || vertex >= n {
			return false, fmt.Sprintf("bus: vertex %d out of range", vertex+1)
		}
		if seen[vertex] {
			return false, fmt.Sprintf("bus: duplicate vertex %d", vertex+1)
		}
		seen[vertex] = true
		next := path[(i+1)%n] - 1
		if weight[vertex][next] == busInf {
			return false, fmt.Sprintf("bus: missing edge %d-%d", vertex+1, next+1)
		}
		sum += weight[vertex][next]
	}
	if sum != cost {
		return false, fmt.Sprintf("bus: path costs %d, output says %d", sum, cost)
	}
	if cost != int64(want[0]) {
		return false, fmt.Sprintf("bus: cost %d, want %d", cost, want[0])
	}
	return true, ""
}

func checkRescheduleExams(input, expected, output []byte) (bool, string) {
	fields := strings.Fields(string(input))
	if len(fields) < 3 {
		return false, "exams: bad input"
	}
	n, errN := strconv.Atoi(fields[0])
	m, errM := strconv.Atoi(fields[1])
	colors := fields[2]
	if errN != nil || errM != nil || n < 1 || m < 0 || len(colors) != n || len(fields) < 3+2*m {
		return false, "exams: bad input"
	}
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		u, errU := strconv.Atoi(fields[3+2*i])
		v, errV := strconv.Atoi(fields[4+2*i])
		if errU != nil || errV != nil || u < 1 || u > n || v < 1 || v > n {
			return false, "exams: bad edge"
		}
		edges[i] = [2]int{u - 1, v - 1}
	}
	wantFields := strings.Fields(string(expected))
	if len(wantFields) == 0 {
		return false, "exams: bad expected output"
	}
	wantImpossible := wantFields[0] == "Impossible"
	gotFields := strings.Fields(string(output))
	if len(gotFields) == 0 {
		return false, "exams: empty output"
	}
	if wantImpossible {
		if len(gotFields) != 1 || gotFields[0] != "Impossible" {
			return false, "exams: expected Impossible"
		}
		return true, ""
	}
	if len(gotFields) != 1 || len(gotFields[0]) != n {
		return false, fmt.Sprintf("exams: got %d color strings, want one of length %d", len(gotFields), n)
	}
	got := gotFields[0]
	for i := 0; i < n; i++ {
		if got[i] != 'R' && got[i] != 'G' && got[i] != 'B' {
			return false, fmt.Sprintf("exams: invalid color %q", got[i])
		}
		if got[i] == colors[i] {
			return false, fmt.Sprintf("exams: vertex %d kept color %q", i+1, got[i])
		}
	}
	for _, edge := range edges {
		if edge[0] == edge[1] || got[edge[0]] == got[edge[1]] {
			return false, fmt.Sprintf("exams: edge %d-%d conflicts", edge[0]+1, edge[1]+1)
		}
	}
	return true, ""
}

func checkCircuitDesign(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 2 {
		return false, "circuit: bad input"
	}
	variables, clausesCount := in[0], in[1]
	if variables < 0 || clausesCount < 0 || len(in) < 2+2*clausesCount {
		return false, "circuit: truncated input"
	}
	clauses := make([][2]int, clausesCount)
	for i := range clauses {
		clauses[i] = [2]int{in[2+2*i], in[3+2*i]}
	}
	wantFields := strings.Fields(string(expected))
	if len(wantFields) == 0 {
		return false, "circuit: bad expected output"
	}
	wantSAT := wantFields[0] == "SATISFIABLE"
	gotFields := strings.Fields(string(output))
	if len(gotFields) == 0 {
		return false, "circuit: empty output"
	}
	if gotFields[0] == "UNSATISFIABLE" {
		if wantSAT {
			return false, "circuit: output UNSATISFIABLE, want SATISFIABLE"
		}
		if len(gotFields) != 1 {
			return false, "circuit: UNSATISFIABLE output has extra tokens"
		}
		return true, ""
	}
	if gotFields[0] != "SATISFIABLE" || !wantSAT {
		return false, "circuit: wrong SAT status"
	}
	if len(gotFields) != variables+1 {
		return false, fmt.Sprintf("circuit: got %d literals, want %d", len(gotFields)-1, variables)
	}
	assignment := make([]bool, variables+1)
	for i := 1; i <= variables; i++ {
		literal, err := strconv.Atoi(gotFields[i])
		if err != nil || literal == 0 || absInt(literal) != i {
			return false, fmt.Sprintf("circuit: bad literal %q for variable %d", gotFields[i], i)
		}
		assignment[i] = literal > 0
	}
	for _, clause := range clauses {
		if !literalValue(clause[0], assignment) && !literalValue(clause[1], assignment) {
			return false, fmt.Sprintf("circuit: clause %d %d is false", clause[0], clause[1])
		}
	}
	return true, ""
}

func checkDiet(input, expected, output []byte) (bool, string) {
	return checkLinearProgram("diet", input, expected, output)
}

func checkAdAllocation(input, expected, output []byte) (bool, string) {
	return checkLinearProgram("ads", input, expected, output)
}

type lpInstance struct {
	constraints [][]float64
	bounds      []float64
	objective   []float64
}

type lpAnswer struct {
	status string
	values []float64
}

func checkLinearProgram(name string, input, expected, output []byte) (bool, string) {
	instance, err := parseLPInstance(input)
	if err != nil {
		return false, fmt.Sprintf("%s: %v", name, err)
	}
	want, err := parseLPAnswer(expected)
	if err != nil {
		return false, fmt.Sprintf("%s: expected: %v", name, err)
	}
	got, err := parseLPAnswer(output)
	if err != nil {
		return false, fmt.Sprintf("%s: %v", name, err)
	}
	if got.status != want.status {
		return false, fmt.Sprintf("%s: status %q, want %q", name, got.status, want.status)
	}
	if want.status != "Bounded solution" {
		return true, ""
	}
	if len(got.values) != len(instance.objective) {
		return false, fmt.Sprintf("%s: got %d values, want %d", name, len(got.values), len(instance.objective))
	}
	if len(want.values) != len(instance.objective) {
		return false, fmt.Sprintf("%s: expected answer has %d values, want %d", name, len(want.values), len(instance.objective))
	}
	for j, value := range got.values {
		if value < -1e-6 {
			return false, fmt.Sprintf("%s: variable %d is negative: %g", name, j+1, value)
		}
	}
	for i, row := range instance.constraints {
		sum := 0.0
		for j, value := range got.values {
			sum += row[j] * value
		}
		if sum > instance.bounds[i]+1e-3 {
			return false, fmt.Sprintf("%s: constraint %d violated: %g > %g", name, i+1, sum, instance.bounds[i])
		}
	}
	gotObjective := dot(instance.objective, got.values)
	wantObjective := dot(instance.objective, want.values)
	if math.Abs(gotObjective-wantObjective) > 1e-6*math.Max(1, math.Abs(wantObjective)) {
		return false, fmt.Sprintf("%s: objective %g, want %g", name, gotObjective, wantObjective)
	}
	return true, ""
}

func parseLPInstance(input []byte) (lpInstance, error) {
	tokens, err := intTokens(input)
	if err != nil || len(tokens) < 2 {
		return lpInstance{}, fmt.Errorf("bad input")
	}
	p, q := tokens[0], tokens[1]
	if p < 0 || q < 0 || len(tokens) != 2+p*q+p+q {
		return lpInstance{}, fmt.Errorf("expected %d constraints and %d variables", p, q)
	}
	index := 2
	constraints := make([][]float64, p)
	for i := 0; i < p; i++ {
		constraints[i] = make([]float64, q)
		for j := 0; j < q; j++ {
			constraints[i][j] = float64(tokens[index])
			index++
		}
	}
	bounds := make([]float64, p)
	for i := 0; i < p; i++ {
		bounds[i] = float64(tokens[index])
		index++
	}
	objective := make([]float64, q)
	for j := 0; j < q; j++ {
		objective[j] = float64(tokens[index])
		index++
	}
	return lpInstance{constraints: constraints, bounds: bounds, objective: objective}, nil
}

func parseLPAnswer(data []byte) (lpAnswer, error) {
	lines := nonEmptyLines(data)
	if len(lines) == 0 {
		return lpAnswer{}, fmt.Errorf("empty output")
	}
	status := strings.TrimSpace(lines[0])
	switch status {
	case "No solution", "Infinity":
		if len(lines) != 1 {
			return lpAnswer{}, fmt.Errorf("unexpected text after %q", status)
		}
		return lpAnswer{status: status}, nil
	case "Bounded solution":
		fields := strings.Fields(strings.Join(lines[1:], " "))
		if len(fields) == 0 {
			return lpAnswer{}, fmt.Errorf("missing solution values")
		}
		values := make([]float64, len(fields))
		for i, field := range fields {
			value, err := strconv.ParseFloat(field, 64)
			if err != nil {
				return lpAnswer{}, fmt.Errorf("bad value %q", field)
			}
			values[i] = value
		}
		return lpAnswer{status: status, values: values}, nil
	default:
		return lpAnswer{}, fmt.Errorf("unknown status %q", status)
	}
}

func dot(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func checkPuzzleAssembly(input, expected, output []byte) (bool, string) {
	inputLines := nonEmptyLines(input)
	side := int(math.Sqrt(float64(len(inputLines))))
	if side < 1 || side*side != len(inputLines) {
		return false, fmt.Sprintf("puzzle: got %d input pieces, want a square", len(inputLines))
	}
	available := make(map[string]int, len(inputLines))
	for _, line := range inputLines {
		piece, ok := parsePuzzlePiece(line)
		if !ok {
			return false, fmt.Sprintf("puzzle: bad input piece %q", strings.TrimSpace(line))
		}
		available[strings.Join(piece[:], ",")]++
	}
	outputLines := nonEmptyLines(output)
	if len(outputLines) != side {
		return false, fmt.Sprintf("puzzle: got %d output rows, want %d", len(outputLines), side)
	}
	grid := make([][4]string, side*side)
	for i, line := range outputLines {
		parts := strings.Split(line, ";")
		if len(parts) != side {
			return false, fmt.Sprintf("puzzle: row %d has %d pieces, want %d", i+1, len(parts), side)
		}
		for j, part := range parts {
			piece, ok := parsePuzzlePiece(part)
			if !ok {
				return false, fmt.Sprintf("puzzle: bad output piece %q", strings.TrimSpace(part))
			}
			key := strings.Join(piece[:], ",")
			if available[key] == 0 {
				return false, fmt.Sprintf("puzzle: piece %q is missing or repeated", key)
			}
			available[key]--
			grid[i*side+j] = piece
		}
	}
	for key, count := range available {
		if count != 0 {
			return false, fmt.Sprintf("puzzle: piece %q left unused", key)
		}
	}
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			piece := grid[i*side+j]
			if j+1 < side && piece[3] != grid[i*side+j+1][1] {
				return false, fmt.Sprintf("puzzle: horizontal mismatch at %d,%d", i+1, j+1)
			}
			if i+1 < side && piece[2] != grid[(i+1)*side+j][0] {
				return false, fmt.Sprintf("puzzle: vertical mismatch at %d,%d", i+1, j+1)
			}
		}
	}
	return true, ""
}

func parsePuzzlePiece(line string) ([4]string, bool) {
	line = strings.TrimSpace(line)
	var piece [4]string
	if len(line) < 2 || line[0] != '(' || line[len(line)-1] != ')' {
		return piece, false
	}
	parts := strings.Split(line[1:len(line)-1], ",")
	if len(parts) != 4 {
		return piece, false
	}
	for i := range parts {
		piece[i] = strings.TrimSpace(parts[i])
		if piece[i] == "" {
			return piece, false
		}
	}
	return piece, true
}

func checkNetworkCirculation(input, expected, output []byte) (bool, string) {
	in, err := intTokens(input)
	if err != nil || len(in) < 2 {
		return false, "circulation: bad input"
	}
	n, m := in[0], in[1]
	if n < 0 || m < 0 || len(in) != 2+4*m {
		return false, "circulation: truncated input"
	}
	type edge struct {
		from, to, lower, capacity int
	}
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		edges[i] = edge{
			from:     in[2+4*i],
			to:       in[3+4*i],
			lower:    in[4+4*i],
			capacity: in[5+4*i],
		}
		if edges[i].from < 1 || edges[i].from > n || edges[i].to < 1 || edges[i].to > n {
			return false, "circulation: edge out of range"
		}
	}
	wantFields := strings.Fields(string(expected))
	if len(wantFields) == 0 {
		return false, "circulation: bad expected output"
	}
	wantFeasible := wantFields[0] == "YES"
	gotFields := strings.Fields(string(output))
	if len(gotFields) == 0 {
		return false, "circulation: empty output"
	}
	if gotFields[0] == "NO" {
		if wantFeasible {
			return false, "circulation: output NO, want YES"
		}
		if len(gotFields) != 1 {
			return false, "circulation: NO output has extra tokens"
		}
		return true, ""
	}
	if gotFields[0] != "YES" || !wantFeasible {
		return false, "circulation: wrong YES/NO status"
	}
	if len(gotFields) != m+1 {
		return false, fmt.Sprintf("circulation: got %d flows, want %d", len(gotFields)-1, m)
	}
	balance := make([]int64, n+1)
	for i, edge := range edges {
		flow, err := strconv.Atoi(gotFields[i+1])
		if err != nil || flow < edge.lower || flow > edge.capacity {
			return false, fmt.Sprintf("circulation: edge %d flow %q outside [%d,%d]", i+1, gotFields[i+1], edge.lower, edge.capacity)
		}
		balance[edge.from] -= int64(flow)
		balance[edge.to] += int64(flow)
	}
	for vertex := 1; vertex <= n; vertex++ {
		if balance[vertex] != 0 {
			return false, fmt.Sprintf("circulation: conservation fails at vertex %d", vertex)
		}
	}
	return true, ""
}

func checkOptimalKmerSize(input, expected, output []byte) (bool, string) {
	reads := make([]string, 0)
	for _, line := range nonEmptyLines(input) {
		reads = append(reads, strings.TrimSpace(line))
	}
	if len(reads) == 0 {
		return false, "kmer: empty input"
	}
	got, err := intTokens(output)
	if err != nil || len(got) != 1 {
		return false, "kmer: output must be one integer"
	}
	k := got[0]
	minLen := len(reads[0])
	for _, read := range reads[1:] {
		if len(read) < minLen {
			minLen = len(read)
		}
	}
	if k < 1 || k >= minLen {
		return false, fmt.Sprintf("kmer: k=%d outside 1..%d", k, minLen-1)
	}
	kmers := make(map[string]bool)
	for _, read := range reads {
		for i := 0; i+k <= len(read); i++ {
			kmers[read[i:i+k]] = true
		}
	}
	if len(kmers) == 0 {
		return false, "kmer: no k-mers"
	}
	next := make(map[string]string, len(kmers))
	inDegree := make(map[string]int, len(kmers))
	for kmer := range kmers {
		prefix, suffix := kmer[:k-1], kmer[1:]
		if _, exists := next[prefix]; exists {
			return false, fmt.Sprintf("kmer: vertex %q has more than one outgoing edge", prefix)
		}
		next[prefix] = suffix
		inDegree[suffix]++
	}
	for vertex := range next {
		if inDegree[vertex] != 1 {
			return false, fmt.Sprintf("kmer: vertex %q has indegree %d", vertex, inDegree[vertex])
		}
	}
	for vertex := range inDegree {
		if _, ok := next[vertex]; !ok {
			return false, fmt.Sprintf("kmer: vertex %q has no outgoing edge", vertex)
		}
	}
	start := ""
	for vertex := range next {
		start = vertex
		break
	}
	seen := make(map[string]bool, len(next))
	current := start
	for {
		if seen[current] {
			break
		}
		seen[current] = true
		nextVertex, ok := next[current]
		if !ok {
			return false, fmt.Sprintf("kmer: broken path at %q", current)
		}
		current = nextVertex
		if current == start {
			break
		}
	}
	if current != start || len(seen) != len(next) {
		return false, "kmer: graph is not a single cycle"
	}
	return true, ""
}

func literalValue(literal int, assignment []bool) bool {
	if literal > 0 {
		return assignment[literal]
	}
	return !assignment[-literal]
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
