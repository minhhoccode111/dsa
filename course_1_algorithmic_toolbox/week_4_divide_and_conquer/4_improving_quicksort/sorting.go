// go run sorting.go

// Sort a sequence that may contain duplicates using 3-way quicksort.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func sortThreeWay(a []int) {
	quickSort3(a, 0, len(a)-1)
}

func quickSort3(a []int, left, right int) {
	if left >= right {
		return
	}

	lt, gt := partition3(a, left, right)
	quickSort3(a, left, lt-1)
	quickSort3(a, gt+1, right)
}

func partition3(a []int, left, right int) (int, int) {
	mid := left + (right-left)/2

	if a[mid] < a[left] {
		a[left], a[mid] = a[mid], a[left]
	}
	if a[right] < a[left] {
		a[left], a[right] = a[right], a[left]
	}
	if a[right] < a[mid] {
		a[mid], a[right] = a[right], a[mid]
	}
	a[left], a[mid] = a[mid], a[left]

	pivot := a[left]
	lt, gt := left, right

	for i := left; i <= gt; {
		switch {
		case a[i] < pivot:
			a[lt], a[i] = a[i], a[lt]
			lt++
			i++
		case a[i] > pivot:
			a[i], a[gt] = a[gt], a[i]
			gt--
		default:
			i++
		}
	}

	return lt, gt
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

	sortThreeWay(elements)

	for i, element := range elements {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(element)
	}
	fmt.Println()
}
