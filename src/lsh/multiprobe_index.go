package lsh

type MultiprobeIndex struct {
	*SimpleIndex
	// A list of perturbations that will be used for lookups.
	probeSeq []TableKey
}

func NewMultiprobeLsh(dim, l, m int, w float64) *MultiprobeIndex {
	index := &MultiprobeIndex{
		SimpleIndex: NewSimpleLsh(dim, l, m, w),
	}
	index.initProbeSequence()
	return index
}

func (index *MultiprobeIndex) initProbeSequence() {
	// TODO(cmei): Initialize sequence.
	index.probeSeq = make([]TableKey, 0)
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
