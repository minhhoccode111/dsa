// go run fibonacci_sum_squares.go

// Compute the last digit of F(0)^2 + F(1)^2 + ... + F(n)^2.
// Identity: the sum equals F(n) * F(n+1).

package main

import (
	"bufio"
	"fmt"
	"os"
)

const pisanoPeriod10 = 60

func fibonacciMod10(n int64) int64 {
	n %= pisanoPeriod10

	previous, current := int64(0), int64(1)
	for i := int64(0); i < n; i++ {
		previous, current = current, (previous+current)%10
	}
	return previous
}

func fibonacciSumSquaresLastDigit(n int64) int64 {
	return fibonacciMod10(n) * fibonacciMod10(n+1) % 10
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	fmt.Println(fibonacciSumSquaresLastDigit(n))
}
