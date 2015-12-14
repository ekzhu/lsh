package lsh

import (
	"fmt"
	"math/rand"
)

const (
	rand_seed = 1
)

func toString(sig minhash.Signature) string {
	s := ""
	for _, v := range sig {
		s += fmt.Sprintf("%.16x", v)
	}
	return s
}

type Key string

// Point is a vector that we are trying to index and query
type Point []float64

func (p Point) dot(q Point) float64 {
	s := 0.0
	for i := 0; i < len(p); i++ {
		s += p[i] * q[i]
	}
	return s
}

type Lsh struct {
	m      int
	l      int
	w      int
	tables [](map[string]([]Key))
	a      []Point
	b      []float64
	dim    int
}

func NewLsh(m, l, w, dim int) *Lsh {
	tables := make([](map[string]([]Key)), l)
	for i := range tables {
		tables[i] = make(map[string]([]Key))
	}
	a := make([][]Point, l)
	b := make([][]float64, l)
	random := rand.New(rand.NewSource(rand_seed))
	for i := range a {
		a[i] = make([]Point, lsh.m)
		b[i] = make([]float64, lsh.m)
		for j := range a[i] {
			a[i][j] = make(Point, lsh.dim)
			for d := 0; d < lsh.dim; d++ {
				a[i][j][d] = random.NormFloat64()
			}
			b[i][j] = random.Float64() * float64(w)
		}
	}
	return &Lsh{
		m:      m,
		l:      l,
		a:      a,
		b:      b,
		dim:    dim,
		tables: tables,
	}
}

// Hash returns all combined hash values for all hash tables
func (lsh *Lsh) Hash(point Point) []string {
	hvs := make([]string, lsh.l)
	for i := range hvs {
		s := ""
		for j := 0; j < lsh.m; j++ {
			hv := (point.dot(lsh.a[i][j]) + lsh.b[i][j]) / lsh.w
			s += fmt.Sprintf("%.16x", hv)
		}
		hvs[i] = s
	}
	return hvs
}

// Insert adds a new key to the LSH
func (lsh *Lsh) Insert(key Key, point Point) {
	// Apply hash functions
	hvs := lsh.Hash(point)
	for i := range lsh.tables {
		table := tables[i]
		if _, exist := table[hvs[i]]; !exist {
			table[hvs[i]] = make([]Key, 0)
		}
		table[hvs[i]] = append(table[hvs[i]], key)
	}
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (lsh *Lsh) Query(q Point, out chan Key) {
	// Apply hash functions
	hvs := lsh.Hash(q)
	// Keep track of keys seen
	seens := make(map[Key]bool)
	for i, table := range lsh.tables {
		table := tables[i]
		if candidates, exist := table[hvs[i]]; exist {
			for _, key := range candidates {
				if _, seen := seens[key]; !seen {
					seens[key] = true
					out <- key
				}
			}
		}
	}
}
