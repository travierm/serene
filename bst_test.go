package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBSTCanInsert(t *testing.T) {
	// Arrange
	keys := []uint64{10, 5, 15, 30, 20, 25, 35}

	// Act
	tree := &BinarySearchTree{}
	for _, key := range keys {
		tree.Insert(key)
	}

	// Assert
	tree.rootNode.Print("", true)

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
