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

func (t *BinarySearchTree) Delete(key uint64) {
	list := t.rootNode.GetOrderedList()
	filteredList := RemoveKey(list, key, func(x uint64, y uint64) bool {
		return x == y
	})

	t.rootNode = balanceNodeRecursive(filteredList, 0, len(filteredList)-1)
}

func (n *BSTNode) CreateBalancedCopy() *BSTNode {
	keys := n.GetOrderedList()

	return balanceNodeRecursive(keys, 0, len(keys)-1)
}

func balanceNodeRecursive(keys []uint64, left, right int) *BSTNode {
	if left > right {
		return nil
	}

	mid := (left + right) / 2
	node := &BSTNode{Key: keys[mid]}

	node.Left = balanceNodeRecursive(keys, left, mid-1)
	node.Right = balanceNodeRecursive(keys, mid+1, right)

	return node
}

// todo: add mutex lock when balancing to help with root replacement while an insert is happening
func (t *BinarySearchTree) Balance() {
	if t.rootNode == nil {
		return
	}

	t.rootNode = t.rootNode.CreateBalancedCopy()
}

func (n *BSTNode) GetOrderedList() []uint64 {
	records := n.GetValues("left", []uint64{})
	records = append(records, n.Key)
	records = append(records, n.GetValues("right", []uint64{})...)

	return records
}

func (n *BSTNode) GetValues(direction string, previousRecords []uint64) []uint64 {

	if direction == "left" && n.Left != nil {
		leftValues := n.Left.GetValues("left", previousRecords)
		rightValues := n.Left.GetValues("right", previousRecords)

		previousRecords := append(previousRecords, leftValues...)
		previousRecords = append(previousRecords, n.Left.Key)
		previousRecords = append(previousRecords, rightValues...)

		return previousRecords
	}

	if direction == "right" && n.Right != nil {
		leftValues := n.Right.GetValues("left", previousRecords)
		rightValues := n.Right.GetValues("right", previousRecords)

		previousRecords := append(previousRecords, leftValues...)
		previousRecords = append(previousRecords, n.Right.Key)
		previousRecords = append(previousRecords, rightValues...)

		return previousRecords
	}

	return previousRecords
}
