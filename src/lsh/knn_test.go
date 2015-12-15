package lsh

import (
	"container/heap"
	"sort"
	"testing"
)

func Test_KHeap(t *testing.T) {
	points := randomPoints(10, 10, 100.0)
	k := 5
	h := NewKHeap(k)
	heap.Init(h)
	q := points[0]
	distances := make([]float64, len(points))
	for i := range points {
		distances[i] = points[i].L2(q)
		c := Candidate{i, distances[i]}
		heap.Push(h, c)
		t.Log(c)
		t.Log(h.candidates)
	}
	if len(h.candidates) != k {
		t.Error("Heap failed to maintain correct number of items")
	}
	sort.Float64s(distances)
	topK := make([]float64, k)
	for i := 0; i < k; i++ {
		c := heap.Pop(h).(Candidate)
		topK[i] = c.distance
	}
	for i := range topK {
		if topK[i] != distances[k-1-i] {
			t.Errorf("Expected order (reverse it) <%s>\nActual order <%s>\n",
				distances[:k], topK)
		}
	}
}

func Test_Knn(t *testing.T) {
	k := 5
	points := randomPoints(20, 10, 100.0)
	ids := make([]int, len(points))
	for i := range points {
		ids[i] = i
	}
	q := points[0]
	// Build ground truth
	distances := make([]float64, len(points))
	for i := range points {
		distances[i] = points[i].L2(q)
	}
	sort.Float64s(distances)
	t.Log("Ground truth distances", distances[:k])
	// Test Knn query
	knn := NewKnn(points, ids)
	out := make(chan int)
	go func() {
		knn.Query(q, k, out)
		close(out)
	}()
	for id := range out {
		// get the point
		p := points[id]
		// get the distance
		d := p.L2(q)
		t.Log(d)
		// check if the distance is indeed within the top k ground truth
		found := false
		for i := 0; i < k; i++ {
			if d == distances[i] {
				found = true
			}
		}
		if !found {
			t.Error("Query did not find top k ground truth distance")
		}
	}
}
