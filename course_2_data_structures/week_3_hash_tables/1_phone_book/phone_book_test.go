package main

import (
	"reflect"
	"testing"
)

func TestProcessPhoneQueries(t *testing.T) {
	queries := []phoneQuery{
		{"add", "911", "police"},
		{"add", "76213", "Mom"},
		{"add", "17239", "Bob"},
		{"find", "76213", ""},
		{"find", "910", ""},
		{"find", "911", ""},
		{"del", "910", ""},
		{"del", "911", ""},
		{"find", "911", ""},
		{"find", "76213", ""},
		{"add", "76213", "daddy"},
		{"find", "76213", ""},
	}

	want := []string{"Mom", "not found", "police", "not found", "Mom", "daddy"}
	if got := processPhoneQueries(queries); !reflect.DeepEqual(got, want) {
		t.Fatalf("processPhoneQueries() = %v, want %v", got, want)
	}
}
