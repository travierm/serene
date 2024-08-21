package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanSaveAndLoadBinaryFile(t *testing.T) {
	tree := NewBPlusTree(3)

	for i := 0; i < 100; i++ {
		key := i
		tree.Insert(key, fmt.Sprintf("Node #%d", key))
	}

	err := tree.SaveToBinaryFile("storage/test/bplustree.bin")
	if err != nil {
		t.Fatalf(err.Error())
	}

	loadedTree, err := LoadFromBinaryFile("storage/test/bplustree.bin")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, tree.root.keys, loadedTree.root.keys)

	ClearTestFolder()
}
