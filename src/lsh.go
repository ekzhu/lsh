package lsh

import (
	"fmt"
	"math/rand"
)

const (
	rand_seed = 1
)

type Key string

type Lsh struct {
	m      int
	l      int
	w      float64
	tables [](map[string]([]Key))
	a      [][]Point
	b      [][]float64
	dim    int
}

func NewLsh(m, l, dim int, w float64) *Lsh {
	tables := make([](map[string]([]Key)), l)
	for i := range tables {
		tables[i] = make(map[string]([]Key))
	}
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
	return &Lsh{
		m:      m,
		l:      l,
		a:      a,
		b:      b,
		w:      w,
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
	// Insert key into all hash tables
	for i, table := range lsh.tables {
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
