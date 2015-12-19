package lsh

import (
	"container/heap"
)

// perturbSetHeap is a min-heap of perturbSetPairs.
type perturbSetHeap []perturbSetPair

type perturbSet map[int]bool

// A pair of perturbation set and its score.
type perturbSetPair struct {
	ps    perturbSet
	score float64
}

func (h perturbSetHeap) Len() int           { return len(h) }
func (h perturbSetHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h perturbSetHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *perturbSetHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(perturbSetPair))
}

func (h *perturbSetHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MultiprobeIndex struct {
	*SimpleIndex
	// A list of perturbations that will be used for lookups.
	probeSeq []TableKey

	// The scores of perturbation values.
	scores []float64

	perturbSets []perturbSet
}

func NewMultiprobeLsh(dim, l, m int, w float64) *MultiprobeIndex {
	index := &MultiprobeIndex{
		SimpleIndex: NewSimpleLsh(dim, l, m, w),
	}
	index.initProbeSequence()
	return index
}

func (index *MultiprobeIndex) initProbeSequence() {
	m := index.SimpleIndex.LshSettings.m
	index.scores = make([]float64, 2*m)
	// Use j's starting from 1 to match the paper.
	for j := 1; j <= m; j++ {
		index.scores[j-1] = float64(j*(j+1)) / float64(4*(m+1)*(m+2))
	}
	for j := m + 1; j <= 2*m; j++ {
		index.scores[j-1] = 1 - float64(2*m+1-j)/float64(m+1) + float64((2*m+1-j)*(2*m+2-j))/float64(4*(m+1)*(m+2))
	}
	index.genPerturbSets()
}

func (index *MultiprobeIndex) getScore(ps *perturbSet) float64 {
	score := 0.0
	for j := range *ps {
		score += index.scores[j-1]
	}
	return score
}

func (ps perturbSet) isValid(m int) bool {
	for key := range ps {
		// At most one perturbation on same index.
		if _, ok := ps[2*m+1-key]; ok {
			return false
		}
		// No keys larger than 2m.
		if key > 2*m {
			return false
		}
	}
	return true
}

func (ps perturbSet) shift() perturbSet {
	next := make(perturbSet)
	max := 0
	for k := range ps {
		if k > max {
			max = k
		}
		next[k] = true
	}
	delete(next, max)
	next[max+1] = true
	return ps
}

func (ps perturbSet) expand() perturbSet {
	next := make(perturbSet)
	max := 0
	for k := range ps {
		if k > max {
			max = k
		}
		next[k] = true
	}
	next[max+1] = true
	return ps
}

func (index *MultiprobeIndex) genPerturbSets(t int) {
	setHeap := make(perturbSetHeap, 1)
	currentTop := map[int]bool{1: true}
	setHeap[0] = perturbSetPair{
		perturbSet: perturbSet,
		score:      index.getScore(&perturbSet),
	}
	heap.Init(&setHeap)
	index.perturbSets = make([]perturbSet, t)
	m := index.SimpleIndex.LshSettings.m

	for i := 0; i < t; i++ {
		for counter := 0; true; counter++ {
			currentTop := heap.Pop(&setHeap).(perturbSetPair)
			heap.Push(&setHeap, currentTop.shift())
			heap.Push(&setHeap, currentTop.expand())

			if currentTop.isValid(m) {
				index.perturbSet[i] = currentTop
				break
			}
			if counter >= 2*m {
				panic("too many iterations, probably infinite loop!")
			}
		}
	}
}

func (index *MultiprobeIndex) Insert(point Point, id int) {
	index.SimpleIndex.Insert(point, id)
}

func (index *MultiprobeIndex) queryHelper(tableKeys []TableKey) []int {
	// Apply hash functions
	hvs := index.toSimpleKeys(tableKeys)

	// Lookup in each table.
	candidatesAll := make([]int, 0)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				candidatesAll = append(candidatesAll, id)
			}
		}
	}
	return candidatesAll
}

// perturb returns the result of applying perturbation on each baseKey.
func (index *MultiprobeIndex) perturb(baseKey []TableKey, perturbation TableKey) []TableKey {
	// TODO(cmei): Apply perturbation
	return baseKey
}

func (index *MultiprobeIndex) QueryK(q Point, k int, out chan int) {
	baseKey := index.Hash(q)
	candidates := make([]int, 0)
	seens := make(map[int]bool)
	for i := 0; i < len(index.probeSeq) && len(candidates) < k; i++ {
		// Generate new hash key based on perturbation.
		perturbedKey := index.perturb(baseKey, index.probeSeq[i])

		// Perform lookup.
		neighbours := index.queryHelper(perturbedKey)

		// Append new candidates to index.
		for _, id := range neighbours {
			if _, seen := seens[id]; !seen {
				candidates = append(candidates, id)
			}
		}
	}
}
