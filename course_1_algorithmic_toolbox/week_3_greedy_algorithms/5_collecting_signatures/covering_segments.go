// go run covering_segments.go

// Find the minimum number of points needed to cover all given segments on a line.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func optimalPoints(segments [][2]int) []int {
	sort.Slice(segments, func(i, j int) bool {
		return segments[i][1] < segments[j][1]
	})

	points := make([]int, 0, len(segments))
	lastPoint := 0
	hasPoint := false

	for _, segment := range segments {
		if !hasPoint || segment[0] > lastPoint {
			lastPoint = segment[1]
			points = append(points, lastPoint)
			hasPoint = true
		}
	}

	return points
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	segments := make([][2]int, n)
	for i := range segments {
		if _, err := fmt.Fscan(reader, &segments[i][0], &segments[i][1]); err != nil {
			return
		}
	}

	points := optimalPoints(segments)

	fmt.Println(len(points))
	for i, point := range points {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(point)
	}
	fmt.Println()
}
