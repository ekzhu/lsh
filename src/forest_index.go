package lsh

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
func (node *TreeNode) recursiveAdd(level int, id int, tableKey TableKey) bool {
	if level == len(tableKey) {
		node.indices = append(node.indices, id)
		return false
	} else {
		// Check if next hash exists in children map. If not, create.
		var next *TreeNode
		hasNewHash := false
		if node, ok := node.children[tableKey[level]]; !ok {
			next = &TreeNode{
				hashKey:  tableKey[level],
				indices:  make([]int, 0),
				children: make(map[int]*TreeNode),
			}
			node.children[tableKey[level]] = next
			hasNewHash = true
		} else {
			next = node
		}
		// Recurse using next level's hash value.
		return hasNewHash || next.recursiveAdd(level+1, id, tableKey)
	}
}

type Tree struct {
	// Number of distinct elements in the tree.
	count int
	// Pointer to the root node.
	root TreeNode
}

func (tree *Tree) insertIntoTree(id int, tableKey TableKey) {
	if tree.root.recursiveAdd(0, id, tableKey) {
		tree.count++
	}
}

func (tree *Tree) lookup(tableKey TableKey) []int {
	indices := make([]int, 0)
	currentNode := &tree.root
	for level := 0; level < len(tableKey); level++ {
		if next, ok := currentNode.children[tableKey[level]]; !ok {
			return indices
		} else {
			currentNode = next
		}
	}
	for child := range currentNode.indices {
		indices = append(indices, child)
	}
	return indices
}

type ForestIndex struct {
	// Embedded type
	*LshSettings
	// Trees.
	trees []Tree
}

func NewLshForest(dim, l, m int, w float64) *ForestIndex {
	trees := make([]Tree, l)
	for _, treeRoot := range trees {
		treeRoot.root = TreeNode{
			indices:  make([]int, 0),
			children: make(map[int]*TreeNode),
		}
	}
	return &ForestIndex{
		LshSettings: NewLshSettings(m, l, dim, w),
		trees:       trees,
	}
}

// Insert adds a point into the LSH Forest index.
func (index *ForestIndex) Insert(point Point, id int) {
	// Apply hash functions.
	hvs := index.Hash(point)
	for treeId, hv := range hvs {
		index.trees[treeId].insertIntoTree(id, hv)
	}
}

// Query searches for candidate keys given the signature
// and writes them to an output channel
func (index *ForestIndex) Query(q Point, out chan int) {
	// Apply hash functions
	hvs := index.Hash(q)
	// Keep track of keys seen
	seens := make(map[int]bool)
	for i, tree := range index.trees {
		for candidate := range tree.lookup(hvs[i]) {
			if _, seen := seens[candidate]; !seen {
				seens[candidate] = true
				out <- candidate
			}
		}
	}
}
