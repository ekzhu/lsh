package lsh

import "testing"

func Test_NewMultiprobeLsh(t *testing.T) {
	lsh := NewMultiprobeLsh(100, 5, 2, 5.0, 64)
	if len(lsh.SimpleIndex.tables) != 5 {
		t.Error("Lsh init fail")
	}
	t.Logf("Scores %v", lsh.scores)
	t.Logf("Perturbation sets: %v", lsh.perturbSets)
	for i, v := range lsh.perturbSets {
		t.Logf("Set: %d, Score: %f, Set contents: %v", i, lsh.getScore(&v), v)
	}

}

/*
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
		result := make(chan int)
		go func() {
			lsh.Query(points[i], result)
			close(result)
		}()
		found := false
		for foundKey := range result {
			if foundKey == key {
				found = true
			}
		}
		if !found {
			t.Error("Query fail")
		}
	}
}

func Test_LshForestParallelQuery(t *testing.T) {
	lsh := NewLshForest(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		lsh.Insert(p, i)
	}
	// Run multiple queries in parallel
	// and writing candidates to the same output
	queries := randomPoints(10, 10, 32.0)
	in := make(chan Point)
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(5)
	// Input thread
	go func() {
		for _, q := range queries {
			in <- q
		}
		close(in)
	}()
	// Worker threads
	for w := 0; w < 5; w++ {
		go func() {
			for q := range in {
				lsh.Query(q, out)
			}
			wg.Done()
		}()
	}
	// Waiter thread
	go func() {
		wg.Wait()
		close(out)
	}()
	// Main thread collecting outputs
	for _ = range out {
	}
}

func Test_LshForestQueryK(t *testing.T) {
	lsh := NewLshForest(100, 5, 5, 5.0)
	points := randomPoints(10, 100, 32.0)
	insertedKeys := make([]int, 10)
	for i, p := range points {
		lsh.Insert(p, i)
		insertedKeys[i] = i
	}

	// Query a single point, should obtain back the entire index.
	result := make(chan int)
	go func() {
		lsh.QueryK(points[0], 10, result)
		close(result)
	}()
	actual := make([]int, 0)
	for key := range result {
		actual = append(actual, key)
	}

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
}*/
