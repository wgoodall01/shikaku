package shikaku

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Board Represents a Shikaku board, containing a grid of Squares.
type Board struct {
	Grid [][]Square
}

// Height returns the height of the board.
func (bo *Board) Height() int {
	return len(bo.Grid)
}

// Width returns the width of the board.
func (bo *Board) Width() int {
	return len(bo.Grid[0])
}

// Size returns the size of the board in a Vec2.
func (bo *Board) Size() Vec2 {
	return Vec2{bo.Width(), bo.Height()}
}

// Get returns the square at the given Vec2{x,y}.
// Panics if the coordinates are out of bounds.
func (bo *Board) Get(pos Vec2) *Square {
	x, y := pos[0], pos[1]
	if x >= bo.Width() || x < 0 || y >= bo.Height() || y < 0 {
		panic("Get() coordinates out of bounds")
	}

	return &bo.Grid[y][x]
}

// BoardVisitor is a function called for a certain set of squares in a
// board, returning true to advance or false to stop iterating.
type BoardVisitor = func(pos Vec2, sq *Square) (advance bool)

// Iter calls visitor for each square in the board.
func (bo *Board) Iter(visitor BoardVisitor) bool {
	return bo.IterIn(ORIGIN, bo.Size(), visitor)
}

// IterWhere calls visitor for each square in the board satisfying predicate, stopping if visitor returns false.
func (bo *Board) IterWhere(predicate func(sq Square) bool, visitor BoardVisitor) (uninterrupted bool) {
	return bo.IterIn(ORIGIN, bo.Size(), func(pos Vec2, sq *Square) bool {
		if predicate(*sq) {
			return visitor(pos, sq)
		}
		return true

	})
}

// IterIn valls visitor for each square in the rectangular range from a (inclusive) to b (exclusive).
//
// Preconditions:
//   a[0] <= b[0]
//   a[1] <= b[1]
// Panics otherwise.
func (bo *Board) IterIn(a, b Vec2, visitor BoardVisitor) (uninterrupted bool) {
	if a[0] == b[0] || a[1] == b[1] {
		return true // Can't iterate through an empty rectangle
	} else if a[0] > b[0] || a[1] > b[1] {
		panic("IterIn() passed negative area")
	}

	var pos Vec2
	for pos[0] = a[0]; pos[0] < b[0]; pos[0]++ {
		for pos[1] = a[1]; pos[1] < b[1]; pos[1]++ {
			sq := bo.Get(pos)
			advance := visitor(pos, sq)
			if !advance {
				return false
			}
		}
	}
	return true
}

// NewBoardFromString creates a new Board from a string of the following format:
//
// Each blank as '--'
// Each given as 'xx', e.g. 00, 04, 16, 99, etc.
//
// Each square separated by a space, and each line by a newline (\n).
//
// For example, a 5x5 could look like this:
//  -- -- 05 -- --
//  -- 04 -- -- --
//  03 02 -- -- --
//  -- -- -- 06 --
//  -- 05 -- -- --
func NewBoardFromString(s string) (*Board, error) {
	s = strings.TrimSpace(s)
	b := new(Board)

	// For each line
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		row := []Square{}
		for _, sqStr := range strings.Fields(line) {
			if sqStr == "--" {
				row = append(row, NewBlank())
			} else {
				area, err := strconv.Atoi(sqStr)
				if err != nil {
					return nil, fmt.Errorf("Couldn't parse board: '%s' isn't an int", sqStr)
				}
				row = append(row, NewGiven(area))
			}
		}
		b.Grid = append(b.Grid, row)
	}

	// Sanity check: all the rows are the same length
	rowLen := len(b.Grid[0])
	for _, row := range b.Grid {
		if len(row) != rowLen {
			return nil, errors.New("Couldn't parse board: not all rows are the same length")
		}
	}

	return b, nil
}

// Solve solves the Shikaku puzzle
/*

For each Given:
	for each factor pair (each way around)
		For each possible placement
			add a potential answer to the blank

For each Blank:
	if it has 0 Possibles, abort with error.
	Count those with len(possible) != 1

If count(len(possible) != 1) is zero
	Done
else
	Repeat everything

*/
func (bo *Board) Solve() error {

	// Finalize if only one solution for anything.
	// So count the number of times something's finalized.
	countFinalized := 0

	// For each Given
	bo.IterWhere(IsGiven, func(pos Vec2, giv *Square) bool {

		// Count possible orientations. If there's only 1, finalize it.
		countPossible := 0
		lastPossible := [2]Vec2{{0, 0}, {0, 0}}

		// For each factor pair...
		factors := Factor(giv.Area)
		for _, area := range factors {

			// ...each way around
			for flip := 0; flip <= 1; flip++ {

				// For each possible placement...
				var ofs Vec2 // Offset of top left corner to Given loc
				for ofs[0] = 0; ofs[0] < area[0]; ofs[0]++ {
					for ofs[1] = 0; ofs[1] < area[1]; ofs[1]++ {
						a := pos.Sub(ofs)
						b := a.Add(area)

						contained := (a.In(ORIGIN, bo.Size()) && b.In(ORIGIN, bo.Size().Add(Vec2{1, 1})))

						// ...That doesn't collide, and that fits
						if contained && !bo.Collides(a, b, giv) {
							// incr possible count, set lastPossible
							countPossible++
							lastPossible = [2]Vec2{a, b}

							// Add a Potential for each square in the area.
							bo.IterIn(a, b, func(pos Vec2, potential *Square) bool {
								if potential != giv {
									potential.Possible = append(potential.Possible, giv)
								}
								return true
							})
						}
					}
				}

				// Flip the factor pair, then try again.
				// If it's a square, don't flip it.
				if area[0] != area[1] {
					area[0], area[1] = area[1], area[0]
				} else {
					break
				}
			}

		}

		// If there's only one solution...
		if countPossible == 1 {
			// Finalize that solution.
			a, b := lastPossible[0], lastPossible[1]
			bo.IterIn(a, b, func(pos Vec2, sq *Square) bool {
				if IsNotFinal(*sq) {
					sq.Final = giv
					sq.Possible = sq.Possible[:0]
					countFinalized++
				}
				return true
			})
		}

		return true // keep going
	})

	// For each Blank
	remaining := 0
	valid := bo.IterWhere(IsNotFinal, func(pos Vec2, blank *Square) bool {
		if len(blank.Possible) == 0 {
			return false
		} else if len(blank.Possible) != 1 {
			remaining++
		}
		return true
	})

	if !valid {
		return errors.New("Invalid board, some squares weren't covered")
	}

	if remaining == 0 {
		// Done. Everything's fine.
		return nil
	}

	// Finalize squares with 1 suggestion, add to the count.
	bo.Iter(func(pos Vec2, sq *Square) bool {
		if IsNotFinal(*sq) && !IsGiven(*sq) {
			if len(sq.Possible) == 1 {
				//Make final.
				sq.Final = sq.Possible[0]
				countFinalized++
			}
		}
		return true
	})

	if countFinalized == 0 {
		// No squares finalized, no progress made. Abort to stop infinite loop.
		return errors.New("couldn't solve unambiguously")
	}

	// Truncate lists of Possible
	bo.IterWhere(IsNotFinal, func(pos Vec2, sq *Square) bool {
		sq.Possible = sq.Possible[:0]
		return true
	})

	// Try refining it again.
	return bo.Solve()
}

// StringGiven returns a string representation of the board, in the same format as NewBoardFromString.
func (bo *Board) StringGiven() string {
	var buf bytes.Buffer

	for _, row := range bo.Grid {
		for _, sq := range row {
			if IsGiven(sq) {
				fmt.Fprintf(&buf, "%02d ", sq.Area)
			} else {
				fmt.Fprintf(&buf, "-- ")
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// String returns a string representation of the board.
func (bo *Board) String() string {
	var buf bytes.Buffer

	// write header
	fmt.Fprint(&buf, "    ")
	for i := 0; i < bo.Width(); i++ {
		fmt.Fprintf(&buf, " %2d", i)
	}
	fmt.Fprint(&buf, "\n\n")

	// Write each line
	for i, row := range bo.Grid {
		fmt.Fprintf(&buf, "%2d  ", i)
		for _, sq := range row {
			if IsGiven(sq) {
				fmt.Fprintf(&buf, " %02d", sq.Area)
			} else {
				if IsFinal(sq) {
					fmt.Fprintf(&buf, "  %1d", sq.Final.Area)
				} else {
					fmt.Fprintf(&buf, "   ")
				}

			}
		}
		fmt.Fprint(&buf, "\n")
	}

	return buf.String()
}

// DebugString prints each cell.
func (bo *Board) DebugString() string {
	var buf bytes.Buffer

	bo.Iter(func(loc Vec2, sq *Square) bool {
		fmt.Fprintf(&buf, "%v: %v\n", loc, sq)
		return true
	})

	return buf.String()
}

// Collides determines if the rectangle bounded by a (inclusively) and
// b (exclusively) overlaps with any Concrete squares other than ignore.
//
// Preconditions:
//   a[0] <= b[0]
//   a[1] <= b[1]
// Panics otherwise.
func (bo *Board) Collides(a, b Vec2, ignore *Square) bool {
	return !bo.IterIn(a, b, func(pos Vec2, sq *Square) bool {
		return (IsGiven(*sq) && sq == ignore) || (IsFinal(*sq) && sq.Final == ignore) || IsNotFinal(*sq)
	})
}
