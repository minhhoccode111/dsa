// go run fibonacci.go

// Compute the nth Fibonacci number.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func fibonacci(n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	previous, current := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		previous, current = current, previous+current
	}
	return current
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	fmt.Println(fibonacci(n))
}
