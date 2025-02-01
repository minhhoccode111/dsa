// go run common_substring.go

// Find a longest common substring using double rolling hashes.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	commonMod1 = int64(1000000007)
	commonMod2 = int64(1000000009)
	commonX    = int64(263)
)

type stringHashes struct {
	hash1  []int64
	hash2  []int64
	power1 []int64
	power2 []int64
}

func buildStringHashes(s string) *stringHashes {
	n := len(s)
	h := &stringHashes{
		hash1:  make([]int64, n+1),
		hash2:  make([]int64, n+1),
		power1: make([]int64, n+1),
		power2: make([]int64, n+1),
	}
	h.power1[0] = 1
	h.power2[0] = 1
	for i := 1; i <= n; i++ {
		h.hash1[i] = (commonX*h.hash1[i-1] + int64(s[i-1])) % commonMod1
		h.hash2[i] = (commonX*h.hash2[i-1] + int64(s[i-1])) % commonMod2
		h.power1[i] = h.power1[i-1] * commonX % commonMod1
		h.power2[i] = h.power2[i-1] * commonX % commonMod2
	}
	return h
}

func (h *stringHashes) get(start, length int) (int64, int64) {
	value1 := (h.hash1[start+length] - h.power1[length]*h.hash1[start]) % commonMod1
	if value1 < 0 {
		value1 += commonMod1
	}
	value2 := (h.hash2[start+length] - h.power2[length]*h.hash2[start]) % commonMod2
	if value2 < 0 {
		value2 += commonMod2
	}
	return value1, value2
}

type hashPair struct {
	first  int64
	second int64
}

func findCommonLength(s, t string, sHashes, tHashes *stringHashes, length int) (int, int, bool) {
	starts := make(map[hashPair]int)
	for i := 0; i+length <= len(s); i++ {
		key := hashPair{}
		key.first, key.second = sHashes.get(i, length)
		if _, exists := starts[key]; !exists {
			starts[key] = i
		}
	}

	bestI, bestJ := -1, -1
	for j := 0; j+length <= len(t); j++ {
		key := hashPair{}
		key.first, key.second = tHashes.get(j, length)
		i, exists := starts[key]
		if !exists {
			continue
		}
		if bestI == -1 || i < bestI || (i == bestI && j < bestJ) {
			bestI, bestJ = i, j
		}
	}

	if bestI == -1 {
		return 0, 0, false
	}
	return bestI, bestJ, true
}

func longestCommonSubstring(s, t string) (int, int, int) {
	if len(s) == 0 || len(t) == 0 {
		return 0, 0, 0
	}

	sHashes := buildStringHashes(s)
	tHashes := buildStringHashes(t)

	low, high := 0, len(s)
	if len(t) < high {
		high = len(t)
	}

	bestI, bestJ, bestLength := 0, 0, 0
	for low <= high {
		mid := low + (high-low)/2
		if mid == 0 {
			low = 1
			continue
		}
		i, j, found := findCommonLength(s, t, sHashes, tHashes, mid)
		if found {
			bestI, bestJ, bestLength = i, j, mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return bestI, bestJ, bestLength
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024), 1<<20)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		i, j, length := longestCommonSubstring(fields[0], fields[1])
		fmt.Println(i, j, length)
	}
}
