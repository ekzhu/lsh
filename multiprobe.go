package lsh

import (
	"container/heap"
	"math/rand"
)

type perturbSet map[int]bool

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
	return next
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
	return next
}

// A pair of perturbation set and its score.
type perturbSetPair struct {
	ps    perturbSet
	score float64
}

// perturbSetHeap is a min-heap of perturbSetPairs.
type perturbSetHeap []perturbSetPair

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

// MultiprobeLsh implements the Multi-probe LSH algorithm by Qin Lv et.al.
// The Multi-probe LSH does not support k-NN query directly.
type MultiprobeLsh struct {
	*BasicLsh
	// The size of our probe sequence.
	t int

	// The scores of perturbation values.
	scores []float64

	perturbSets []perturbSet

	// Each hash table has a list of perturbation vectors
	// each perturbation vector is list of -+ 1 or 0 that will
	// be applied to the hashTableKey of the query hash value
	// t x l x m
	perturbVecs [][][]int
}

// NewMultiprobeLsh creates a new Multi-probe LSH for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, and w is the slot size for the
// family of LSH functions.
// t is the number of perturbation vectors that will be applied to
// each query.
// Increasing t increases the running time of the Query function.
func NewMultiprobeLsh(dim, l, m int, w float64, t int) *MultiprobeLsh {
	index := &MultiprobeLsh{
		BasicLsh: NewBasicLsh(dim, l, m, w),
		t:        t,
	}
	index.initProbeSequence()
	return index
}

func (index *MultiprobeLsh) initProbeSequence() {
	m := index.m
	index.scores = make([]float64, 2*m)
	// Use j's starting from 1 to match the paper.
	for j := 1; j <= m; j++ {
		index.scores[j-1] = float64(j*(j+1)) / float64(4*(m+1)*(m+2))
	}
	for j := m + 1; j <= 2*m; j++ {
		index.scores[j-1] = 1 - float64(2*m+1-j)/float64(m+1) + float64((2*m+1-j)*(2*m+2-j))/float64(4*(m+1)*(m+2))
	}
	index.genPerturbSets()
	index.genPerturbVecs()
}

func (index *MultiprobeLsh) getScore(ps *perturbSet) float64 {
	score := 0.0
	for j := range *ps {
		score += index.scores[j-1]
	}
	return score
}

func (index *MultiprobeLsh) genPerturbSets() {
	setHeap := make(perturbSetHeap, 1)
	start := perturbSet{1: true}
	setHeap[0] = perturbSetPair{
		ps:    start,
		score: index.getScore(&start),
	}
	heap.Init(&setHeap)
	index.perturbSets = make([]perturbSet, index.t)
	m := index.m

	for i := 0; i < index.t; i++ {
		for counter := 0; true; counter++ {
			currentTop := heap.Pop(&setHeap).(perturbSetPair)
			nextShift := currentTop.ps.shift()
			heap.Push(&setHeap, perturbSetPair{
				ps:    nextShift,
				score: index.getScore(&nextShift),
			})
			nextExpand := currentTop.ps.expand()
			heap.Push(&setHeap, perturbSetPair{
				ps:    nextExpand,
				score: index.getScore(&nextExpand),
			})

			if currentTop.ps.isValid(m) {
				index.perturbSets[i] = currentTop.ps
				break
			}
			if counter >= 2*m {
				panic("too many iterations, probably infinite loop!")
			}
		}
	}
}

func (index *MultiprobeLsh) genPerturbVecs() {
	// First we need to generate the permutation tables
	// that maps the ids of the unit perturbation in each
	// perturbation set to the index of the unit hash
	// value
	perms := make([][]int, index.l)
	for i := range index.tables {
		random := rand.New(rand.NewSource(int64(i)))
		perm := random.Perm(index.m)
		perms[i] = make([]int, index.m*2)
		for j := 0; j < index.m; j++ {
			perms[i][j] = perm[j]
		}
		for j := 0; j < index.m; j++ {
			perms[i][j+index.m] = perm[index.m-1-j]
		}
	}

	// Generate the vectors
	index.perturbVecs = make([][][]int, len(index.perturbSets))
	for i, ps := range index.perturbSets {
		perTableVecs := make([][]int, index.l)
		for j := range perTableVecs {
			vec := make([]int, index.m)
			for k := range ps {
				mapped_ind := perms[j][k-1]
				if k > index.m {
					// If it is -1
					vec[mapped_ind] = -1
				} else {
					// if it is +1
					vec[mapped_ind] = 1
				}
			}
			perTableVecs[j] = vec
		}
		index.perturbVecs[i] = perTableVecs
	}
}

func (index *MultiprobeLsh) queryHelper(tableKeys []hashTableKey, out chan<- string) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(tableKeys)

	// Lookup in each table.
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				out <- id
			}
		}
	}
}

// perturb returns the result of applying perturbation on each baseKey.
func (index *MultiprobeLsh) perturb(baseKey []hashTableKey, perturbation [][]int) []hashTableKey {
	if len(baseKey) != len(perturbation) {
		panic("Number tables does not match with number of perturb vecs")
	}
	perturbedTableKeys := make([]hashTableKey, len(baseKey))
	for i, p := range perturbation {
		perturbedTableKeys[i] = make(hashTableKey, index.m)
		for j, h := range baseKey[i] {
			perturbedTableKeys[i][j] = h + p[j]
		}
	}
	return perturbedTableKeys
}

// Query finds the ids of nearest neighbour candidates,
// given the query point
func (index *MultiprobeLsh) Query(q Point) []string {
	// Hash
	baseKey := index.hash(q)
	// Query
	results := make(chan string)
	go func() {
		defer close(results)
		for i := 0; i < len(index.perturbVecs)+1; i++ {
			perturbedTableKeys := baseKey
			if i != 0 {
				// Generate new hash key based on perturbation.
				perturbedTableKeys = index.perturb(baseKey, index.perturbVecs[i-1])
			}
			// Perform lookup.
			index.queryHelper(perturbedTableKeys, results)
		}
	}()
	seen := make(map[string]bool)
	for id := range results {
		if _, exist := seen[id]; exist {
			continue
		}
		seen[id] = true
	}
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}
