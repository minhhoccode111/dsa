package main

import "testing"

func TestRope(t *testing.T) {
	rope := newRope("hlelowrold")
	rope.process(1, 1, 2)
	if got := rope.result(); got != "hellowrold" {
		t.Fatalf("after first move = %q, want %q", got, "hellowrold")
	}
	rope.process(6, 6, 7)
	if got := rope.result(); got != "helloworld" {
		t.Fatalf("after second move = %q, want %q", got, "helloworld")
	}
}

func TestRopeEdges(t *testing.T) {
	rope := newRope("abc")
	rope.process(0, 0, 2)
	if got := rope.result(); got != "bca" {
		t.Fatalf("move first to end = %q, want %q", got, "bca")
	}

	rope = newRope("abc")
	rope.process(2, 2, 0)
	if got := rope.result(); got != "cab" {
		t.Fatalf("move last to front = %q, want %q", got, "cab")
	}

	rope = newRope("abcdef")
	rope.process(1, 3, 0)
	if got := rope.result(); got != "bcdaef" {
		t.Fatalf("move middle to front = %q, want %q", got, "bcdaef")
	}
}
