package main

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func compareOutput(got, want []byte, opts Options, mode compareMode) (bool, string) {
	if opts.Exact {
		if bytes.Equal(got, want) {
			return true, ""
		}
		return false, "byte-exact output differs"
	}
	if mode == modeSmoke {
		return true, ""
	}

	gotLines := normalizedLines(got)
	wantLines := normalizedLines(want)

	switch mode {
	case modeUnorderedLines:
		gotLines = append([]string(nil), gotLines...)
		wantLines = append([]string(nil), wantLines...)
		sort.Strings(gotLines)
		sort.Strings(wantLines)
		return compareLines(gotLines, wantLines, opts.Float)
	case modeDropFirstUnordered:
		if len(gotLines) > 0 {
			gotLines = gotLines[1:]
		}
		gotLines = append([]string(nil), gotLines...)
		wantLines = append([]string(nil), wantLines...)
		sort.Strings(gotLines)
		sort.Strings(wantLines)
		return compareLines(gotLines, wantLines, opts.Float)
	case modeUnorderedTokens:
		gotTokens := allTokens(gotLines)
		wantTokens := allTokens(wantLines)
		sort.Strings(gotTokens)
		sort.Strings(wantTokens)
		return compareTokenSlices(gotTokens, wantTokens, opts.Float)
	case modeRotationLines:
		return compareRotations(gotLines, wantLines)
	default:
		return compareLines(gotLines, wantLines, opts.Float)
	}
}

func normalizedLines(data []byte) []string {
	value := strings.ReplaceAll(string(data), "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	lines := strings.Split(value, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func compareLines(gotLines, wantLines []string, tolerance float64) (bool, string) {
	if len(gotLines) != len(wantLines) {
		return false, fmt.Sprintf("line count: got %d, want %d", len(gotLines), len(wantLines))
	}
	for i := range wantLines {
		if tolerance > 0 {
			if equalTokens(gotLines[i], wantLines[i], tolerance) {
				continue
			}
			return false, lineDetail(i, gotLines[i], wantLines[i])
		}
		if gotLines[i] == wantLines[i] {
			continue
		}
		return false, lineDetail(i, gotLines[i], wantLines[i])
	}
	return true, ""
}

func equalTokens(gotLine, wantLine string, tolerance float64) bool {
	ok, _ := compareTokenSlices(strings.Fields(gotLine), strings.Fields(wantLine), tolerance)
	return ok
}

func compareTokenSlices(got, want []string, tolerance float64) (bool, string) {
	if len(got) != len(want) {
		return false, fmt.Sprintf("token count: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] == want[i] {
			continue
		}
		if !isFloatToken(got[i]) && !isFloatToken(want[i]) {
			return false, fmt.Sprintf("token %d: got %q, want %q", i+1, got[i], want[i])
		}
		first, errFirst := strconv.ParseFloat(got[i], 64)
		second, errSecond := strconv.ParseFloat(want[i], 64)
		if errFirst != nil || errSecond != nil || math.IsNaN(first) || math.IsNaN(second) ||
			math.IsInf(first, 0) || math.IsInf(second, 0) {
			return false, fmt.Sprintf("token %d: got %q, want %q", i+1, got[i], want[i])
		}
		scale := math.Max(1, math.Max(math.Abs(first), math.Abs(second)))
		if math.Abs(first-second) > tolerance*scale {
			return false, fmt.Sprintf("token %d: got %q, want %q", i+1, got[i], want[i])
		}
	}
	return true, ""
}

func compareRotations(gotLines, wantLines []string) (bool, string) {
	if len(gotLines) != len(wantLines) {
		return false, fmt.Sprintf("line count: got %d, want %d", len(gotLines), len(wantLines))
	}
	for i := range wantLines {
		got := strings.Fields(gotLines[i])
		want := strings.Fields(wantLines[i])
		if !rotationEqual(got, want) {
			return false, lineDetail(i, gotLines[i], wantLines[i])
		}
	}
	return true, ""
}

func rotationEqual(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	if len(got) == 0 {
		return true
	}
	if len(got) == 1 {
		return stringRotation(got[0], want[0])
	}
	for offset := 0; offset < len(want); offset++ {
		match := true
		for i := range want {
			if got[i] != want[(i+offset)%len(want)] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func stringRotation(got, want string) bool {
	if len(got) != len(want) || len(got) == 0 {
		return len(got) == len(want)
	}
	for offset := 0; offset < len(want); offset++ {
		match := true
		for i := 0; i < len(want); i++ {
			if got[i] != want[(i+offset)%len(want)] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func allTokens(lines []string) []string {
	var result []string
	for _, line := range lines {
		result = append(result, strings.Fields(line)...)
	}
	return result
}

func isFloatToken(value string) bool {
	return strings.ContainsAny(value, ".eE")
}

func lineDetail(line int, got, want string) string {
	const max = 120
	if len(got) > max {
		got = got[:max] + "..."
	}
	if len(want) > max {
		want = want[:max] + "..."
	}
	return fmt.Sprintf("line %d: got %q, want %q", line+1, got, want)
}
