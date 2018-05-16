package shikaku

import (
	"sort"
	"testing"
)

func TestFactor(t *testing.T) {
	// In a list of factors, makes sure all pairs are in pair[0] <= pair[1] order.
	swapFactorList := func(facs []Pair) {
		for _, p := range facs {
			if p[0] > p[1] {
				tmp := p[1]
				p[1] = p[0]
				p[0] = tmp
			}
		}
	}

	// Sort a list of factors by their first term
	sortFactorList := func(facs []Pair) {
		sort.Slice(facs, func(i, j int) bool { return facs[i][0] < facs[j][0] })
	}

	type testCase struct {
		X       int
		Factors []Pair
	}

	tests := []testCase{
		{5, []Pair{{1, 5}}},
		{42, []Pair{{1, 42}, {2, 21}, {3, 14}, {6, 7}}},
		{16, []Pair{{1, 16}, {2, 8}, {4, 4}}},
		{32, []Pair{{1, 32}, {2, 16}, {4, 8}}},
		{33, []Pair{{1, 33}, {3, 11}}},
		{100, []Pair{{1, 100}, {2, 50}, {4, 25}, {5, 20}, {10, 10}}},
	}

	for _, test := range tests {
		x := test.X
		testFactors := Factor(x)
		realFactors := test.Factors

		// Swap factors, lowest first.
		swapFactorList(testFactors)
		swapFactorList(realFactors)

		// Sort lists
		sortFactorList(testFactors)
		sortFactorList(realFactors)

		// Compare lists and assert
		failed := false

		if len(testFactors) != len(realFactors) {
			failed = true
		} else {
			for i := range testFactors {
				if testFactors[i] != realFactors[i] {
					failed = true
					break
				}
			}
		}

		if failed {
			t.Logf("Incorrect factors for %d", x)
			t.Logf("       got: %v", testFactors)
			t.Logf("  expected: %v", realFactors)
			t.Fail()
		}
	}
}
