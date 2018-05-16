package shikaku

// Square represents a Shikaku square, part of a Board.
type Square interface {
}

// Given represents a given value, part of the original puzzle.
//
// Given.Area is the area of a rectangle which must enclose the Given.
type Given struct {
	Area int
}

// Blank represents an empty square, which could possibly belong to any number of Givens.
type Blank struct {
	Possible []*Given
}
