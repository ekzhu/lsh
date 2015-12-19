package lsh

import (
	"container/heap"
	"time"
)

type Candidate struct {
	id       int
	distance float64
}

type KHeap struct {
	k          int
	candidates []Candidate
}

func (h KHeap) Len() int {
	return len(h.candidates)
}

func (h KHeap) Less(i, j int) bool {
	// We want to pop out the candidate with largest distance
	// so we use greater than here
	return h.candidates[i].distance > h.candidates[j].distance
}

func (h KHeap) Swap(i, j int) {
	h.candidates[i], h.candidates[j] = h.candidates[j], h.candidates[i]
}

func (h *KHeap) Push(x interface{}) {
	c := x.(Candidate)
	if len(h.candidates) < h.k {
		h.candidates = append(h.candidates, c)
		return
	}
	// Check if we can still push to the top-k heap when it is full
	if h.candidates[0].distance > c.distance {
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
	h := make([]Candidate, 0)
	return &KHeap{k, h}
}

type Knn struct {
	data []Point
	ids  []int
}

func NewKnn(data []Point, ids []int) *Knn {
	if len(data) != len(ids) {
		panic("Mismatch between size of data and ids")
	}
	return &Knn{data, ids}
}

// Query outputs the top-k closest points from the query point
// to the chanel out. The sequence of output is NOT sorted by
// distance.
func (knn *Knn) Query(q Point, k int, out chan int) {
	kheap := NewKHeap(k)
	heap.Init(kheap)
	for i, p := range knn.data {
		d := p.L2(q)
		heap.Push(kheap, Candidate{knn.ids[i], d})
	}
	for i := range kheap.candidates {
		out <- kheap.candidates[i].id
	}
}

// RunKnn executes the KNN experiment
func RunKnn(datafile, output string, k, nQuery, nWorker int, parser *PointParser) {
	// Load data
	nData := CountPoint(datafile, parser.ByteLen)
	iter := NewDataPointIterator(datafile, parser)
	data := make([]Point, nData)
	ids := make([]int, nData)
	for i := 0; i < nData; i++ {
		p, err := iter.Next()
		if err != nil {
			panic(err.Error())
		}
		data[i] = p.Point
		ids[i] = p.Id
	}

	// Run Knn
	knn := NewKnn(data, ids)
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan int)
		go func() {
			knn.Query(q.Point, k, out)
			close(out)
		}()
		r := make([]int, 0)
		for i := range out {
			r = append(r, i)
		}
		dur := time.Since(start)
		ns := make(Neighbours, len(r))
		for i := range r {
			ns[i] = Neighbour{
				Id:       r[i],
				Distance: q.Point.L2(data[i]),
			}
		}
		return QueryResult{
			QueryId:    q.Id,
			Neighbours: ns,
			Time:       float64(dur) / float64(time.Millisecond),
		}
	}
	// Select queries
	queryIds := SelectQueries(nData, nQuery)
	iter = NewQueryPointIterator(datafile, parser, queryIds)
	// Run queries in parallel
	results := ParallelQueryIndex(iter, queryFunc, nWorker)
	DumpJson(output, results)
}

// RunKnn executes the KNN experiment
func RunKnnSampleAllPair(datafile, output string, nSample, nWorker int, parser *PointParser) {
	// Load data
	nData := CountPoint(datafile, parser.ByteLen)
	pointIds := SelectQueries(nData, nSample)
	iter := NewQueryPointIterator(datafile, parser, pointIds)
	data := make([]Point, nSample)
	ids := make([]int, nSample)
	for i := 0; i < nSample; i++ {
		p, err := iter.Next()
		if err != nil {
			panic(err.Error())
		}
		data[i] = p.Point
		ids[i] = p.Id
	}

	// Run Knn
	knn := NewKnn(data, ids)
	queryFunc := func(q DataPoint) QueryResult {
		start := time.Now()
		out := make(chan int)
		go func() {
			knn.Query(q.Point, nSample, out)
			close(out)
		}()
		r := make([]int, 0)
		for i := range out {
			r = append(r, i)
		}
		dur := time.Since(start)
		ns := make(Neighbours, len(r))
		for i := range r {
			ns[i] = Neighbour{
				Id:       r[i],
				Distance: q.Point.L2(data[i]),
			}
		}
		return QueryResult{
			QueryId:    q.Id,
			Neighbours: ns,
			Time:       float64(dur) / float64(time.Millisecond),
		}
	}
	// Select queries
	queryIds := SelectQueries(nData, nSample)
	iter = NewQueryPointIterator(datafile, parser, queryIds)
	// Run queries in parallel
	results := ParallelQueryIndex(iter, queryFunc, nWorker)
	DumpJson(output, results)
}
