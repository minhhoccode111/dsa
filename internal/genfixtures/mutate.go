package main

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var numberPattern = regexp.MustCompile(`-?[0-9]+`)

func mutate(rng *rand.Rand, input []byte, seeds [][]byte) []byte {
	data := append([]byte(nil), input...)
	steps := 1 + rng.Intn(3)
	for i := 0; i < steps; i++ {
		switch rng.Intn(7) {
		case 0:
			data = mutateNumber(rng, data)
		case 1:
			data = mutateLine(rng, data)
		case 2:
			data = mutateBytes(rng, data)
		case 3:
			data = truncate(rng, data)
		case 4:
			data = splice(rng, data, seeds)
		case 5:
			data = mutateTokens(rng, data)
		case 6:
			data = appendNoise(rng, data)
		}
	}
	return data
}

func mutateNumber(rng *rand.Rand, data []byte) []byte {
	indices := numberPattern.FindAllIndex(data, -1)
	if len(indices) == 0 {
		return data
	}
	index := indices[rng.Intn(len(indices))]
	original := string(data[index[0]:index[1]])
	value, err := strconv.ParseInt(original, 10, 64)
	if err != nil {
		value = 0
	}

	var replacement string
	switch rng.Intn(12) {
	case 0:
		replacement = "0"
	case 1:
		replacement = "1"
	case 2:
		replacement = "2"
	case 3:
		replacement = "3"
	case 4:
		replacement = strconv.FormatInt(value+1, 10)
	case 5:
		replacement = strconv.FormatInt(value-1, 10)
	case 6:
		replacement = strconv.FormatInt(value*2, 10)
	case 7:
		if value > 1 {
			replacement = strconv.FormatInt(value/2, 10)
		} else {
			replacement = "10"
		}
	case 8:
		replacement = strconv.Itoa(rng.Intn(10))
	case 9:
		replacement = strconv.Itoa(rng.Intn(1000))
	case 10:
		if value > 0 {
			replacement = strconv.FormatInt(rng.Int63n(value+2), 10)
		} else {
			replacement = strconv.Itoa(rng.Intn(100))
		}
	default:
		boundaries := []int64{10, 100, 1000, 10000, 100000}
		replacement = strconv.FormatInt(boundaries[rng.Intn(len(boundaries))], 10)
	}

	result := make([]byte, 0, len(data)-len(original)+len(replacement))
	result = append(result, data[:index[0]]...)
	result = append(result, replacement...)
	result = append(result, data[index[1]:]...)
	return result
}

func mutateLine(rng *rand.Rand, data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return data
	}
	switch rng.Intn(4) {
	case 0:
		at := rng.Intn(len(lines))
		lines = append(lines[:at+1], append([]string{lines[at]}, lines[at+1:]...)...)
	case 1:
		if len(lines) > 1 {
			at := rng.Intn(len(lines))
			lines = append(lines[:at], lines[at+1:]...)
		}
	case 2:
		if len(lines) > 1 {
			a := rng.Intn(len(lines))
			b := rng.Intn(len(lines))
			lines[a], lines[b] = lines[b], lines[a]
		}
	case 3:
		at := rng.Intn(len(lines) + 1)
		line := strconv.Itoa(rng.Intn(100))
		lines = append(lines[:at], append([]string{line}, lines[at:]...)...)
	}
	return []byte(strings.Join(lines, "\n"))
}

func mutateTokens(rng *rand.Rand, data []byte) []byte {
	tokens := strings.Fields(string(data))
	if len(tokens) == 0 {
		return data
	}
	switch rng.Intn(3) {
	case 0:
		if len(tokens) > 1 {
			at := rng.Intn(len(tokens))
			tokens = append(tokens[:at], tokens[at+1:]...)
		}
	case 1:
		at := rng.Intn(len(tokens))
		tokens = append(tokens[:at+1], append([]string{tokens[at]}, tokens[at+1:]...)...)
	case 2:
		tokens = append(tokens, strconv.Itoa(rng.Intn(1000)))
	}
	return []byte(strings.Join(tokens, " ") + "\n")
}

func mutateBytes(rng *rand.Rand, data []byte) []byte {
	alphabet := []byte("0123456789 \nABCdef$-")
	if len(data) == 0 {
		return []byte{alphabet[rng.Intn(len(alphabet))]}
	}
	at := rng.Intn(len(data))
	switch rng.Intn(3) {
	case 0:
		data = append([]byte(nil), data...)
		data[at] = alphabet[rng.Intn(len(alphabet))]
	case 1:
		value := alphabet[rng.Intn(len(alphabet))]
		data = append(data[:at], append([]byte{value}, data[at:]...)...)
	case 2:
		if len(data) > 1 {
			data = append(data[:at], data[at+1:]...)
		}
	}
	return data
}

func truncate(rng *rand.Rand, data []byte) []byte {
	if len(data) <= 1 {
		return data
	}
	at := 1 + rng.Intn(len(data)-1)
	if newline := strings.LastIndexByte(string(data[:at]), '\n'); newline > 0 && rng.Intn(2) == 0 {
		at = newline + 1
	}
	return append([]byte(nil), data[:at]...)
}

func splice(rng *rand.Rand, data []byte, seeds [][]byte) []byte {
	if len(seeds) == 0 {
		return data
	}
	other := seeds[rng.Intn(len(seeds))]
	if len(other) == 0 {
		return data
	}
	linesA := strings.Split(string(data), "\n")
	linesB := strings.Split(string(other), "\n")
	cutA := rng.Intn(len(linesA) + 1)
	cutB := rng.Intn(len(linesB) + 1)
	result := append([]string(nil), linesA[:cutA]...)
	result = append(result, linesB[cutB:]...)
	return []byte(strings.Join(result, "\n"))
}

func appendNoise(rng *rand.Rand, data []byte) []byte {
	result := append([]byte(nil), data...)
	if len(result) > 0 && result[len(result)-1] != '\n' {
		result = append(result, '\n')
	}
	values := []string{"0", "1", "2", "10", "100", "A", "AA", "$"}
	return append(result, []byte(values[rng.Intn(len(values))]+"\n")...)
}
