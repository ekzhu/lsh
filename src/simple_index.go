package lsh

import (
	"fmt"
)

// A table in the simple index is a lookup from a TableKey to a value.
type Table map[TableKey]Value

type SimpleIndex struct {
	*LshSettings
	// Number of distinct hashes in the index.
	count int
	// Hash tables.
	tables []Table
}

func NewSimpleLsh(dim, l, m int, w float64) *SimpleIndex {
	tables := make([]Table, l)
	for i := range tables {
		tables[i] = make(Table)
	}
	return &SimpleIndex{
		LshSettings: NewLshSettings(dim, m, l, w),
		m:           m,
		l:           l,
		a:           a,
		b:           b,
		w:           w,
		dim:         dim,
		tables:      tables,
	}
}

// Insert adds a new key to the LSH
func (index *SimpleIndex) Insert(key Key, point Point) {
	// Apply hash functions
	hvs := index.Hash(point)
	// Insert key into all hash tables
	for i, table := range index.tables {
		if _, exist := table[hvs[i]]; !exist {
			table[hvs[i]] = make([]Key, 0)
		}
		table[hvs[i]] = append(table[hvs[i]], key)
	}
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (index *SimpleIndex) Query(q Point, out chan Key) {
	// Apply hash functions
	hvs := index.Hash(q)
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
