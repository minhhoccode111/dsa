// go run change.go

// Compute the minimum number of coins needed to change the given value into
// coins with denominations 1, 5, and 10.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func change(money int) int {
	return money/10 + (money%10)/5 + money%5
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var money int
	if _, err := fmt.Fscan(reader, &money); err != nil {
		return
	}

	fmt.Println(change(money))
}
