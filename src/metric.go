package lsh

import "math"

// Point is a vector that we are trying to index and query
type Point []float64

// dot returns the dot product of two Point vectors
func (p Point) dot(q Point) float64 {
	s := 0.0
	for i := 0; i < len(p); i++ {
		s += p[i] * q[i]
	}
	return s
}

// l2 returns the L2 distance of two Point vectors
func (p Point) l2(q Point) float64 {
	s := 0.0
	for i := 0; i < len(p); i++ {
		d := p[i] - q[i]
		s += d * d
	}
	return math.Sqrt(s)
}
