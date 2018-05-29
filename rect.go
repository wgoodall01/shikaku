package shikaku

import "fmt"

type Rect struct {
	// A, B are Vec2s representing the top-left and bottom-right corner of the Rect.
	A, B Vec2

	// Given is the location of the given square which the Rect surrounds.
	Given Vec2
}

func (r Rect) String() string {
	return fmt.Sprintf("%v:%v-%v", r.Given, r.A, r.B)
}

func (r Rect) Size() Vec2 {
	return r.B.Sub(r.A)
}

func (r Rect) Height() int {
	return r.Size()[1]
}

func (r Rect) Width() int {
	return r.Size()[0]
}
