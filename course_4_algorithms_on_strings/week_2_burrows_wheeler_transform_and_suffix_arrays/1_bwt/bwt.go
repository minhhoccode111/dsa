// go run bwt.go

// Construct the Burrows-Wheeler transform of a string.

package main

import (
	"bufio"
	"fmt"
	"os"

	"dsa/course_4_algorithms_on_strings/internal/strs"
)

func bwt(text string) string {
	n := len(text)
	if n == 0 {
		return ""
	}
	data := []byte(text)
	sa := strs.SuffixArray(data)
	result := make([]byte, n)
	for i, suffix := range sa {
		previous := suffix - 1
		if previous < 0 {
			previous = n - 1
		}
		result[i] = data[previous]
	}
	return string(result)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var text string
	if _, err := fmt.Fscan(reader, &text); err != nil {
		return
	}
	fmt.Println(bwt(text))
}
