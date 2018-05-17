package shikaku

import (
	"testing"
)

func TestBoardParseUnevenLengths(t *testing.T) {
	_, err := NewBoardFromString(`
		-- 02 --
		-- -- -- 04
		-- 03 --
	`)

	if err == nil {
		t.Fatal("Didn't blow up parsing a board with uneven row lengths")
	}
}

func TestBoardParseNotIntegers(t *testing.T) {
	_, err := NewBoardFromString(`
		-- nope what
	`)

	if err == nil {
		t.Fatal("Didn't blow up parsing a board with non-integer integers")
	}
}

func TestBoardParseValid(t *testing.T) {
	bo, err := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	if err != nil {
		t.Fatal("Failed to parse valid board")
	}

	t1 := bo.Grid[0][2].Area == 5
	t2 := bo.Grid[2][1].Area == 2
	t3 := IsNotFinal(bo.Grid[4][4])
	t4 := IsNotFinal(bo.Grid[0][0])

	if !(t1 && t2 && t3 && t4) {
		t.Fatal("Didn't parse valid board correctly.")
	}

	if t.Failed() {
		t.Log("\n" + bo.String())
	}
}

func TestDimensions(t *testing.T) {
	bo, _ := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
	`)

	if bo.Height() != 4 {
		t.Log("Incorrect height")
		t.Fail()
	}

	if bo.Width() != 5 {
		t.Log("Incorrect width")
		t.Fail()
	}

	size := bo.Size()

	if size[0] != 5 {
		t.Log("Incorrect width from Size")
		t.Fail()
	}

	if size[1] != 4 {
		t.Log("Incorrect height from Size")
		t.Fail()
	}

	if t.Failed() {
		t.Log("\n" + bo.String())
	}
}

func TestCollides(t *testing.T) {
	bo, _ := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	// Top left 3x3, collides.
	if bo.Collides(Vec2{0, 0}, Vec2{3, 3}, bo.Get(Vec2{1, 1})) != true {
		t.Fail()
		t.Log("Doesn't collide when it should")
	}

	// Bottom right 3x3, no collides.
	if bo.Collides(Vec2{2, 2}, Vec2{5, 5}, bo.Get(Vec2{3, 3})) != false {
		t.Fail()
		t.Log("Collides when it shouldn't")
	}

	if t.Failed() {
		t.Log("\n" + bo.String())
	}
}

func TestIter(t *testing.T) {
	bo, _ := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	// Returns false when terminated
	t1 := bo.Iter(func(pos Vec2, sq *Square) bool {
		return false
	})
	if t1 {
		t.Log("Iter returns true (should be false) when terminated")
		t.Fail()
	}

	// Gets 5 sometime going across a row
	t2 := bo.IterIn(Vec2{0, 4}, Vec2{3, 5}, func(pos Vec2, sq *Square) bool {
		if IsGiven(*sq) && sq.Area == 5 {
			return false
		}
		return true
	})
	if t2 {
		t.Log("IterIn doesn't find a Given in a row")
		t.Fail()
	}

	// Gets 6 sometime going down a column
	t3 := bo.IterIn(Vec2{3, 0}, Vec2{4, 5}, func(pos Vec2, sq *Square) bool {
		if IsGiven(*sq) && sq.Area == 6 {
			return false
		}
		return true
	})
	if t3 {
		t.Log("IterIn doesn't find a Given in a column")
		t.Fail()
	}

	// IterWhere finds the 2
	t4 := bo.IterWhere(IsGiven, func(pos Vec2, sq *Square) bool {
		if sq.Area == 2 {
			return false
		}
		return true
	})
	if t4 {
		t.Log("IterWhere doesn't find a Given")
		t.Fail()
	}

}

func TestSolve(t *testing.T) {
	bo, _ := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	origStr := bo.String()

	err := bo.Solve()
	if err != nil {
		t.Error("Couldn't find solution to solvable puzzle:", err)
	}

	// Bad excuse for a snapshot test
	correctStr := `      0  1  2  3  4

 0    5  5 05  5  5
 1    3 04  4  4  4
 2   03 02  6  6  6
 3    3  2  6 06  6
 4    5 05  5  5  5`

	if t.Failed() {
		t.Log("Original:\n" + origStr)
		t.Log("Correct:\n" + correctStr)
		t.Log("Actual:\n" + bo.String())
	}
}
