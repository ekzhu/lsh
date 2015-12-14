package lsh

import (
	"fmt"
	"math/rand"
)

const (
	rand_seed = 1
)

// Key is a way to index into a table.
type TableKey string

// Value is an index into the input dataset.
type Value []int

type LshSettings struct {
	// Dimensionality of the input data.
	dim int
	// Number of tables.
	l int
	// Number of hash functions for each table.
	m int
	// Shared constant for each table.
	w float64

	// Hash function params for each (l, m).
	a [][]Point
	b [][]float64
}

// NewLshSettings initializes the LSH settings.
func NewLshSettings(dim, l, m int, w float64) *LshSettings {
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
	return &LshSettings{
		dim: dim,
		l:   l,
		m:   m,
		a:   a,
		b:   b,
		w:   w,
	}
}

// Hash returns all combined hash values for all hash tables.
func (lsh *LshSettings) Hash(point Point) []TableKey {
	hvs := make([]TableKey, lsh.l)
	for i := range hvs {
		s := ""
		for j := 0; j < lsh.m; j++ {
			hv := (point.dot(lsh.a[i][j]) + lsh.b[i][j]) / lsh.w
			s += fmt.Sprintf("%.16x", hv)
		}
		hvs[i] = TableKey(s)
	}
	return hvs
}
