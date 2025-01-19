// go run placing_parentheses.go

// Parenthesize an arithmetic expression to maximize its value.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func evaluate(a, b int, operator byte) int {
	switch operator {
	case '+':
		return a + b
	case '-':
		return a - b
	default: // '*'
		return a * b
	}
}

func maxArithmeticValue(expression string) int {
	numbers := make([]int, 0, len(expression)/2+1)
	operators := make([]byte, 0, len(expression)/2)

	for i := 0; i < len(expression); i++ {
		if i%2 == 0 {
			numbers = append(numbers, int(expression[i]-'0'))
		} else {
			operators = append(operators, expression[i])
		}
	}

	n := len(numbers)
	maxDP := make([][]int, n)
	minDP := make([][]int, n)
	for i := 0; i < n; i++ {
		maxDP[i] = make([]int, n)
		minDP[i] = make([]int, n)
		maxDP[i][i] = numbers[i]
		minDP[i][i] = numbers[i]
	}

	const infinity = 1 << 60

	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			maxValue := -infinity
			minValue := infinity

			for k := i; k < j; k++ {
				operator := operators[k]
				candidates := [4]int{
					evaluate(maxDP[i][k], maxDP[k+1][j], operator),
					evaluate(maxDP[i][k], minDP[k+1][j], operator),
					evaluate(minDP[i][k], maxDP[k+1][j], operator),
					evaluate(minDP[i][k], minDP[k+1][j], operator),
				}

				for _, candidate := range candidates {
					if candidate > maxValue {
						maxValue = candidate
					}
					if candidate < minValue {
						minValue = candidate
					}
				}
			}

			maxDP[i][j] = maxValue
			minDP[i][j] = minValue
		}
	}

	return maxDP[0][n-1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var expression string
	if _, err := fmt.Fscan(reader, &expression); err != nil {
		return
	}

	fmt.Println(maxArithmeticValue(expression))
}
