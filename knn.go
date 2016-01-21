package lsh

import (
	"container/heap"
	"sort"
	"time"
)

type KHeap struct {
	k          int
	candidates []Neighbour
}

func (h KHeap) Len() int {
	return len(h.candidates)
}

func (h KHeap) Less(i, j int) bool {
	// We want to pop out the candidate with largest distance
	// so we use greater than here
	return h.candidates[i].Distance > h.candidates[j].Distance
}

func (h KHeap) Swap(i, j int) {
	h.candidates[i], h.candidates[j] = h.candidates[j], h.candidates[i]
}

func (h *KHeap) Push(x interface{}) {
	c := x.(Neighbour)
	if len(h.candidates) < h.k {
		h.candidates = append(h.candidates, c)
		return
	}
	// Check if we can still push to the top-k heap when it is full
	if h.candidates[0].Distance > c.Distance {
		heap.Pop(h)
		heap.Push(h, c)
	}
	// sliently return when push is not possible
	return
}

func (h *KHeap) Pop() interface{} {
	k := len(h.candidates)
	x := h.candidates[k-1]
	h.candidates = h.candidates[0 : k-1]
	return x
}

func NewKHeap(k int) *KHeap {
	h := make([]Neighbour, 0)
	return &KHeap{k, h}
}

type Knn struct {
	points []DataPoint
}

func NewKnn(points []DataPoint) *Knn {
	return &Knn{points}
}

// Query outputs the top-k closest points from the query point
// to the chanel out. The sequence of output is NOT sorted by
// distance.
func (knn *Knn) Query(q Point, k int, out chan Neighbour) {
	kheap := NewKHeap(k)
	heap.Init(kheap)
	for _, p := range knn.points {
		d := p.Point.L2(q)
		heap.Push(kheap, Neighbour{p.Id, d})
	}
	for i := range kheap.candidates {
		out <- kheap.candidates[i]
	}
}

// RunKnn executes the KNN experiment
func RunKnn(data []DataPoint, queries []DataPoint,
	output string, k, nWorker int) {
	knn := NewKnn(data)
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan Neighbour)
		go func() {
			knn.Query(q.Point, k, out)
			close(out)
		}()
		ns := make(Neighbours, 0)
		for i := range out {
			ns = append(ns, i)
		}
		dur := time.Since(start)
		sort.Sort(ns)
		return QueryResult{
			QueryId:    q.Id,
			Neighbours: ns,
			Time:       float64(dur) / float64(time.Millisecond),
		}
	}
	results := ParallelQueryIndex(queries, queryFunc, nWorker)
	DumpJson(output, results)
}

// RunKnn executes the KNN experiment
func RunKnnSampleAllPair(data []DataPoint, output string, nWorker int) {
	knn := NewKnn(data)
	nSample := len(data)
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan Neighbour)
		go func() {
			knn.Query(q.Point, nSample, out)
			close(out)
		}()
		ns := make(Neighbours, 0)
		for i := range out {
			ns = append(ns, i)
		}
		dur := time.Since(start)
		sort.Sort(ns)
		return QueryResult{
			QueryId:    q.Id,
			Neighbours: ns,
			Time:       float64(dur) / float64(time.Millisecond),
		}
	}
	results := ParallelQueryIndex(data, queryFunc, nWorker)
	DumpJson(output, results)
}
