// go run closest.go

// Find the closest pair of points in a set of points on a plane.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct {
	x, y int
}

const maxDistance int64 = 1<<63 - 1

func distanceSquared(a, b Point) int64 {
	dx := int64(a.x) - int64(b.x)
	dy := int64(a.y) - int64(b.y)
	return dx*dx + dy*dy
}

func closestPair(points []Point) float64 {
	if len(points) < 2 {
		return 0
	}

	sorted := append([]Point(nil), points...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].x != sorted[j].x {
			return sorted[i].x < sorted[j].x
		}
		return sorted[i].y < sorted[j].y
	})

	return math.Sqrt(float64(closestDistance(sorted)))
}

func closestDistance(points []Point) int64 {
	n := len(points)
	if n < 2 {
		return maxDistance
	}

	if n <= 3 {
		best := maxDistance
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if d := distanceSquared(points[i], points[j]); d < best {
					best = d
				}
			}
		}
		return best
	}

	mid := n / 2
	midX := points[mid].x

	best := closestDistance(points[:mid])
	if rightBest := closestDistance(points[mid:]); rightBest < best {
		best = rightBest
	}

	strip := make([]Point, 0, n)
	for _, point := range points {
		dx := int64(point.x) - int64(midX)
		if dx*dx < best {
			strip = append(strip, point)
		}
	}

	sort.Slice(strip, func(i, j int) bool { return strip[i].y < strip[j].y })

	for i := 0; i < len(strip); i++ {
		for j := i + 1; j < len(strip) && j <= i+7; j++ {
			dy := int64(strip[j].y) - int64(strip[i].y)
			if dy*dy >= best {
				break
			}
			if d := distanceSquared(strip[i], strip[j]); d < best {
				best = d
			}
		}
	}

	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	points := make([]Point, n)
	for i := range points {
		if _, err := fmt.Fscan(reader, &points[i].x, &points[i].y); err != nil {
			return
		}
	}

	fmt.Printf("%.9f\n", closestPair(points))
}
