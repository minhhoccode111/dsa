package assembly

import (
	"reflect"
	"testing"
)

func TestFindEulerianCycle(t *testing.T) {
	edges := [][2]int{{0, 1}, {1, 2}, {2, 0}}
	path, ok := FindEulerianCycle(3, edges)
	if !ok {
		t.Fatal("expected a cycle")
	}
	if len(path) != len(edges)+1 || path[0] != path[len(path)-1] {
		t.Fatalf("path = %v", path)
	}
	if path, ok := FindEulerianCycle(3, [][2]int{{0, 1}, {1, 0}, {0, 1}}); ok {
		t.Fatalf("unbalanced graph returned cycle %v", path)
	}
}

func TestFindEulerianCycleDisconnected(t *testing.T) {
	edges := [][2]int{{0, 1}, {1, 2}, {2, 0}, {3, 4}, {4, 3}}
	if path, ok := FindEulerianCycle(5, edges); ok {
		t.Fatalf("disconnected graph returned cycle %v", path)
	}
}

func TestCircularSuperstring(t *testing.T) {
	reads := []string{"AAC", "ACG", "GAA", "GTT", "TCG"}
	got := CircularSuperstring(reads)
	if len(got) != 8 {
		t.Fatalf("length = %d, want 8", len(got))
	}
	for _, read := range reads {
		if !containsCircular(got, read) {
			t.Fatalf("%q does not contain %q", got, read)
		}
	}
}

func TestUnitigContigs(t *testing.T) {
	kmers := []string{"ATG", "ATG", "TGT", "TGG", "CAT", "GGA", "GAT", "AGA"}
	got := UnitigContigs(kmers)
	want := []string{"AGA", "ATG", "ATG", "CAT", "GAT", "TGGA", "TGT"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UnitigContigs = %v, want %v", got, want)
	}
}

func TestUnitigPureCycle(t *testing.T) {
	kmers := []string{"ATG", "TGC", "GCA", "CAT"}
	got := UnitigContigs(kmers)
	if len(got) != 1 || len(got[0]) != len(kmers) || !containsCircular(got[0], kmers[0]) {
		t.Fatalf("UnitigContigs = %v", got)
	}
}

func containsCircular(haystack, needle string) bool {
	doubled := haystack + haystack
	for i := 0; i+len(needle) <= len(doubled); i++ {
		if doubled[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
