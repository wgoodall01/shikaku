package shikaku

import "fmt"

// Square represents a Shikaku square, part of a Board.
type Square struct {
	// Area of the rectangle which must enclose the Square.
	// For a given square, Area != 0. For a blank, Area <= 0.
	Area int

	// Final is a blank square's parent, if it's known for certain
	// If !nil, then Possible should be nil.
	Final *Square

	// For a blank square, possible values for its parent square.
	Possible []*Square
}

// NewBlank creates a blank Square
func NewBlank() Square {
	return Square{
		Area:     0,
		Possible: make([]*Square, 0, 5),
	}
}

// NewGiven creates a given Square with an area
func NewGiven(area int) Square {
	return Square{
		Area:     area,
		Possible: nil,
	}
}

// String returns a string representation of a Given
func (sq Square) String() string {
	if IsGiven(sq) {
		return fmt.Sprintf("Given(%d)", sq.Area)
	}

	if IsFinal(sq) {
		return fmt.Sprintf("Final(%v)", sq.Final.Area)
	}

	return fmt.Sprintf("Blank%v", sq.Possible)
}

// IsNotFinal returns !IsFinal(sq)
func IsNotFinal(sq Square) bool {
	return !IsFinal(sq)
}

// IsFinal returns true if a square's value is known, and false if it isn't.
func IsFinal(sq Square) bool {
	return sq.Final != nil || IsGiven(sq)
}

// IsAny returns true.
// (for use as predicate, passed to Board.IterWhere)
func IsAny(sq Square) bool {
	return true
}

// IsGiven returns if a square is Given
func IsGiven(sq Square) bool {
	return sq.Area > 0
}
