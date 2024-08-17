package main

import (
	"fmt"

	"github.com/spy16/kiwi/index/bptree"
)

type IndexStorageType string

const (
	File   IndexStorageType = "file"
	Memory IndexStorageType = "memory"
)

type Index struct {
	tree *bptree.BPlusTree
}

func NewIndex(name string, driverType IndexStorageType) (*Index, error) {
	var (
		err  error
		tree *bptree.BPlusTree // Adjust the type to match the actual return type of bptree.Open
	)

	switch driverType {
	case Memory:
		tree, err = bptree.Open(":memory:", nil)
	case File:
		tree, err = bptree.Open(name, nil)
	default:
		return nil, fmt.Errorf("invalid driver type")
	}

	if err != nil {
		return nil, err
	}

	return &Index{tree: tree}, nil
}

func (index *Index) Put(key []byte, value uint64) error {
	return index.tree.Put(key, value)
}

func (index *Index) Get(key []byte) (uint64, error) {
	return index.tree.Get(key)
}

func (index *Index) Del(key []byte) (uint64, error) {
	return index.tree.Del(key)
}

// from the given key scan left or right depending on reverse flag
func (index *Index) Scan(key []byte, reverse bool, scanFn func(key []byte, v uint64) bool) error {
	return index.tree.Scan(key, reverse, scanFn)
}
