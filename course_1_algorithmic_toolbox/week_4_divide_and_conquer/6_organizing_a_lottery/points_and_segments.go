// go run points_and_segments.go

// For each point, compute the number of segments that contain it.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	eventStart = iota
	eventPoint
	eventEnd
)

type event struct {
	coordinate int
	kind       int
	index      int
}

func pointsCovered(starts, ends, points []int) []int {
	events := make([]event, 0, 2*len(starts)+len(points))

	for i := range starts {
		events = append(events,
			event{coordinate: starts[i], kind: eventStart},
			event{coordinate: ends[i], kind: eventEnd},
		)
	}
	for i, point := range points {
		events = append(events, event{coordinate: point, kind: eventPoint, index: i})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].coordinate != events[j].coordinate {
			return events[i].coordinate < events[j].coordinate
		}
		return events[i].kind < events[j].kind
	})

	result := make([]int, len(points))
	active := 0
	for _, event := range events {
		switch event.kind {
		case eventStart:
			active++
		case eventEnd:
			active--
		case eventPoint:
			result[event.index] = active
		}
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	starts := make([]int, n)
	ends := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &starts[i], &ends[i]); err != nil {
			return
		}
	}

	points := make([]int, m)
	for i := range points {
		if _, err := fmt.Fscan(reader, &points[i]); err != nil {
			return
		}
	}

	result := pointsCovered(starts, ends, points)
	strs := make([]string, len(result))
	for i, value := range result {
		strs[i] = fmt.Sprint(value)
	}
	fmt.Println(strings.Join(strs, " "))
}
