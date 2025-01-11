// go run two_sum.go

package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumOfTwoDigits(a, b int) int {
	return a + b
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}

	fmt.Println(sumOfTwoDigits(a, b))
}
