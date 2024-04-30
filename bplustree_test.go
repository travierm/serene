package main

import (
	"fmt"
	"testing"
)

func TestCanInsertNodes(t *testing.T) {
	//keys := []int{10, 5, 15, 30, 20, 25, 35}
	keys := []int{10, 5, 15, 30, 12}

	tree := &BPlusTree{
		minChildren: 2,
		maxChildren: 3,
	}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	tree.Print()
}

func TestCanFindLeaf(t *testing.T) {
	n1 := &BPlusTreeNode{keys: []int{1, 2, 3}, isLeaf: true}
	n2 := &BPlusTreeNode{keys: []int{4, 5, 6}, isLeaf: true}
	n3 := &BPlusTreeNode{keys: []int{7, 8, 9}, isLeaf: true}
	n4 := &BPlusTreeNode{keys: []int{10, 11, 12}, isLeaf: true}

	root := &BPlusTreeNode{
		keys:     []int{4, 7},
		children: []*BPlusTreeNode{n1, n2, n3},
		isLeaf:   false,
	}

	root.next = n4 // example of linking nodes sequentially
	tree := &BPlusTree{root: root, minChildren: 2, maxChildren: 3}

	tree.Print()
}
