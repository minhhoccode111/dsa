// go run error_prone.go

// Assemble a circular genome from error-prone reads using overlap graphs.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dsa/course_6_genome_assembly/internal/assembly"
)

func hammingWithin(first, second string, errors int) bool {
	if len(first) != len(second) {
		return false
	}
	mismatches := 0
	for i := range first {
		if first[i] != second[i] {
			mismatches++
			if mismatches > errors {
				return false
			}
		}
	}
	return true
}

func removeDuplicates(reads []string, readLength, errors int) []string {
	skip := make([]bool, len(reads))
	for i := 0; i < len(reads)-1; i++ {
		if skip[i] {
			continue
		}
		for j := i + 1; j < len(reads); j++ {
			if skip[j] {
				continue
			}
			if hammingWithin(reads[i], reads[j], errors) {
				skip[j] = true
			}
		}
	}
	result := make([]string, 0, len(reads))
	for i, read := range reads {
		if !skip[i] {
			result = append(result, read)
		}
	}
	return result
}

func findOverlap(first, second string, previous, errors int) (int, bool) {
	for size := len(first) - 1; size > previous; size-- {
		mismatches := 0
		start := len(first) - size
		valid := true
		for i := 0; i < size; i++ {
			if first[start+i] != second[i] {
				mismatches++
				if mismatches > errors {
					valid = false
					break
				}
			}
		}
		if valid {
			return size, true
		}
	}
	return 0, false
}

func assembleGreedy(data []string, readLength, errors int) string {
	size := len(data)
	processed := make([]bool, size)
	overlap := make([]int, size)
	nextID := make([]int, size)
	for i := range nextID {
		nextID[i] = -1
	}

	current := 0
	for step := 0; step < size-1; step++ {
		processed[current] = true
		for j := 1; j < size; j++ {
			if processed[j] {
				continue
			}
			if found, ok := findOverlap(data[current], data[j], overlap[current], errors); ok {
				nextID[current] = j
				overlap[current] = found
			}
		}
		if nextID[current] < 0 {
			return ""
		}
		current = nextID[current]
	}

	nextID[current] = 0
	if found, ok := findOverlap(data[current], data[0], 0, errors); ok {
		overlap[current] = found
	} else {
		overlap[current] = 0
	}
	lastOverlap := overlap[current]

	type placement struct {
		read  int
		start int
	}
	placements := make([]placement, 0, 2*size)
	position := 0
	placements = append(placements, placement{read: 0, start: 0})
	current = 0
	for step := 0; step < size-1; step++ {
		position += readLength - overlap[current]
		placements = append(placements, placement{read: nextID[current], start: position})
		current = nextID[current]
	}

	genomeLength := position + readLength
	for step := size - 1; step < size*2; step++ {
		position += readLength - overlap[current]
		placements = append(placements, placement{read: nextID[current], start: position})
		current = nextID[current]
	}

	var builder strings.Builder
	for i := 0; i < genomeLength; i++ {
		counts := [256]int{}
		for _, placement := range placements {
			offset := i - placement.start
			if offset >= 0 && offset < readLength {
				counts[data[placement.read][offset]]++
			}
		}
		bestChar := byte(0)
		bestCount := 0
		for _, char := range []byte{'A', 'C', 'G', 'T'} {
			if counts[char] > bestCount {
				bestCount = counts[char]
				bestChar = char
			}
		}
		if bestChar == 0 {
			for char := 0; char < 256; char++ {
				if counts[char] > bestCount {
					bestCount = counts[char]
					bestChar = byte(char)
				}
			}
		}
		builder.WriteByte(bestChar)
	}
	result := builder.String()
	if lastOverlap > 0 && lastOverlap < len(result) {
		result = result[lastOverlap:]
	}
	return result
}

func assembleErrorProne(reads []string) string {
	if len(reads) == 0 {
		return ""
	}
	readLength := len(reads[0])
	errors := readLength / 50
	data := removeDuplicates(reads, readLength, errors)
	if len(data) <= 9 {
		if genome := assembly.CircularSuperstring(data); genome != "" {
			return genome
		}
	}
	return assembleGreedy(data, readLength, errors)
}

func parseReads(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 1 {
		if count, err := strconv.Atoi(lines[0]); err == nil && count == len(lines)-1 {
			return lines[1:]
		}
	}
	return lines
}

func main() {
	reads := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			reads = append(reads, line)
		}
	}
	if len(reads) > 1 {
		if count, err := strconv.Atoi(reads[0]); err == nil && count == len(reads)-1 {
			reads = reads[1:]
		}
	}
	if len(reads) > 0 {
		fmt.Println(assembleErrorProne(reads))
	}
}
