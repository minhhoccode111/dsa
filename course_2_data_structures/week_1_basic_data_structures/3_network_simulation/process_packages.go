// go run process_packages.go

// Simulate processing of network packages with a bounded buffer.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Request struct {
	arrivedAt     int
	timeToProcess int
}

func processRequests(bufferSize int, requests []Request) []int {
	startTimes := make([]int, 0, len(requests))
	finishTimes := make([]int, 0, len(requests))
	head := 0

	for _, request := range requests {
		for head < len(finishTimes) && finishTimes[head] <= request.arrivedAt {
			head++
		}

		if len(finishTimes)-head >= bufferSize {
			startTimes = append(startTimes, -1)
			continue
		}

		start := request.arrivedAt
		if len(finishTimes) > head && finishTimes[len(finishTimes)-1] > start {
			start = finishTimes[len(finishTimes)-1]
		}

		startTimes = append(startTimes, start)
		finishTimes = append(finishTimes, start+request.timeToProcess)
	}

	return startTimes
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var bufferSize, n int
	if _, err := fmt.Fscan(reader, &bufferSize, &n); err != nil {
		return
	}

	requests := make([]Request, n)
	for i := range requests {
		if _, err := fmt.Fscan(reader, &requests[i].arrivedAt, &requests[i].timeToProcess); err != nil {
			return
		}
	}

	for _, start := range processRequests(bufferSize, requests) {
		fmt.Println(start)
	}
}
