// go run substring_equality.go

// Answer substring equality queries using double hashing.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	substringMod1 = int64(1000000007)
	substringMod2 = int64(1000000009)
	substringX    = int64(263)
)

type substringSolver struct {
	s      string
	hash1  []int64
	hash2  []int64
	power1 []int64
	power2 []int64
}

func newSubstringSolver(s string) *substringSolver {
	n := len(s)
	solver := &substringSolver{
		s:      s,
		hash1:  make([]int64, n+1),
		hash2:  make([]int64, n+1),
		power1: make([]int64, n+1),
		power2: make([]int64, n+1),
	}
	solver.power1[0] = 1
	solver.power2[0] = 1

	for i := 1; i <= n; i++ {
		solver.hash1[i] = (substringX*solver.hash1[i-1] + int64(s[i-1])) % substringMod1
		solver.hash2[i] = (substringX*solver.hash2[i-1] + int64(s[i-1])) % substringMod2
		solver.power1[i] = solver.power1[i-1] * substringX % substringMod1
		solver.power2[i] = solver.power2[i-1] * substringX % substringMod2
	}

	return solver
}

func substringHash(hashes, powers []int64, start, length int, mod int64) int64 {
	value := (hashes[start+length] - powers[length]*hashes[start]) % mod
	if value < 0 {
		value += mod
	}
	return value
}

func (solver *substringSolver) ask(a, b, length int) bool {
	return substringHash(solver.hash1, solver.power1, a, length, substringMod1) ==
		substringHash(solver.hash1, solver.power1, b, length, substringMod1) &&
		substringHash(solver.hash2, solver.power2, a, length, substringMod2) ==
			substringHash(solver.hash2, solver.power2, b, length, substringMod2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	solver := newSubstringSolver(s)
	answers := make([]string, q)
	for i := range answers {
		var a, b, length int
		if _, err := fmt.Fscan(reader, &a, &b, &length); err != nil {
			return
		}
		if solver.ask(a, b, length) {
			answers[i] = "Yes"
		} else {
			answers[i] = "No"
		}
	}

	fmt.Println(strings.Join(answers, "\n"))
}
