package shikaku

import "fmt"

type Rect struct {
	// A, B are Vec2s representing the top-left and bottom-right corner of the Rect.
	A, B Vec2

	// Given is the given square which the Rect surrounds.
	Given *Square
}

// NewRect creates a Rect around giv through points a and b.
func NewRect(bo *Board, a, b Vec2, giv *Square) Rect {
	return Rect{
		A:     a,
		B:     b,
		Given: giv,
	}
}

func (r Rect) String() string {
	return fmt.Sprintf("%v-%v-%v", r.Given.Area, r.A, r.B)
}
