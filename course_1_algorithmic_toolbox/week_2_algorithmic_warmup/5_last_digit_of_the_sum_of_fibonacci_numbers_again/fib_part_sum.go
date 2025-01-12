// go run fib_part_sum.go

// Compute the last digit of F(m) + F(m+1) + ... + F(n).
// Identity: the partial sum equals F(n+2) - F(m+1).

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

func fibonacciPartialSumLastDigit(m, n int64) int64 {
	difference := fibonacciMod10(n+2) - fibonacciMod10(m+1)
	return (difference + 10) % 10
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var m, n int64
	if _, err := fmt.Fscan(reader, &m, &n); err != nil {
		return
	}

	fmt.Println(fibonacciPartialSumLastDigit(m, n))
}
