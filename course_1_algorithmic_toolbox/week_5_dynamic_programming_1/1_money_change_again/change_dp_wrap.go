// go run change_dp_wrap.go

// Compute the minimum number of coins needed to change the given value into
// coins with denominations 1, 3, and 4.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func minCoins(money int) int {
	coins := []int{1, 3, 4}
	dp := make([]int, money+1)

	for value := 1; value <= money; value++ {
		dp[value] = value
		for _, coin := range coins {
			if value >= coin && dp[value-coin]+1 < dp[value] {
				dp[value] = dp[value-coin] + 1
			}
		}
	}

	return dp[money]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var money int
	if _, err := fmt.Fscan(reader, &money); err != nil {
		return
	}

	fmt.Println(minCoins(money))
}
