package main

import (
	"fmt"
	"sort"
	"strings"
)

type BPlusTree struct {
	root        *BPlusTreeNode
	minChildren uint
	maxChildren uint
}

type BPlusTreeNode struct {
	keys     []int
	children []*BPlusTreeNode
	isLeaf   bool
	next     *BPlusTreeNode // next node sequentially
}

func (tree *BPlusTree) insertIntoParent(node *BPlusTreeNode, key int, newNode *BPlusTreeNode) {
	parent := tree.findParent(tree.root, nil, node)

	// find the position in the parent to store key
	index := sort.Search(len(parent.keys), func(i int) bool { return key < parent.keys[i] })

	// Insert the new key and node into the parent
	parent.keys = append(parent.keys[:index], append([]int{key}, parent.keys[index:]...)...)
	parent.children = append(parent.children[:index+1], append([]*BPlusTreeNode{newNode}, parent.children[index+1:]...)...)

	if len(parent.keys) > int(tree.maxChildren) {
		tree.splitNode(parent)
	}
}

func (tree *BPlusTree) findParent(current, parent, child *BPlusTreeNode) *BPlusTreeNode {
	// no children to look through
	if current.isLeaf {
		return nil
	}

	// Iterate over the children of the current node to check if any of them is the child
	// we are looking for. If we find the child, return the parent node as it is the
	// direct parent of the child.
	for _, c := range current.children {
		if c == child {
			return parent
		}
	}

	// If the child was not found directly under the current node, recursively search deeper
	// into each child. This block explores all subtrees under the current node to find where
	// the child might be located.
	for i := range current.children {
		if result := tree.findParent(current.children[i], current, child); result != nil {
			return result
		}
	}

	return nil
}
func (tree *BPlusTree) Insert(key int, value interface{}) {
	if tree.root == nil {
		tree.root = &BPlusTreeNode{
			keys:   []int{},
			isLeaf: true,
		}
	}

	if tree.root.isLeaf {
		tree.root.keys = append(tree.root.keys, key)

		sort.Ints(tree.root.keys)
		if len(tree.root.keys) > int(tree.maxChildren) {
			tree.SplitRoot()
		}

		return
	}

	leaf := tree.findLeaf(key)
	// this is breaking when doing the 2nd split of the tree
	leaf.keys = append(leaf.keys, key)
	sort.Ints(leaf.keys)

	if len(leaf.keys) > int(tree.maxChildren) {
		tree.splitNode(leaf)
	}
}

func (tree *BPlusTree) findLeaf(key int) *BPlusTreeNode {
	node := tree.root

	for !node.isLeaf {
		i := sort.Search(len(node.keys), func(i int) bool {
			return key <= node.keys[i]
		})
		node = node.children[i]
	}
	return node
}

func (tree *BPlusTree) Print() {
	if tree.root != nil {
		tree.root.Print("", true)
	}
}

func (tree *BPlusTree) SplitRoot() {
	if tree.root == nil {
		return
	}

	// haven't filled up the node yet
	root := tree.root
	if len(root.keys) < int(tree.maxChildren-1) {
		return
	}

	newRoot := &BPlusTreeNode{
		keys:     []int{},
		children: []*BPlusTreeNode{},
		isLeaf:   false,
	}

	// split the node in two
	mid := len(root.keys) / 2
	leftKeys := root.keys[:mid]
	rightKeys := root.keys[mid:]

	leftChildren := []*BPlusTreeNode{}
	rightChildren := []*BPlusTreeNode{}

	if len(root.children) > 0 {
		leftChildren = root.children[:mid+1]
		rightChildren = root.children[mid+1:]
	}

	// create children
	leftChild := &BPlusTreeNode{
		keys:     leftKeys,
		children: leftChildren,
		isLeaf:   root.isLeaf,
	}
	rightChild := &BPlusTreeNode{
		keys:     rightKeys,
		children: rightChildren,
		isLeaf:   root.isLeaf,
	}

	if root.isLeaf {
		leftChild.next = rightChild
		rightChild.next = root.next
	}

	newRoot.keys = append(newRoot.keys, root.keys[mid])
	newRoot.children = append(newRoot.children, leftChild, rightChild)

	tree.root = newRoot
}

// Node Functions

func (tree *BPlusTree) splitNode(node *BPlusTreeNode) {
	if len(node.keys) <= int(tree.maxChildren) {
		return
	}

	midIndex := len(node.keys) / 2

	newNode := &BPlusTreeNode{
		keys:     make([]int, 0, len(node.keys[midIndex+1:])), // create slice with exact length
		children: nil,
		isLeaf:   node.isLeaf,
		next:     nil,
	}

	// append right half of node the new node
	newNode.keys = append(newNode.keys, node.keys[midIndex+1:]...)

	// retain the left half in the original node
	node.keys = node.keys[:midIndex]

	if node.isLeaf {
		newNode.next = node.next
		newNode.children = nil
		node.next = newNode
	} else {
		newNode.children = append(newNode.children, node.children[midIndex+1:]...)
		node.children = node.children[:midIndex+1]
	}

	// promote middle key to parent
	promoteKey := node.keys[midIndex]
	node.keys = node.keys[:len(node.keys)-1]

	if node == tree.root {
		// create new root
		newRoot := &BPlusTreeNode{
			keys:     []int{promoteKey},
			children: []*BPlusTreeNode{node, newNode},
			isLeaf:   false,
		}

		tree.root = newRoot
	} else {
		tree.insertIntoParent(node, promoteKey, newNode)
	}
}

func (node *BPlusTreeNode) Print(prefix string, isTail bool) {
	if node == nil {
		return
	}

	// Define the connector based on whether the node is at the tail or not
	connector := "├── "
	if isTail {
		connector = "└── "
	}

	// Print the current node's keys
	keysStr := fmt.Sprintf("[%s]", strings.Join(strings.Fields(fmt.Sprint(node.keys)), ", "))
	fmt.Println(prefix + connector + keysStr)

	// Update the prefix for child nodes
	if isTail {
		prefix += "    " // Space for a tail
	} else {
		prefix += "│   " // Continuation for non-tail
	}

	// Recursive call for all children except the last
	for i := 0; i < len(node.children)-1; i++ {
		node.children[i].Print(prefix, false)
	}

	// Handle the last child separately to mark it as tail
	if len(node.children) > 0 {
		node.children[len(node.children)-1].Print(prefix, true)
	}
}
