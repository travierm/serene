package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type BPlusTree struct {
	root      *BPlusTreeNode
	maxDegree uint // max num of children
	maxKeys   uint
}

type BPlusTreeNode struct {
	keys     []int
	children []*BPlusTreeNode
	isLeaf   bool
	next     *BPlusTreeNode // next node sequentially
}

func NewBPlusTree(maxDegree uint) *BPlusTree {
	return &BPlusTree{
		maxDegree: maxDegree,
		maxKeys:   maxDegree - 1,
	}
}

/**
 * node - the node that is splitting
 * key - promote key
 * newNode - newNode created when splitting
 */
func (tree *BPlusTree) insertIntoParent(node *BPlusTreeNode, key int, newNode *BPlusTreeNode) {
	parent := tree.findParent(tree.root, node)
	if parent == nil { // This could happen if node is the root
		newRoot := &BPlusTreeNode{
			keys:     []int{key},
			children: []*BPlusTreeNode{node, newNode},
			isLeaf:   false,
		}
		tree.root = newRoot
		return
	}

	// Insert key and new node in the parent node
	index := sort.Search(len(parent.keys), func(i int) bool { return key < parent.keys[i] })
	parent.keys = append(parent.keys, 0)             // Expand the slice by adding a zero value
	copy(parent.keys[index+1:], parent.keys[index:]) // Shift elements right
	parent.keys[index] = key

	// Insert new node in the children slice
	parent.children = append(parent.children, nil)             // Expand the slice by adding a nil value
	copy(parent.children[index+2:], parent.children[index+1:]) // Shift elements right
	parent.children[index+1] = newNode

	// If the parent is now too large, split the parent
	if len(parent.keys) > int(tree.maxKeys) {
		if parent == tree.root {
			tree.SplitRoot()
		} else {
			tree.splitNode(parent)
		}
	}
}

func (tree *BPlusTree) findParent(current, child *BPlusTreeNode) *BPlusTreeNode {
	if current == nil || current.isLeaf {
		return nil
	}

	for _, c := range current.children {
		if c == child {
			return current
		}
	}

	for _, c := range current.children {
		foundParent := tree.findParent(c, child)
		if foundParent != nil {
			return foundParent
		}
	}

	return nil
}
func (tree *BPlusTree) Insert(key int, value interface{}) {

	if tree.root == nil {
		tree.root = &BPlusTreeNode{
			keys:   []int{key}, // Directly insert the first key
			isLeaf: true,
		}
		return
	}

	leaf := tree.findLeaf(key)
	// Insert the key in sorted order without needing to sort each time
	index := sort.SearchInts(leaf.keys, key)
	leaf.keys = append(leaf.keys[:index], append([]int{key}, leaf.keys[index:]...)...)

	if len(tree.root.keys) > int(tree.maxKeys) {
		tree.SplitRoot()
		return
	}

	if len(leaf.keys) > int(tree.maxKeys) {
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
	if len(root.keys) < int(tree.maxKeys-1) {
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
	if len(node.keys) <= int(tree.maxKeys) {
		return
	}

	midIndex := len(node.keys) / 2
	promoteKey := node.keys[midIndex]

	if len(node.children) > int(tree.maxKeys) {
		midIndex = midIndex + 1
	}

	// Create new node with the right half of the keys
	newNode := &BPlusTreeNode{
		keys:   append([]int{}, node.keys[midIndex:]...), // this is fucked
		isLeaf: node.isLeaf,
		next:   node.next,
	}
	node.keys = node.keys[:midIndex]

	if node.isLeaf {
		node.next = newNode // Maintain linked list of leaves
	} else {
		// If not a leaf, also split children
		newNode.children = node.children[midIndex+1:]
		node.children = node.children[:midIndex+1]
	}

	if node == tree.root {
		// Create new root when the root is split
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

// RenderAsJSON returns a JSON representation of the entire B+ tree.
func (tree *BPlusTree) RenderAsJSON() (string, error) {
	if tree.root == nil {
		return "{}", nil
	}
	data, err := json.Marshal(tree.root)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MarshalJSON customizes the JSON representation of a BPlusTreeNode.
func (node *BPlusTreeNode) MarshalJSON() ([]byte, error) {
	type Alias BPlusTreeNode
	return json.Marshal(&struct {
		Keys     []int            `json:"keys"`
		Children []*BPlusTreeNode `json:"children,omitempty"` // omit if no children
		IsLeaf   bool             `json:"isLeaf"`
		Next     *BPlusTreeNode   `json:"next,omitempty"` // omit if no next node
		*Alias
	}{
		Keys:     node.keys,
		Children: node.children,
		IsLeaf:   node.isLeaf,
		Next:     node.next,
		Alias:    (*Alias)(node),
	})
}
