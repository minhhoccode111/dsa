// go run fibonacci_last_digit.go

// Compute the last digit of the nth Fibonacci number.
// The last digits repeat with the Pisano period of 10, which is 60.

package main

import (
	"bufio"
	"fmt"
	"os"
)

const pisanoPeriod10 = 60

func fibonacciLastDigit(n int64) int64 {
	n %= pisanoPeriod10

	previous, current := int64(0), int64(1)
	for i := int64(0); i < n; i++ {
		previous, current = current, (previous+current)%10
	}
	return previous
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	fmt.Println(fibonacciLastDigit(n))
}
