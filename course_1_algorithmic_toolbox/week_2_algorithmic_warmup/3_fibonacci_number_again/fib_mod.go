// go run fib_mod.go

// Compute the nth Fibonacci number modulo m using the Pisano period.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func pisanoPeriod(m int64) int64 {
	previous, current := int64(0), int64(1)
	for i := int64(0); i < 6*m; i++ {
		previous, current = current, (previous+current)%m
		if previous == 0 && current == 1 {
			return i + 1
		}
	}
	return 6 * m
}

func fibonacciMod(n, m int64) int64 {
	if m == 1 {
		return 0
	}

	n %= pisanoPeriod(m)

	previous, current := int64(0), int64(1)
	for i := int64(0); i < n; i++ {
		previous, current = current, (previous+current)%m
	}
	return previous
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	fmt.Println(fibonacciMod(n, m))
}
