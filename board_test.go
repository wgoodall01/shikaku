package shikaku

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
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

	cupaloy.SnapshotT(t, bo)

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

	cupaloy.SnapshotT(t, bo.Height(), bo.Width(), bo.Size())

	if t.Failed() {
		t.Log("Failed snapshot\n" + bo.String())
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
	if bo.Collides(Rect{Vec2{0, 0}, Vec2{3, 3}, Vec2{1, 1}}) != true {
		t.Fail()
		t.Log("Doesn't collide when it should")
	}

	// Bottom right 3x3, no collides.
	if bo.Collides(Rect{Vec2{2, 2}, Vec2{5, 5}, Vec2{3, 3}}) != false {
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
	makeTest := func(boString string) {
		t.Run("Board", func(t *testing.T) {
			bo, _ := NewBoardFromString(boString)
			origStr := bo.String()

			err := bo.Solve()
			if err != nil {
				t.Error("Couldn't find solution to solvable puzzle:", err)
			}

			cupaloy.SnapshotT(t, bo.String(), bo.DebugString())

			if t.Failed() {
				t.Log("Original:\n" + origStr)
				t.Log("Actual:\n" + bo.String())
				t.Log("Debug:\n" + bo.DebugString())
			}
		})
	}

	makeTest(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	makeTest(`
		03 -- -- 02 --
		05 -- -- -- --
		02 -- 03 -- 06
		-- -- -- -- -- 
		-- 04 -- -- --
	`)

	makeTest(`
		-- -- -- -- -- -- -- 08
		07 -- -- -- -- -- -- --
		-- -- -- -- -- -- 02 03
		-- 08 -- -- -- 04 -- --
		-- 04 -- -- -- -- 08 --
		-- -- 03 03 -- -- -- --
		06 -- -- -- 02 04 -- 02
		-- -- -- -- -- -- -- --
	`)
}
