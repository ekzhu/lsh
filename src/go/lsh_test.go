package lsh

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

// randomPoints returns a slice of point vectors,
// each element of every point vector is drawn from a uniform
// distribution over [0, max)
func randomPoints(n, dim int, max float64) []Point {
	random := rand.New(rand.NewSource(1))
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i] = make(Point, dim)
		for d := 0; d < dim; d++ {
			points[i][d] = random.Float64() * max
		}
	}
	return points
}

func Test_NewLsh(t *testing.T) {
	lsh := NewLsh(5, 5, 100, 5.0)
	if len(lsh.tables) != 5 {
		t.Error("Lsh init fail")
	}
}

func Test_Insert(t *testing.T) {
	lsh := NewLsh(5, 5, 100, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		key := Key(fmt.Sprintf("%d", i))
		lsh.Insert(key, p)
	}
	for _, table := range lsh.tables {
		if len(table) == 0 {
			t.Error("Insert fail")
		}
	}
}

func Test_Query(t *testing.T) {
	lsh := NewLsh(5, 5, 100, 5.0)
	points := randomPoints(10, 100, 32.0)
	insertedKeys := make([]Key, 10)
	for i, p := range points {
		key := Key(fmt.Sprintf("%d", i))
		lsh.Insert(key, p)
		insertedKeys[i] = key
	}
	// Use the inserted points as queries, and
	// verify that we can get back each query itself
	for i, key := range insertedKeys {
		result := make(chan Key)
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

func Test_ParallelQuery(t *testing.T) {
	lsh := NewLsh(5, 5, 100, 5.0)
	points := randomPoints(10, 100, 32.0)
	for i, p := range points {
		key := Key(fmt.Sprintf("%d", i))
		lsh.Insert(key, p)
	}
	// Run multiple queries in parallel
	// and writing candidates to the same output
	queries := randomPoints(10, 10, 32.0)
	in := make(chan Point)
	out := make(chan Key)
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
