package twosat

import (
	"math/rand"
	"testing"
)

func TestTwosatAgainstBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for round := 0; round < 5000; round++ {
		n := 1 + rng.Intn(6)
		m := rng.Intn(12)
		cs := make([]Clause, m)
		for i := range cs {
			a := 1 + rng.Intn(n)
			if rng.Intn(2) == 0 {
				a = -a
			}
			b := 1 + rng.Intn(n)
			if rng.Intn(2) == 0 {
				b = -b
			}
			cs[i] = Clause{a, b}
		}
		got, ok := Solve(n, cs)
		want := false
		for mask := 0; mask < 1<<n; mask++ {
			good := true
			for _, c := range cs {
				sat := false
				for _, l := range c {
					v := mask&(1<<(absInt(l)-1)) != 0
					if l < 0 {
						v = !v
					}
					if v {
						sat = true
					}
				}
				if !sat {
					good = false
					break
				}
			}
			if good {
				want = true
				break
			}
		}
		if got != nil {
			for _, c := range cs {
				sat := false
				for _, l := range c {
					v := got[absInt(l)]
					if l < 0 {
						v = !v
					}
					if v {
						sat = true
					}
				}
				if !sat {
					t.Fatalf("round %d bad assignment %v clauses %v", round, got, cs)
				}
			}
		}
		if ok != want {
			t.Fatalf("round %d got %v want %v clauses %v", round, ok, want, cs)
		}
	}
}
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
