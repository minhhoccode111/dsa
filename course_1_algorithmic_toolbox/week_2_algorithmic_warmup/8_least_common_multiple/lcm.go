// go run lcm.go

// Compute the least common multiple using lcm(a, b) = a / gcd(a, b) * b.
// Divide before multiplying to keep the product inside int64.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / gcd(a, b) * b
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var a, b int64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}

	fmt.Println(lcm(a, b))
}
