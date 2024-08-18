package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCanInsertNodes(t *testing.T) {
// 	//keys := []int{10, 5, 15, 30, 20, 25, 35}
// 	//keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// 	keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// 	tree := NewBPlusTree(3)

// 	for _, key := range keys {
// 		tree.Insert(key, fmt.Sprintf("Node #%d", key))
// 		tree.Print()
// 	}
// }

func TestIsShapedCorrectlyAfter2(t *testing.T) {
	tree := NewBPlusTree(3)
	keys := []int{1, 2}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	assert.True(t, len(tree.root.keys) == 2, "root has one key")
	assert.True(t, tree.root.keys[0] == 1, "root key is 1")
	assert.True(t, tree.root.keys[1] == 2, "root key is 1")
	assert.True(t, len(tree.root.children) == 0, "root has no children")

	tree.Print()
}

func TestIsShapedCorrectlyAfter3(t *testing.T) {
	tree := NewBPlusTree(3)
	keys := []int{1, 2, 3}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	assert.True(t, len(tree.root.keys) == 1, "root has one key")
	assert.True(t, tree.root.keys[0] == 2, "root key is 2")
	assert.Equal(t, []int{1}, tree.root.children[0].keys, "1 is the left most child")
	assert.Equal(t, []int{2, 3}, tree.root.children[1].keys, "2,3 is the right most children")

	tree.Print()
}

func TestIsShapedCorrectlyAfter4(t *testing.T) {
	tree := NewBPlusTree(3)
	keys := []int{1, 2, 3, 4}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	assert.True(t, len(tree.root.keys) == 2, "root has two keys")
	assert.True(t, len(tree.root.children) == 3, "root has three children")
	assert.Equal(t, []int{2, 3}, tree.root.keys)
	assert.Equal(t, []int{1}, tree.root.children[0].keys)
	assert.Equal(t, []int{2}, tree.root.children[1].keys)
	assert.Equal(t, []int{3, 4}, tree.root.children[2].keys)

	tree.Print()
}

func TestIsShapedCorrectlyAfter5(t *testing.T) {
	tree := NewBPlusTree(3)
	keys := []int{1, 2, 3, 4, 5}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	// root
	assert.Equal(t, []int{3}, tree.root.keys)

	// 2nd layer
	assert.Equal(t, []int{2}, tree.root.children[0].keys) // right
	assert.Equal(t, []int{4}, tree.root.children[1].keys) // left

	// 3rd layer
	// left
	assert.Equal(t, []int{1}, tree.root.children[0].children[0].keys)
	assert.Equal(t, []int{2}, tree.root.children[0].children[1].keys)

	// right
	assert.Equal(t, []int{3}, tree.root.children[1].children[0].keys)
	assert.Equal(t, []int{4, 5}, tree.root.children[1].children[1].keys)

	tree.Print()
}

func TestIsShapedCorrectlyAfter6(t *testing.T) {
	tree := NewBPlusTree(3)
	keys := []int{1, 2, 3, 4, 5, 6}

	for _, key := range keys {
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	// root
	assert.Equal(t, []int{3}, tree.root.keys)

	// 2nd layer
	assert.Equal(t, []int{2}, tree.root.children[0].keys)    // right
	assert.Equal(t, []int{4, 5}, tree.root.children[1].keys) // left

	// 3rd layer
	// left
	assert.Equal(t, []int{1}, tree.root.children[0].children[0].keys)
	assert.Equal(t, []int{2}, tree.root.children[0].children[1].keys)

	// right
	assert.Equal(t, []int{3}, tree.root.children[1].children[0].keys)
	assert.Equal(t, []int{4}, tree.root.children[1].children[1].keys)
	assert.Equal(t, []int{5, 6}, tree.root.children[1].children[2].keys)

	tree.Print()
}

// func TestIsShapedCorrectlyAfter7(t *testing.T) {
// 	tree := NewBPlusTree(3)
// 	keys := []int{1, 2, 3, 4, 5, 6, 7}

// 	for _, key := range keys {
// 		tree.Insert(key, fmt.Sprintf("Node #%d", key))
// 	}

// 	tree.Print()
// 	return

// 	// root
// 	assert.Equal(t, []int{3, 5}, tree.root.keys)

// 	// 2nd layer
// 	assert.Equal(t, []int{2}, tree.root.children[0].keys) // right
// 	assert.Equal(t, []int{4}, tree.root.children[1].keys) // middle
// 	assert.Equal(t, []int{6}, tree.root.children[2].keys) // left

// 	// 3rd layer
// 	// left
// 	assert.Equal(t, []int{1}, tree.root.children[0].children[0].keys)
// 	assert.Equal(t, []int{2}, tree.root.children[0].children[1].keys)

// 	// middle
// 	assert.Equal(t, []int{3}, tree.root.children[1].children[0].keys)
// 	assert.Equal(t, []int{4}, tree.root.children[1].children[1].keys)

// 	// right
// 	assert.Equal(t, []int{5}, tree.root.children[2].children[0].keys)
// 	assert.Equal(t, []int{6, 7}, tree.root.children[2].children[1].keys)

// 	tree.Print()
// }

// func TestCanHandleOneHundredNodes(t *testing.T) {
// 	tree := NewBPlusTree(3)

// 	for i := 0; i < 8; i++ {
// 		tree.Insert(i, fmt.Sprintf("Node #%d", i))
// 	}

// 	tree.Print()
// }

// func TestCanFindLeaf(t *testing.T) {
// 	n1 := &BPlusTreeNode{keys: []int{1, 2, 3}, isLeaf: true}
// 	n2 := &BPlusTreeNode{keys: []int{4, 5, 6}, isLeaf: true}
// 	n3 := &BPlusTreeNode{keys: []int{7, 8, 9}, isLeaf: true}
// 	n4 := &BPlusTreeNode{keys: []int{10, 11, 12}, isLeaf: true}

// 	root := &BPlusTreeNode{
// 		keys:     []int{4, 7},
// 		children: []*BPlusTreeNode{n1, n2, n3},
// 		isLeaf:   false,
// 	}

// 	root.next = n4 // example of linking nodes sequentially
// 	tree := NewBPlusTree(3)

// 	tree.Print()
// }
