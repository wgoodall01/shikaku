package shikaku

import (
	"fmt"
	"math"
)

// Vec2 represents a 2-dimensional integer vector
type Vec2 [2]int

// ORIGIN is Vec2{0,0}
var ORIGIN = Vec2{0, 0}

// Add returns a Vec2 representing the sum of v and x.
func (v Vec2) Add(x Vec2) Vec2 {
	return Vec2{v[0] + x[0], v[1] + x[1]}
}

// Sub returns a Vec2 representing the difference of v and x.
func (v Vec2) Sub(x Vec2) Vec2 {
	return Vec2{v[0] - x[0], v[1] - x[1]}
}

// In returns true if lo <= v < hi in both dimensions, false otherwise.
func (v Vec2) In(lo, hi Vec2) bool {
	return lo[0] <= v[0] && lo[1] <= v[1] && v[0] < hi[0] && v[1] < hi[1]
}

// Transpose the vector. Flips v[0] and v[1].
func (v Vec2) Transpose() Vec2 {
	return Vec2{v[1], v[0]}
}

// String returns a [x,y] representation of the Vec2.
func (v Vec2) String() string {
	return fmt.Sprintf("[%d,%d]", v[0], v[1])
}

// Factor finds all the integer factor pairs of x, sorted by the smallest factor.
func Factor(x int) []Vec2 {
	factors := []Vec2{}

	// Max of any one factor is sqrt(x)
	max := int(math.Floor(math.Sqrt(float64(x))))

	for div := 1; div <= max; div++ {
		mod := x % div
		if mod == 0 {
			factors = append(factors, Vec2{div, x / div})
		}
	}

	return factors
}
