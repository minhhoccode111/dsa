// go run largest_number.go

// Compile the largest integer by concatenating the given integers.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func largestNumber(numbers []string) string {
	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i]+numbers[j] > numbers[j]+numbers[i]
	})

	result := strings.Join(numbers, "")
	if result != "" && result[0] == '0' {
		return "0"
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	numbers := make([]string, n)
	for i := range numbers {
		if _, err := fmt.Fscan(reader, &numbers[i]); err != nil {
			return
		}
	}

	fmt.Println(largestNumber(numbers))
}
