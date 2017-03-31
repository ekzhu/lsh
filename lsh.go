package lsh

import (
	"math"
	"math/rand"
)

const (
	rand_seed = 1
)

// Key is a way to index into a table.
type hashTableKey []int

// Value is an index into the input dataset.
type hashTableBucket []string

type lshParams struct {
	// Dimensionality of the input data.
	dim int
	// Number of hash tables.
	l int
	// Number of hash functions for each table.
	m int
	// Shared constant for each table.
	w float64

	// Hash function params for each (l, m).
	a [][]Point
	b [][]float64
}

// NewLshParams initializes the LSH settings.
func newLshParams(dim, l, m int, w float64) *lshParams {
	// Initialize hash params.
	a := make([][]Point, l)
	b := make([][]float64, l)
	random := rand.New(rand.NewSource(rand_seed))
	for i := range a {
		a[i] = make([]Point, m)
		b[i] = make([]float64, m)
		for j := range a[i] {
			a[i][j] = make(Point, dim)
			for d := 0; d < dim; d++ {
				a[i][j][d] = random.NormFloat64()
			}
			b[i][j] = random.Float64() * float64(w)
		}
	}
	return &lshParams{
		dim: dim,
		l:   l,
		m:   m,
		a:   a,
		b:   b,
		w:   w,
	}
}

// Hash returns all combined hash values for all hash tables.
func (lsh *lshParams) hash(point Point) []hashTableKey {
	hvs := make([]hashTableKey, lsh.l)
	for i := range hvs {
		s := make(hashTableKey, lsh.m)
		for j := 0; j < lsh.m; j++ {
			hv := (point.Dot(lsh.a[i][j]) + lsh.b[i][j]) / lsh.w
			s[j] = int(math.Floor(hv))
		}
		hvs[i] = s
	}
	return hvs
}
