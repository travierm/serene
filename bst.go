package main

import (
	"errors"
	"fmt"
	"strconv"
)

type BinarySearchTree struct {
	rootNode      *BSTNode
	LastHopsCheck int
}

type BSTNode struct {
	Left  *BSTNode
	Right *BSTNode
	Key   uint64
}

func (n *BSTNode) Print(prefix string, isTail bool) {
	if n.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		n.Right.Print(newPrefix, false)
	}

	fmt.Println(prefix + (func() string {
		if isTail {
			return "└── "
		}
		return "┌── "
	}()) + strconv.FormatUint(n.Key, 10))

	if n.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		n.Left.Print(newPrefix, true)
	}
}

func (n *BSTNode) Insert(key uint64) error {
	if key > n.Key {
		// right has not been set
		if n.Right == nil {
			n.Right = &BSTNode{
				Key: key,
			}

			return nil
		}

		n.Right.Insert(key)
	}

	if key < n.Key {
		// left has not been set
		if n.Left == nil {
			n.Left = &BSTNode{
				Key: key,
			}

			return nil
		}

		n.Left.Insert(key)
	}

	return errors.New("duplicate key passed")
}

func (t *BinarySearchTree) Exists(key uint64) bool {
	exists, hops := t.rootNode.Exists(key, 0)

	t.LastHopsCheck = hops

	return exists
}

func (n *BSTNode) Exists(key uint64, hop int) (bool, int) {
	if n.Key == key {
		return true, hop
	}

	if key < n.Key && n.Left != nil {
		return n.Left.Exists(key, hop+1)
	}

	if key > n.Key && n.Right != nil {
		return n.Right.Exists(key, hop+1)
	}

	return false, hop
}

func (t *BinarySearchTree) Insert(key uint64) {
	if t.rootNode == nil {
		t.rootNode = &BSTNode{
			Key: key,
		}
	}

	t.rootNode.Insert(key)
}

// todo: add mutex lock when balancing to help with root replacement while an insert is happening
func (t *BinarySearchTree) Balance() {
	if t.rootNode == nil {
		return
	}

	records := t.rootNode.GetOrderedList()

	mid := len(records) / 2  //find mid element
	left := records[:mid]    //store left slice elements in left variable
	right := records[mid+1:] // right always stores the middle number and we need to skip over it

	insertList := make([]uint64, 0, len(records))
	insertList = append(insertList, left...) // Add left elements
	insertList = append(insertList, right...)

	newRoot := BSTNode{
		Key: records[mid],
	}

	for _, node := range insertList {
		newRoot.Insert(node)
	}

	t.rootNode = &newRoot
}

func (n *BSTNode) GetOrderedList() []uint64 {
	records := n.GetValues("left", []uint64{})
	records = append(records, n.Key)
	records = append(records, n.GetValues("right", []uint64{})...)

	return records
}

func (n *BSTNode) GetValues(direction string, previousRecords []uint64) []uint64 {
	if direction == "left" {
		if n.Left != nil {
			leftValues := n.Left.GetValues("left", previousRecords)
			rightValues := n.Left.GetValues("right", previousRecords)

			previousRecords := append(previousRecords, leftValues...)
			previousRecords = append(previousRecords, n.Left.Key)
			previousRecords = append(previousRecords, rightValues...)

			return previousRecords
		}

		return previousRecords
	}

	if direction == "right" {
		if n.Right != nil {
			leftValues := n.Right.GetValues("left", previousRecords)
			rightValues := n.Right.GetValues("right", previousRecords)

			previousRecords := append(previousRecords, leftValues...)
			previousRecords = append(previousRecords, n.Right.Key)
			previousRecords = append(previousRecords, rightValues...)

			return previousRecords
		}

		return previousRecords
	}

	return previousRecords
}
