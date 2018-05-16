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
				row = append(row, Blank{})
			} else {
				area, err := strconv.Atoi(sqStr)
				if err != nil {
					return nil, fmt.Errorf("Couldn't parse board: '%s' isn't an int", sqStr)
				}
				row = append(row, Given{Area: area})
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
func (b *Board) Solve() {
	// TODO: implement
}

// Reset removes all Possibles from the board.
func (b *Board) Reset() {
	for _, row := range b.Grid {
		for _, sq := range row {
			switch v := sq.(type) {
			case Blank:
				v.Possible = nil
			}
		}
	}
}

// String returns a string representation of the board, in the same format as NewBoardFromString.
func (b *Board) String() string {
	var buf bytes.Buffer

	for _, row := range b.Grid {
		for _, sq := range row {
			switch v := sq.(type) {
			case Given:
				buf.WriteString(fmt.Sprintf("%02d", v.Area))
			case Blank:
				buf.WriteString("--")
			default:
				panic("Encountered a strange square type")
			}
			buf.WriteString(" ")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// Collides determines if the rectangle bounded (inclusively) by nw and se overlaps with any
func (b *Board) Collides(nw, se Point) bool {
	//TODO implement this
	return false
}
