// go run plan_party.go

// Find the maximum-weight independent set in a company hierarchy tree.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func maxPartyFun(n int, weights []int, edges [][2]int) int64 {
	if n == 0 {
		return 0
	}
	adjacency := make([][]int, n)
	for _, edge := range edges {
		u, v := edge[0]-1, edge[1]-1
		adjacency[u] = append(adjacency[u], v)
		adjacency[v] = append(adjacency[v], u)
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	parent[0] = -1
	order := make([]int, 0, n)
	stack := []int{0}
	for len(stack) > 0 {
		vertex := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, vertex)
		for _, next := range adjacency[vertex] {
			if parent[next] == -2 {
				parent[next] = vertex
				stack = append(stack, next)
			}
		}
	}

	include := make([]int64, n)
	exclude := make([]int64, n)
	for i := len(order) - 1; i >= 0; i-- {
		vertex := order[i]
		include[vertex] = int64(weights[vertex])
		exclude[vertex] = 0
		for _, next := range adjacency[vertex] {
			if parent[next] != vertex {
				continue
			}
			include[vertex] += exclude[next]
			if include[next] > exclude[next] {
				exclude[vertex] += include[next]
			} else {
				exclude[vertex] += exclude[next]
			}
		}
	}
	if include[0] > exclude[0] {
		return include[0]
	}
	return exclude[0]
}

func parseParty(input string) (int, []int, [][2]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, nil, err
	}
	weights := make([]int, n)
	for i := range weights {
		if _, err := fmt.Fscan(reader, &weights[i]); err != nil {
			return 0, nil, nil, err
		}
	}
	edges := make([][2]int, n-1)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return 0, nil, nil, err
		}
	}
	return n, weights, edges, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	weights := make([]int, n)
	for i := range weights {
		if _, err := fmt.Fscan(reader, &weights[i]); err != nil {
			return
		}
	}
	edges := make([][2]int, n-1)
	for i := range edges {
		if _, err := fmt.Fscan(reader, &edges[i][0], &edges[i][1]); err != nil {
			return
		}
	}
	fmt.Println(maxPartyFun(n, weights, edges))
}
