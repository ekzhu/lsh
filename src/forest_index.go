package lsh

import (
	"fmt"
)

type TreeNode struct {
	// Hash key for this intermediate node. nil/empty for root nodes.
	hashKey Key
	// A list of indices to the source dataset.
	indices []int
	// Child nodes.
	children []TreeNode
}

type Tree struct {
	// Number of distinct elements in the tree.
	count int
	// Pointer to the root node.
	root TreeNode
}

type ForestIndex struct {
	// Embedded type
	*Lsh
	// Number of leaves in the tree.
	count int
	// Trees.
	trees []Tree
}

func NewLshForest(dim, l, m int, w float64) *ForestIndex {
	trees := make([]Tree, l)
	return &ForestIndex{
		Lsh:   NewLsh(m, l, dim, w),
		count: 0,
		trees: trees,
	}
}

// Inserts a point into the index.
func (index *ForestIndex) Insert(point Point) {

}
