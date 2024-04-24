package main

import (
	"fmt"
	"sync"

	"github.com/google/btree"
)

type BTreeIndex[T any] struct {
	tree *btree.BTree
	lock sync.RWMutex
}

type indexItem[T any] struct {
	key    uint64
	record *Record[T]
}

// Less implements the btree.Item interface for indexItem
func (item indexItem[T]) Less(than btree.Item) bool {
	return item.key < than.(indexItem[T]).key
}

// NewBTreeIndex initializes a B-tree with the given degree.
func NewBTreeIndex[T any](degree int) *BTreeIndex[T] {
	return &BTreeIndex[T]{
		tree: btree.New(degree),
	}
}

// Insert adds a record to the B-tree index.
func (idx *BTreeIndex[T]) Insert(record *Record[T]) {
	idx.lock.Lock()
	defer idx.lock.Unlock()

	item := indexItem[T]{key: record.ID, record: record}
	idx.tree.ReplaceOrInsert(item)
}

// Delete removes a record from the B-tree index by key.
func (idx *BTreeIndex[T]) Delete(id uint64) {
	idx.lock.Lock()
	defer idx.lock.Unlock()

	item := indexItem[T]{key: id}
	idx.tree.Delete(item)
}

// Find retrieves a record from the B-tree index by key.
func (idx *BTreeIndex[T]) Find(id uint64) *Record[T] {
	idx.lock.RLock()
	defer idx.lock.RUnlock()

	item := indexItem[T]{key: id}
	if foundItem := idx.tree.Get(item); foundItem != nil {
		return foundItem.(indexItem[T]).record
	}
	return nil
}

// Integration with the Heap
func (h *Heap[T]) InsertWithIndex(record *Record[T], disableWAL ...bool) error {
	if err := h.Insert(record, disableWAL...); err != nil {
		return err
	}
	h.index.Insert(record)
	return nil
}

func (h *Heap[T]) DeleteWithIndex(id uint64, disableWAL ...bool) error {
	if err := h.Delete(id, disableWAL...); err != nil {
		return err
	}
	h.index.Delete(id)
	return nil
}

func (h *Heap[T]) FindByIdWithinIndex(id uint64) (*Record[T], error) {
	record := h.index.Find(id)
	if record == nil {
		return nil, fmt.Errorf("record with ID %d not found", id)
	}
	return record, nil
}
