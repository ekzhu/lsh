package lsh

import (
	"fmt"
	"sync"
)

type basicHashTableKey string

type hashTable map[basicHashTableKey]hashTableBucket

type BasicLsh struct {
	*lshParams
	// Number of distinct hashes in the index.
	count int
	// Hash tables.
	tables []hashTable
}

func NewBasicLsh(dim, l, m int, w float64) *BasicLsh {
	tables := make([]hashTable, l)
	for i := range tables {
		tables[i] = make(hashTable)
	}
	return &BasicLsh{
		lshParams: newLshParams(dim, l, m, w),
		count:     0,
		tables:    tables,
	}
}

func (index *BasicLsh) toBasicHashTableKeys(keys []hashTableKey) []basicHashTableKey {
	basicKeys := make([]basicHashTableKey, index.l)
	for i, key := range keys {
		s := ""
		for _, hashVal := range key {
			s += fmt.Sprintf("%.16x", hashVal)
		}
		basicKeys[i] = basicHashTableKey(s)
	}
	return basicKeys
}

// Insert adds a new key to the LSH
func (index *BasicLsh) Insert(point Point, id int) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(point))
	// Insert key into all hash tables
	var wg sync.WaitGroup
	for i := range index.tables {
		hv := hvs[i]
		table := index.tables[i]
		wg.Add(1)
		go func(table hashTable, hv basicHashTableKey) {
			if _, exist := table[hv]; !exist {
				table[hv] = make(hashTableBucket, 0)
			}
			table[hv] = append(table[hv], id)
			wg.Done()
		}(table, hv)
	}
	wg.Wait()
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (index *BasicLsh) Query(q Point, out chan int) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(q))
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
