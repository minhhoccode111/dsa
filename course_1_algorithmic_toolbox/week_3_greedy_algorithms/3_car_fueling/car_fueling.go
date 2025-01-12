// go run car_fueling.go

// Compute the minimum number of gas tank refills to get from one city to another.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func minRefills(distance, tank int, stops []int) int {
	positions := make([]int, 0, len(stops)+1)
	positions = append(positions, stops...)
	positions = append(positions, distance)

	refills := 0
	lastPosition := 0
	index := 0

	for lastPosition+tank < distance {
		currentPosition := lastPosition
		for index < len(positions) && positions[index] <= lastPosition+tank {
			currentPosition = positions[index]
			index++
		}

		if currentPosition == lastPosition {
			return -1
		}

		lastPosition = currentPosition
		refills++
	}

	return refills
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var distance, tank, n int
	if _, err := fmt.Fscan(reader, &distance, &tank, &n); err != nil {
		return
	}

	stops := make([]int, n)
	for i := range stops {
		if _, err := fmt.Fscan(reader, &stops[i]); err != nil {
			return
		}
	}

	fmt.Println(minRefills(distance, tank, stops))
}
