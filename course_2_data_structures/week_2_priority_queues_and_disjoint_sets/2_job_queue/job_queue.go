// go run job_queue.go

// Simulate processing jobs in parallel.

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type AssignedJob struct {
	Worker    int
	StartedAt int64
}

type workerState struct {
	free  int64
	index int
}

type workerHeap []workerState

func (h workerHeap) Len() int { return len(h) }
func (h workerHeap) Less(i, j int) bool {
	if h[i].free != h[j].free {
		return h[i].free < h[j].free
	}
	return h[i].index < h[j].index
}
func (h workerHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *workerHeap) Push(x any)   { *h = append(*h, x.(workerState)) }
func (h *workerHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]
	return last
}

func assignJobs(nWorkers int, jobs []int) []AssignedJob {
	result := make([]AssignedJob, 0, len(jobs))
	workers := make(workerHeap, nWorkers)
	for i := range workers {
		workers[i] = workerState{free: 0, index: i}
	}
	heap.Init(&workers)

	for _, duration := range jobs {
		worker := heap.Pop(&workers).(workerState)
		result = append(result, AssignedJob{Worker: worker.index, StartedAt: worker.free})
		heap.Push(&workers, workerState{free: worker.free + int64(duration), index: worker.index})
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var nWorkers, nJobs int
	if _, err := fmt.Fscan(reader, &nWorkers, &nJobs); err != nil {
		return
	}

	jobs := make([]int, nJobs)
	for i := range jobs {
		if _, err := fmt.Fscan(reader, &jobs[i]); err != nil {
			return
		}
	}

	for _, job := range assignJobs(nWorkers, jobs) {
		fmt.Println(job.Worker, job.StartedAt)
	}
}
