package main

import "testing"

func TestStackWithMax(t *testing.T) {
	stack := newStackWithMax()
	stack.Push(2)
	stack.Push(1)
	if got := stack.Max(); got != 2 {
		t.Fatalf("Max() = %d, want 2", got)
	}
	stack.Pop()
	if got := stack.Max(); got != 2 {
		t.Fatalf("Max() after pop = %d, want 2", got)
	}

	stack.Push(5)
	stack.Push(5)
	if got := stack.Max(); got != 5 {
		t.Fatalf("Max() with duplicates = %d, want 5", got)
	}
	stack.Pop()
	if got := stack.Max(); got != 5 {
		t.Fatalf("Max() after duplicate pop = %d, want 5", got)
	}
	stack.Pop()
	if got := stack.Max(); got != 2 {
		t.Fatalf("Max() after duplicate removed = %d, want 2", got)
	}
}
