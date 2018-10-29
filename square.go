package shikaku

import "fmt"

// Square represents a Shikaku square, part of a Board.
type Square struct {
	// Area of the rectangle which must enclose the Square.
	// For a given square, Area != 0. For a blank, Area <= 0.
	Area int

	// Final is a blank square's parent rectangle, if it's known for certain
	// If !nil, then Possible should be nil.
	// If set to this square, a Given has been finalized.
	Final Rect

	// For a blank square, possible values for its parent Rects
	Possible []Rect
}

// NewBlank creates a blank Square
func NewBlank() Square {
	return Square{
		Area:     0,
		Possible: make([]Rect, 0, 5),
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
	str := ""
	if IsGiven(sq) {
		str += fmt.Sprintf("Given(%d) ", sq.Area)
	}

	if IsFinal(sq) {
		str += fmt.Sprintf("Final(%v) ", sq.Final)
	}

	if len(sq.Possible) > 0 {
		str += fmt.Sprintf("Blank%v ", sq.Possible)
	}

	return str[:len(str)-1]
}

// AddPossible adds a given Rect to the list of possiblities, ignoring duplicates.
// Returns true when a unique possibility is added.
func (sq *Square) AddPossible(r Rect) bool {
	// See if it's already in the list of possibles
	for _, p := range sq.Possible {
		if p == r {
			return false
		}
	}

	// Add it to the list
	sq.Possible = append(sq.Possible, r)
	return true
}

// IsNotFinal returns !IsFinal(sq)
func IsNotFinal(sq Square) bool {
	return !IsFinal(sq)
}

// IsFinal returns true if a square's value is known, and false if it isn't.
func IsFinal(sq Square) bool {
	return sq.Final != Rect{} || IsGiven(sq)
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

// IsUnsolvedGiven returns if a square is Given and not finalized
func IsUnsolvedGiven(sq Square) bool {
	return IsGiven(sq) && (sq.Final == Rect{})
}
