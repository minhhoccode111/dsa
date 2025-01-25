// go run check_brackets.go

// Find the first error in bracket usage, or report Success.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type bracket struct {
	char     byte
	position int
}

func matching(left, right byte) bool {
	return left == '(' && right == ')' ||
		left == '[' && right == ']' ||
		left == '{' && right == '}'
}

func findMismatch(text string) int {
	stack := make([]bracket, 0, len(text))

	for i := 0; i < len(text); i++ {
		char := text[i]

		switch char {
		case '(', '[', '{':
			stack = append(stack, bracket{char: char, position: i + 1})
		case ')', ']', '}':
			if len(stack) == 0 {
				return i + 1
			}

			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !matching(top.char, char) {
				return i + 1
			}
		}
	}

	if len(stack) > 0 {
		return stack[0].position
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return
	}

	text := strings.TrimRight(line, "\r\n")
	result := findMismatch(text)
	if result == 0 {
		fmt.Println("Success")
	} else {
		fmt.Println(result)
	}
}
