// go run stack_with_max.go

// Implement a stack that also reports its maximum in constant time.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type stackWithMax struct {
	values []int
	maxima []int
}

func newStackWithMax() *stackWithMax {
	return &stackWithMax{}
}

func (s *stackWithMax) Push(value int) {
	s.values = append(s.values, value)
	if len(s.maxima) == 0 || value >= s.maxima[len(s.maxima)-1] {
		s.maxima = append(s.maxima, value)
	}
}

func (s *stackWithMax) Pop() {
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	if value == s.maxima[len(s.maxima)-1] {
		s.maxima = s.maxima[:len(s.maxima)-1]
	}
}

func (s *stackWithMax) Max() int {
	return s.maxima[len(s.maxima)-1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var numQueries int
	if _, err := fmt.Fscan(reader, &numQueries); err != nil {
		return
	}

	stack := newStackWithMax()
	for i := 0; i < numQueries; i++ {
		var query string
		if _, err := fmt.Fscan(reader, &query); err != nil {
			return
		}

		switch query {
		case "push":
			var value int
			if _, err := fmt.Fscan(reader, &value); err != nil {
				return
			}
			stack.Push(value)
		case "pop":
			stack.Pop()
		case "max":
			fmt.Println(stack.Max())
		}
	}
}
