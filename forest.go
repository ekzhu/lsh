package lsh

import (
	"fmt"
	"sync"
)

type treeNode struct {
	// Hash key for this intermediate node. nil/empty for root nodes.
	hashKey int
	// A list of ids to the source dataset, only leaf nodes have non-empty ids.
	ids []string
	// Child nodes, keyed by partial hash value.
	children map[int]*treeNode
}

func (node *treeNode) recursiveDelete() {
	for _, child := range node.children {
		if len((child).children) > 0 {
			(child).recursiveDelete()
		}
		if len(child.ids) > 0 {
			node.ids = nil
		}
	}
	node.ids = nil
	node.children = nil
}

// recursiveAdd recurses down the tree to find the correct location to insert id.
// Returns whether a new hash value was added.
func (node *treeNode) recursiveAdd(level int, id string, tableKey hashTableKey) bool {
	if level == len(tableKey) {
		node.ids = append(node.ids, id)
		return false
	}
	// Check if next hash exists in children map. If not, create.
	var next *treeNode
	hasNewHash := false
	if nextNode, ok := node.children[tableKey[level]]; !ok {
		next = &treeNode{
			hashKey:  tableKey[level],
			ids:      make([]string, 0),
			children: make(map[int]*treeNode),
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

func tab(times int) {
	for i := 0; i < times; i++ {
		fmt.Print("    ")
	}
}

func (node *treeNode) dump(level int) {
	tab(level)
	fmt.Printf("{ (%v): ids %v ", node.hashKey, node.ids)
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

type prefixTree struct {
	// Number of distinct elements in the tree.
	count int
	// Pointer to the root node.
	root *treeNode
}

func (tree *prefixTree) insertIntoTree(id string, tableKey hashTableKey) {
	if tree.root.recursiveAdd(0, id, tableKey) {
		tree.count++
	}
}

// lookup find ids and write them to out channel
func (tree *prefixTree) lookup(maxLevel int, tableKey hashTableKey,
	done <-chan struct{}, out chan<- string) {
	currentNode := tree.root
	for level := 0; level < len(tableKey) && level < maxLevel; level++ {
		if next, ok := currentNode.children[tableKey[level]]; ok {
			currentNode = next
		} else {
			return
		}
	}

	// Grab all ids of nodes descendent from the current node.
	queue := []*treeNode{currentNode}
	for len(queue) > 0 {
		// Add node's ids to main list.
		for _, id := range queue[0].ids {
			select {
			case out <- id:
			case <-done:
				return
			}
		}

		// Add children.
		for _, child := range queue[0].children {
			queue = append(queue, child)
		}

		// Done with head.
		queue = queue[1:]
	}
}

// LshForest implements the LSH Forest algorithm by Mayank Bawa et.al.
// It supports both nearest neighbour candidate query and k-NN query.
type LshForest struct {
	// Embedded type
	*lshParams
	// Trees.
	trees []prefixTree
}

// NewLshForest creates a new LSH Forest for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, w is the slot size for the
// family of LSH functions.
func NewLshForest(dim, l, m int, w float64) *LshForest {
	trees := make([]prefixTree, l)
	for i := range trees {
		trees[i].count = 0
		trees[i].root = &treeNode{
			hashKey:  0,
			ids:      make([]string, 0),
			children: make(map[int]*treeNode),
		}
	}
	return &LshForest{
		lshParams: newLshParams(dim, l, m, w),
		trees:     trees,
	}
}

// Delete releases the memory used by this index.
func (index *LshForest) Delete() {
	for _, tree := range index.trees {
		(*tree.root).recursiveDelete()
	}
}

// Insert adds a new data point to the LSH Forest.
// id is the unique identifier for the data point.
func (index *LshForest) Insert(point Point, id string) {
	// Apply hash functions.
	hvs := index.hash(point)
	// Parallel insert
	var wg sync.WaitGroup
	wg.Add(len(index.trees))
	for i := range index.trees {
		hv := hvs[i]
		tree := &(index.trees[i])
		go func(tree *prefixTree, hv hashTableKey) {
			tree.insertIntoTree(id, hv)
			wg.Done()
		}(tree, hv)
	}
	wg.Wait()
}

// Helper that queries all trees and returns an channel ids.
func (index *LshForest) queryHelper(maxLevel int, tableKeys []hashTableKey, done <-chan struct{}, out chan<- string) {
	var wg sync.WaitGroup
	wg.Add(len(index.trees))
	for i := range index.trees {
		key := tableKeys[i]
		tree := index.trees[i]
		go func() {
			tree.lookup(maxLevel, key, done, out)
			wg.Done()
		}()
	}
	wg.Wait()
}

// Query finds at top-k ids of approximate nearest neighbour candidates,
// in unsorted order, given the query point.
func (index *LshForest) Query(q Point, k int) []string {
	// Apply hash functions
	hvs := index.hash(q)
	// Query
	results := make(chan string)
	done := make(chan struct{})
	go func() {
		for maxLevels := index.m; maxLevels >= 0; maxLevels-- {
			select {
			case <-done:
				return
			default:
				index.queryHelper(maxLevels, hvs, done, results)
			}
		}
		close(results)
	}()
	seen := make(map[string]bool)
	for id := range results {
		if len(seen) >= k {
			break
		}
		if _, exist := seen[id]; exist {
			continue
		}
		seen[id] = true
	}
	close(done)
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// Dump prints out the index for debugging
func (index *LshForest) dump() {
	for i, tree := range index.trees {
		fmt.Printf("Tree %d (%d hash values):\n", i, tree.count)
		tree.root.dump(0)
	}
}
