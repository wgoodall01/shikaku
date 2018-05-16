package shikaku

import "math"

// Point represents coordinates into a Board.
type Point struct {
	X int
	Y int
}

// Pair represents a pair of two integers, such as factors of another integer.
type Pair [2]int

// Factor finds all the integer factor pairs of x, sorted by the smallest factor.
func Factor(x int) []Pair {
	factors := []Pair{}

	// Max of any one factor is sqrt(x)
	max := int(math.Floor(math.Sqrt(float64(x))))

	for div := 1; div <= max; div++ {
		mod := x % div
		if mod == 0 {
			factors = append(factors, Pair{div, x / div})
		}
	}

	return factors
}
