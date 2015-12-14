package lsh

import (
	"fmt"
)

type SimpleIndexKey string

// A table in the simple index is a lookup from a TableKey to a value.
type Table map[SimpleIndexKey]Value

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
		LshSettings: NewLshSettings(dim, l, m, w),
		count:       0,
		tables:      tables,
	}
}

func (index *SimpleIndex) toSimpleKeys(keys []TableKey) []SimpleIndexKey {
	simpleKeys := make([]SimpleIndexKey, index.l)
	for i, key := range keys {
		s := ""
		for _, hashVal := range key {
			s += fmt.Sprintf("%.16x", hashVal)
		}
		simpleKeys[i] = SimpleIndexKey(s)
	}
	return simpleKeys
}

// Insert adds a new key to the LSH
func (index *SimpleIndex) Insert(point Point, id int) {
	// Apply hash functions
	hvs := index.toSimpleKeys(index.Hash(point))
	// Insert key into all hash tables
	for i, table := range index.tables {
		if _, exist := table[hvs[i]]; !exist {
			table[hvs[i]] = make(Value, 0)
		}
		table[hvs[i]] = append(table[hvs[i]], id)
	}
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (index *SimpleIndex) Query(q Point, out chan int) {
	// Apply hash functions
	hvs := index.toSimpleKeys(index.Hash(q))
	// Keep track of keys seen
	seens := make(map[int]bool)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				if _, seen := seens[id]; !seen {
					seens[id] = true
					out <- id
				}
			}
		}
	}
}
