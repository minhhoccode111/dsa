package main

import (
	"reflect"
	"testing"
)

func TestAssignJobs(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		jobs    []int
		want    []AssignedJob
	}{
		{
			"course sample",
			2,
			[]int{1, 2, 3, 4, 5},
			[]AssignedJob{{0, 0}, {1, 0}, {0, 1}, {1, 2}, {0, 4}},
		},
		{
			"single worker",
			1,
			[]int{1, 2, 3},
			[]AssignedJob{{0, 0}, {0, 1}, {0, 3}},
		},
		{
			"zero duration reuses first worker",
			2,
			[]int{0, 0, 0},
			[]AssignedJob{{0, 0}, {0, 0}, {0, 0}},
		},
		{
			"no jobs",
			2,
			nil,
			[]AssignedJob{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := assignJobs(tt.workers, tt.jobs); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("assignJobs(%d, %v) = %v, want %v", tt.workers, tt.jobs, got, tt.want)
			}
		})
	}
}
