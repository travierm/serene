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
