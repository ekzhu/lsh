package lsh

import (
	"math/rand"
)

// randomPoints returns a slice of point vectors,
// each element of every point vector is drawn from a uniform
// distribution over [0, max)
func randomPoints(n, dim int, max float64) []Point {
	random := rand.New(rand.NewSource(1))
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i] = make(Point, dim)
		for d := 0; d < dim; d++ {
			points[i][d] = random.Float64() * max
		}
	}
	return points
}
