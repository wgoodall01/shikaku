package shikaku

import (
	"sort"
	"testing"
)

func TestIn(t *testing.T) {
	lo := Vec2{1, 1}
	hi := Vec2{3, 3}

	if (Vec2{2, 2}).In(lo, hi) != true {
		t.Log("1<2<3 didn't work")
		t.Fail()
	}
}

func TestMath(t *testing.T) {
	if (Vec2{1, 2}).Add(Vec2{3, 5}) != (Vec2{4, 7}) {
		t.Fail()
		t.Log("Can't add properly")
	}

	if (Vec2{4, 7}).Sub(Vec2{3, 5}) != (Vec2{1, 2}) {
		t.Fail()
		t.Log("Can't subtract properly")
	}
}

func TestFactor(t *testing.T) {
	// In a list of factors, makes sure all pairs are in pair[0] <= pair[1] order.
	swapFactorList := func(facs []Vec2) {
		for _, p := range facs {
			if p[0] > p[1] {
				tmp := p[1]
				p[1] = p[0]
				p[0] = tmp
			}
		}
	}

	// Sort a list of factors by their first term
	sortFactorList := func(facs []Vec2) {
		sort.Slice(facs, func(i, j int) bool { return facs[i][0] < facs[j][0] })
	}

	type testCase struct {
		X       int
		Factors []Vec2
	}

	tests := []testCase{
		{5, []Vec2{{1, 5}}},
		{42, []Vec2{{1, 42}, {2, 21}, {3, 14}, {6, 7}}},
		{16, []Vec2{{1, 16}, {2, 8}, {4, 4}}},
		{32, []Vec2{{1, 32}, {2, 16}, {4, 8}}},
		{33, []Vec2{{1, 33}, {3, 11}}},
		{100, []Vec2{{1, 100}, {2, 50}, {4, 25}, {5, 20}, {10, 10}}},
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
