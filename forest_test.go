package lsh

import "testing"

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
		lsh.Insert(p, i)
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
	insertedKeys := make([]int, 10)
	for i, p := range points {
		lsh.Insert(p, i)
		insertedKeys[i] = i
	}

	// Use the inserted points as queries, and
	// verify that we can get back each query itself
	for i, key := range insertedKeys {
		done := make(chan struct{})
		found := false
		for foundKey := range lsh.Query(points[i], done) {
			if foundKey == key {
				found = true
			}
		}
		if !found {
			t.Error("Query fail")
		}
		close(done)
	}
}

func Test_LshForestQueryKnn(t *testing.T) {
	lsh := NewLshForest(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	insertedKeys := make([]int, 10)
	for i, p := range points {
		lsh.Insert(p, i)
		insertedKeys[i] = i
	}
	done := make(chan struct{})
	// Query a single point, should obtain back the entire index.
	actual := make([]int, 0)
	for key := range lsh.QueryKnn(points[0], 10, done) {
		actual = append(actual, key)
	}
	close(done)
	// Use the inserted points as queries, and
	// verify that we can get back each query itself
	for _, key := range insertedKeys {
		found := false
		for foundKey := range actual {
			if foundKey == key {
				found = true
			}
		}
		if !found {
			t.Error("Query fail")
		}
	}
}
