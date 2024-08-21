package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (tree *BPlusTree) SaveToBinaryFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write tree metadata
	err = binary.Write(file, binary.LittleEndian, tree.maxDegree)
	if err != nil {
		return fmt.Errorf("error writing maxDegree: %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, tree.maxKeys)
	if err != nil {
		return fmt.Errorf("error writing maxKeys: %v", err)
	}

	// Write the tree structure
	err = tree.writeBinaryNode(file, tree.root)
	if err != nil {
		return fmt.Errorf("error writing tree structure: %v", err)
	}

	return nil
}

func (tree *BPlusTree) writeBinaryNode(w io.Writer, node *BPlusTreeNode) error {
	if node == nil {
		return binary.Write(w, binary.LittleEndian, false) // Write a flag indicating null node
	}

	err := binary.Write(w, binary.LittleEndian, true) // Write a flag indicating non-null node
	if err != nil {
		return err
	}

	// Write node data
	err = binary.Write(w, binary.LittleEndian, node.isLeaf)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(node.keys)))
	if err != nil {
		return err
	}

	for _, key := range node.keys {
		err = binary.Write(w, binary.LittleEndian, int32(key))
		if err != nil {
			return err
		}
	}

	if !node.isLeaf {
		for _, child := range node.children {
			err = tree.writeBinaryNode(w, child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadFromBinaryFile loads the B+ tree structure from a binary file
func LoadFromBinaryFile(filename string) (*BPlusTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var maxDegree, maxKeys uint32
	err = binary.Read(file, binary.LittleEndian, &maxDegree)
	if err != nil {
		return nil, fmt.Errorf("error reading maxDegree: %v", err)
	}
	err = binary.Read(file, binary.LittleEndian, &maxKeys)
	if err != nil {
		return nil, fmt.Errorf("error reading maxKeys: %v", err)
	}

	tree := &BPlusTree{
		maxDegree: uint32(maxDegree),
		maxKeys:   uint32(maxKeys),
	}

	tree.root, err = tree.readBinaryNode(file)
	if err != nil {
		return nil, fmt.Errorf("error reading tree structure: %v", err)
	}

	return tree, nil
}

func (tree *BPlusTree) readBinaryNode(r io.Reader) (*BPlusTreeNode, error) {
	var nodeExists bool
	err := binary.Read(r, binary.LittleEndian, &nodeExists)
	if err != nil {
		return nil, err
	}

	if !nodeExists {
		return nil, nil
	}

	node := &BPlusTreeNode{}

	err = binary.Read(r, binary.LittleEndian, &node.isLeaf)
	if err != nil {
		return nil, err
	}

	var keyCount uint32
	err = binary.Read(r, binary.LittleEndian, &keyCount)
	if err != nil {
		return nil, err
	}

	node.keys = make([]int, keyCount)
	for i := range node.keys {
		var key int32
		err = binary.Read(r, binary.LittleEndian, &key)
		if err != nil {
			return nil, err
		}
		node.keys[i] = int(key)
	}

	if !node.isLeaf {
		node.children = make([]*BPlusTreeNode, keyCount+1)
		for i := range node.children {
			node.children[i], err = tree.readBinaryNode(r)
			if err != nil {
				return nil, err
			}
		}
	}

	return node, nil
}
