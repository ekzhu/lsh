package lsh

import (
	"strconv"
	"testing"
)

func Test_NewLshForest(t *testing.T) {
	lsh := NewLshForest(5, 5, 100, 5.0)
	if len(lsh.trees) != 5 {
		t.Error("Lsh init fail")
	}
}

func Test_LshForestInsert(t *testing.T) {
	lsh := NewLshForest(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		lsh.Insert(p, strconv.Itoa(i))
	}
	for _, trees := range lsh.trees {
		if trees.count == 0 {
			t.Error("Insert fail")
		}
	}
}

func Test_LshForestQuery(t *testing.T) {
	lsh := NewLshForest(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	insertedKeys := make([]string, 10)
	for i, p := range points {
		lsh.Insert(p, strconv.Itoa(i))
		insertedKeys[i] = strconv.Itoa(i)
	}

	// Use the inserted points as queries, and
	// verify that we can get back each query itself
	for i, key := range insertedKeys {
		found := false
		for _, foundKey := range lsh.Query(points[i], 5) {
			if foundKey == key {
				found = true
			}
		}
		if !found {
			t.Error("Query fail")
		}
	}
}
