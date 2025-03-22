// go run school_bus.go

// Find the shortest circular route through all vertices with Held-Karp DP.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const busInfinity int64 = 1 << 62

func shortestBusRoute(n int, graph [][]int64) (int64, []int, bool) {
	if n == 1 {
		return 0, []int{0}, true
	}
	states := 1 << n
	dp := make([]int64, states*n)
	for i := range dp {
		dp[i] = busInfinity
	}
	parent := make([]int8, states*n)
	for i := range parent {
		parent[i] = -1
	}
	dp[(1<<0)*n+0] = 0

	for mask := 1; mask < states; mask++ {
		if mask&1 == 0 {
			continue
		}
		for last := 0; last < n; last++ {
			if mask&(1<<last) == 0 {
				continue
			}
			current := dp[mask*n+last]
			if current == busInfinity {
				continue
			}
			for next := 1; next < n; next++ {
				if mask&(1<<next) != 0 || graph[last][next] == busInfinity {
					continue
				}
				nextMask := mask | 1<<next
				candidate := current + graph[last][next]
				if candidate < dp[nextMask*n+next] {
					dp[nextMask*n+next] = candidate
					parent[nextMask*n+next] = int8(last)
				}
			}
		}
	}

	full := states - 1
	best := busInfinity
	lastVertex := -1
	for last := 1; last < n; last++ {
		if dp[full*n+last] == busInfinity || graph[last][0] == busInfinity {
			continue
		}
		candidate := dp[full*n+last] + graph[last][0]
		if candidate < best {
			best = candidate
			lastVertex = last
		}
	}
	if lastVertex < 0 {
		return -1, nil, false
	}

	path := make([]int, 0, n)
	mask := full
	current := lastVertex
	for current != 0 {
		path = append(path, current)
		previous := int(parent[mask*n+current])
		mask ^= 1 << current
		current = previous
	}
	path = append(path, 0)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return best, path, true
}

func parseBus(input string) (int, [][]int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return 0, nil, err
	}
	graph := make([][]int64, n)
	for i := range graph {
		graph[i] = make([]int64, n)
		for j := range graph[i] {
			graph[i][j] = busInfinity
		}
	}
	for i := 0; i < m; i++ {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return 0, nil, err
		}
		graph[u-1][v-1] = weight
		graph[v-1][u-1] = weight
	}
	return n, graph, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	graph := make([][]int64, n)
	for i := range graph {
		graph[i] = make([]int64, n)
		for j := range graph[i] {
			graph[i][j] = busInfinity
		}
	}
	for i := 0; i < m; i++ {
		var u, v int
		var weight int64
		if _, err := fmt.Fscan(reader, &u, &v, &weight); err != nil {
			return
		}
		graph[u-1][v-1] = weight
		graph[v-1][u-1] = weight
	}
	cost, path, ok := shortestBusRoute(n, graph)
	if !ok {
		fmt.Println(-1)
		return
	}
	parts := make([]string, len(path))
	for i, vertex := range path {
		parts[i] = strconv.Itoa(vertex + 1)
	}
	fmt.Println(cost)
	fmt.Println(strings.Join(parts, " "))
}
