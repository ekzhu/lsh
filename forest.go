package lsh

import (
	"fmt"
	"sync"
)

type TreeNode struct {
	// Hash key for this intermediate node. nil/empty for root nodes.
	hashKey int
	// A list of indices to the source dataset.
	indices []int
	// Child nodes, keyed by partial hash value.
	children map[int]*TreeNode
}

// recursiveAdd recurses down the tree to find the correct location to insert id.
// Returns whether a new hash value was added.
func (node *TreeNode) recursiveAdd(level int, id int, tableKey hashTableKey) bool {
	if level == len(tableKey) {
		node.indices = append(node.indices, id)
		return false
	} else {
		// Check if next hash exists in children map. If not, create.
		var next *TreeNode
		hasNewHash := false
		if nextNode, ok := node.children[tableKey[level]]; !ok {
			next = &TreeNode{
				hashKey:  tableKey[level],
				indices:  make([]int, 0),
				children: make(map[int]*TreeNode),
			}
			node.children[tableKey[level]] = next
			hasNewHash = true
		} else {
			next = nextNode
		}
		// Recurse using next level's hash value.
		recursive := next.recursiveAdd(level+1, id, tableKey)
		return hasNewHash || recursive
	}
}

func tab(times int) {
	for i := 0; i < times; i++ {
		fmt.Print("    ")
	}
}

func (node *TreeNode) dump(level int) {
	tab(level)
	fmt.Printf("{ (%d): indices %o ", node.hashKey, node.indices)
	if len(node.children) > 0 {
		fmt.Printf("[\n")
		for _, v := range node.children {
			v.dump(level + 1)
		}
		tab(level)
		fmt.Print("] }\n")
	} else {
		fmt.Print("}\n")
	}
}

type Tree struct {
	// Number of distinct elements in the tree.
	count int
	// Pointer to the root node.
	root *TreeNode
}

func (tree *Tree) insertIntoTree(id int, tableKey hashTableKey) {
	if tree.root.recursiveAdd(0, id, tableKey) {
		tree.count++
	}
}

func (tree *Tree) lookup(maxLevel int, tableKey hashTableKey) []int {
	indices := make([]int, 0)
	currentNode := tree.root
	// fmt.Println(tableKey)
	for level := 0; level < len(tableKey) && level < maxLevel; level++ {
		if next, ok := currentNode.children[tableKey[level]]; !ok {
			return indices
		} else {
			currentNode = next
			// fmt.Printf("Found hash key %d at level %d, current hash %d\n", tableKey[level], level, currentNode.hashKey)
		}
	}

	// Grab all indices of nodes descendent from the current node.
	queue := []*TreeNode{currentNode}
	for len(queue) > 0 {
		// Add node's indices to main list.
		indices = append(indices, queue[0].indices...)

		// Add children.
		for _, child := range queue[0].children {
			queue = append(queue, child)
		}

		// Done with head.
		queue = queue[1:]
	}
	// fmt.Printf("Result: %o\n", indices)
	return indices
}

type ForestIndex struct {
	// Embedded type
	*lshParams
	// Trees.
	trees []Tree
}

func NewLshForest(dim, l, m int, w float64) *ForestIndex {
	trees := make([]Tree, l)
	for i, _ := range trees {
		trees[i].count = 0
		trees[i].root = &TreeNode{
			hashKey:  0,
			indices:  make([]int, 0),
			children: make(map[int]*TreeNode),
		}
	}
	return &ForestIndex{
		lshParams: newLshParams(dim, l, m, w),
		trees:     trees,
	}
}

// Insert adds a point into the LSH Forest index.
func (index *ForestIndex) Insert(point Point, id int) {
	// Apply hash functions.
	hvs := index.Hash(point)
	// Parallel insert
	var wg sync.WaitGroup
	for i := range index.trees {
		hv := hvs[i]
		tree := &(index.trees[i])
		wg.Add(1)
		go func(tree *Tree, hv hashTableKey) {
			tree.insertIntoTree(id, hv)
			wg.Done()
		}(tree, hv)
	}
	wg.Wait()
}

// Helper that queries all trees and returns an array of distinct indices.
func (index *ForestIndex) queryHelper(maxLevel int, tableKeys []hashTableKey) []int {
	// Keep track of keys seen
	indices := make([]int, 0)
	seens := make(map[int]bool)
	for i, tree := range index.trees {
		for _, candidate := range tree.lookup(maxLevel, tableKeys[i]) {
			if _, seen := seens[candidate]; !seen {
				seens[candidate] = true
				indices = append(indices, candidate)
			}
		}
	}
	return indices
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (index *ForestIndex) Query(q Point, out chan int) {
	// Apply hash functions
	hvs := index.Hash(q)
	for _, candidate := range index.queryHelper(index.m, hvs) {
		out <- candidate
	}
}

// QueryK queries for the top k approximate closest neighbours.
func (index *ForestIndex) QueryKnn(q Point, k int, out chan int) {
	// Apply hash functions
	hvs := index.Hash(q)
	candidates := make([]int, 0)
	for maxLevels := index.m; maxLevels >= 0; maxLevels-- {
		candidates = index.queryHelper(maxLevels, hvs)
		// Enough candidates at this level, so we can rank and return.
		if len(candidates) >= k {
			break
		}
	}
	for _, candidate := range candidates {
		out <- candidate
	}
}

// Dump prints out the index for debugging
func (index *ForestIndex) Dump() {
	for i, tree := range index.trees {
		fmt.Printf("Tree %d (%d hash values):\n", i, tree.count)
		tree.root.dump(0)
	}
}
