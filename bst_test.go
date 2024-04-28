package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanBalanceBST(t *testing.T) {
	// Arrange
	keys := []uint64{10, 5, 15, 30, 20, 25, 35}

	// Act
	tree := &BinarySearchTree{}
	for _, key := range keys {
		tree.Insert(key)
	}

	tree.Exists(25)
	initialHops := tree.LastHopsCheck

	fmt.Printf("\n initial Tree \n")
	tree.rootNode.Print("", true)

	tree.Balance()
	tree.Exists(25)
	balancedHops := tree.LastHopsCheck

	fmt.Printf("\n balanced Tree \n")
	tree.rootNode.Print("", true)

	// Assert
	assert.Equal(t, 4, initialHops, "initial hops")
	assert.Equal(t, 2, balancedHops, "balanced hops")
}

func TestBSTCanDelete(t *testing.T) {
	keys := []uint64{10, 5, 15, 30, 20, 25, 35}

	// Act
	tree := &BinarySearchTree{}
	for _, key := range keys {
		tree.Insert(key)
	}

	tree.Balance()

	assert.Equal(t, true, tree.Exists(25))
	tree.Delete(25)

	assert.Equal(t, false, tree.Exists(25))
}

func TestBSTCanInsert(t *testing.T) {
	// Arrange
	keys := []uint64{10, 5, 15, 30, 20, 25, 35}

	// Act
	tree := &BinarySearchTree{}
	for _, key := range keys {
		tree.Insert(key)
	}

	// Assert
	assert.Equal(t, true, tree.Exists(20), fmt.Sprintf("%d exists", 20))

	for _, key := range keys {
		assert.Equal(t, true, tree.Exists(key), fmt.Sprintf("%d exists", key))
	}

	// Hop checker
	assert.Equal(t, true, tree.Exists(35))
	assert.Equal(t, 3, tree.LastHopsCheck)

	assert.Equal(t, true, tree.Exists(30))
	assert.Equal(t, 2, tree.LastHopsCheck)

	assert.Equal(t, false, tree.Exists(36))
	assert.Equal(t, 3, tree.LastHopsCheck)

	// Tree miss check
	assert.Equal(t, false, tree.Exists(0))
	assert.Equal(t, false, tree.Exists(1))
	assert.Equal(t, false, tree.Exists(9))
	assert.Equal(t, false, tree.Exists(21))
	assert.Equal(t, false, tree.Exists(36))

}
