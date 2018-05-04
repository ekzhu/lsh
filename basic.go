package lsh

import (
	"fmt"
	"sync"
)

type basicHashTableKey string

type hashTable map[basicHashTableKey]hashTableBucket

// BasicLsh implements the original LSH algorithm for L2 distance.
type BasicLsh struct {
	*lshParams
	// Hash tables.
	tables []hashTable
}

// NewBasicLsh creates a basic LSH for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, w is the slot size for the
// family of LSH functions.
func NewBasicLsh(dim, l, m int, w float64) *BasicLsh {
	tables := make([]hashTable, l)
	for i := range tables {
		tables[i] = make(hashTable)
	}
	return &BasicLsh{
		lshParams: newLshParams(dim, l, m, w),
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

// Insert adds a new data point to the LSH.
// id is the unique identifier for the data point.
func (index *BasicLsh) Insert(point Point, id string) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(point))
	// Insert key into all hash tables
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for i := range index.tables {
		hv := hvs[i]
		table := index.tables[i]
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

// Query finds the ids of approximate nearest neighbour candidates,
// in un-sorted order, given the query point,
func (index *BasicLsh) Query(q Point) []string {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(q))
	// Keep track of keys seen
	seen := make(map[string]bool)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				if _, exist := seen[id]; exist {
					continue
				}
				seen[id] = true
			}
		}
	}
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// Delete removes a new data point to the LSH.
// id is the unique identifier for the data point.
func (index *BasicLsh) Delete(id string) {
	// Delete key from all hash tables
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for i := range index.tables {
		table := index.tables[i]
		go func(table hashTable) {
			for tableIndex, bucket := range table {
				for index, identifier := range bucket {
					if id == identifier {
						table[tableIndex] = remove(bucket, index)
						if len(table[tableIndex]) == 0 {
							delete(table, tableIndex)
						}
					}
				}
			}
			wg.Done()
		}(table)
	}
	wg.Wait()
}

func remove(original []string, index int) []string {
	original[index] = original[len(original)-1]
	original = original[:len(original)-1]
	return original
}
