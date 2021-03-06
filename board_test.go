package shikaku

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

var testBoards = []string{
	`-- -- 05 -- --
	 -- 04 -- -- --
	 03 02 -- -- --
	 -- -- -- 06 --
	 -- 05 -- -- --`,

	`03 -- -- 02 --
	 05 -- -- -- --
	 02 -- 03 -- 06
	 -- -- -- -- -- 
	 -- 04 -- -- --`,

	`-- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- --
	 -- -- 64 -- -- -- -- --
	 -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- --`,

	`-- -- -- -- -- -- -- 08
	 07 -- -- -- -- -- -- --
	 -- -- -- -- -- -- 02 03
	 -- 08 -- -- -- 04 -- --
	 -- 04 -- -- -- -- 08 --
	 -- -- 03 03 -- -- -- --
	 06 -- -- -- 02 04 -- 02
	 -- -- -- -- -- -- -- --`,

	`-- -- -- 04 17 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 03 -- -- 02
	 02 -- -- 22 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 02 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 04 -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- -- -- 24 02 -- -- -- -- -- 02 -- 04 -- -- 
	 -- -- -- -- -- -- -- -- -- 13 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- 09 -- -- -- -- -- 21 -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- 02 -- 02 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- 20 -- -- -- -- -- -- -- -- -- -- -- 08 -- -- -- 80 -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 95 -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- 12 -- 21 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- 03 -- -- -- -- 03 -- -- -- 03 -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- 02 -- -- 03 -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- 03 -- 02 -- 02 -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- 32 -- -- -- 04 -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- 04 -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- 03 -- 06 -- -- -- -- -- 03 -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- 15 -- -- 03 -- 09 -- -- -- -- -- -- -- -- 02 -- -- -- -- -- --
	 -- 44 -- -- -- -- -- -- 04 -- -- -- 06 -- -- -- 02 -- -- 02 -- -- -- -- --
	 -- -- -- -- -- -- -- 26 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 03 02 -- -- -- -- --
	 -- -- -- -- 09 -- -- -- -- -- -- -- -- 05 -- -- -- 03 -- -- -- -- -- -- --
	 -- -- -- -- 03 -- -- 06 -- -- -- 28 -- -- -- -- -- -- -- -- -- -- -- -- --
	 -- -- -- -- -- -- -- -- -- 09 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --`,

	`-- 02 02 -- -- -- 02
	 -- 03 -- -- 06 -- --
	 06 -- 05 02 -- -- --
	 -- -- -- -- 02 -- --
	 -- 03 -- -- 04 -- 05
	 -- -- -- -- -- 04 --
	 -- -- 03 -- -- -- --`,

	`12 -- -- -- -- -- -- --
	 -- -- -- -- -- 04 -- --
	 -- -- -- -- 06 -- -- --
	 -- -- 08 -- -- -- 08 --
	 -- -- -- -- -- -- -- 03
	 -- 09 -- 01 -- 04 -- --
	 -- -- -- 04 -- -- -- --
	 -- -- -- -- 03 -- -- 02`,
}

var testBadBoards = []string{
	`12 01 -- -- -- -- -- --
	 -- -- -- -- -- 04 -- --
	 -- -- -- -- 06 -- -- --
	 -- -- 08 -- -- -- 08 --
	 -- -- -- -- -- -- -- 03
	 -- 09 -- 01 -- 04 -- --
	 -- -- -- 04 -- -- -- --
	 -- -- -- -- 03 -- -- 02`,
}

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
	for _, boString := range testBoards {
		t.Run("Board", func(t *testing.T) {
			bo, err := NewBoardFromString(boString)

			if err != nil {
				t.Fatal("Failed to parse valid board")
			}

			cupaloy.SnapshotT(t, bo)

			if t.Failed() {
				t.Log("\n" + bo.String())
			}
		})
	}
}

func TestDimensions(t *testing.T) {
	for _, boString := range testBoards {
		t.Run("Board", func(t *testing.T) {
			bo, _ := NewBoardFromString(boString)
			cupaloy.SnapshotT(t, bo.Height(), bo.Width(), bo.Size())

			if t.Failed() {
				t.Log("Failed snapshot\n" + bo.String())
			}
		})

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
	for _, boString := range testBoards {
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
}

func TestBadSolve(t *testing.T) {
	for _, boString := range testBadBoards {
		t.Run("Board", func(t *testing.T) {
			bo, _ := NewBoardFromString(boString)
			origStr := bo.String()

			err := bo.Solve()
			t.Log("Error:", err)
			if err == nil {
				t.Error("Should have failed while solving bad puzzle")
			}

			cupaloy.SnapshotT(t, err, bo.String(), bo.DebugString())

			if t.Failed() {
				t.Log("Original:\n" + origStr)
				t.Log("Actual:\n" + bo.String())
				t.Log("Debug:\n" + bo.DebugString())
			}
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	for _, boString := range testBoards {
		b.Run("Board", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bo, _ := NewBoardFromString(boString)
				err := bo.Solve()
				if err != nil {
					b.Fatalf("Solve failed, see TestSolve for details")
				}
			}
		})
	}
}
