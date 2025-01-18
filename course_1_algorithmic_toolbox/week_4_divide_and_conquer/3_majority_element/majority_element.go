// go run majority_element.go

// Check whether a sequence contains an element that appears more than half
// of the time.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func majorityElement(elements []int) int {
	if len(elements) == 0 {
		return 0
	}

	candidate := 0
	count := 0
	for _, element := range elements {
		if count == 0 {
			candidate = element
			count = 1
		} else if element == candidate {
			count++
		} else {
			count--
		}
	}

	count = 0
	for _, element := range elements {
		if element == candidate {
			count++
		}
	}

	if count*2 > len(elements) {
		return 1
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	elements := make([]int, n)
	for i := range elements {
		if _, err := fmt.Fscan(reader, &elements[i]); err != nil {
			return
		}
	}

	fmt.Println(majorityElement(elements))
}
