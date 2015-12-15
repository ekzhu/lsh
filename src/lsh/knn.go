package lsh

import "container/heap"

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
	return &KHeap{k, make([]Candidate, 0)}
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
	for i, p := range knn.data {
		d := p.L2(q)
		heap.Push(kheap, Candidate{knn.ids[i], d})
	}
	for i := range kheap.candidates {
		out <- kheap.candidates[i].id
	}
}
