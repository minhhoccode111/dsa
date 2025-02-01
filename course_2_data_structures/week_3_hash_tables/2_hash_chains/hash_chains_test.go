package main

import (
	"reflect"
	"testing"
)

func TestProcessHashQueries(t *testing.T) {
	queries := []hashQuery{
		{"add", 0, "world"},
		{"add", 0, "HellO"},
		{"check", 4, ""},
		{"find", 0, "World"},
		{"find", 0, "world"},
		{"del", 0, "world"},
		{"check", 4, ""},
		{"del", 0, "HellO"},
		{"add", 0, "luck"},
		{"add", 0, "GooD"},
		{"check", 2, ""},
		{"del", 0, "good"},
	}

	want := []string{"HellO world", "no", "yes", "HellO", "GooD luck"}
	if got := processHashQueries(5, queries); !reflect.DeepEqual(got, want) {
		t.Fatalf("processHashQueries() = %v, want %v", got, want)
	}
}

func TestHashQueriesDuplicateAdd(t *testing.T) {
	queries := []hashQuery{
		{"add", 0, "a"},
		{"add", 0, "a"},
		{"check", 0, ""},
	}
	want := []string{"a"}
	if got := processHashQueries(1, queries); !reflect.DeepEqual(got, want) {
		t.Fatalf("duplicate add produced %v, want %v", got, want)
	}
}
