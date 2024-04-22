package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalWriting(t *testing.T) {
	// create a new wal
	wal := NewWAL("storage/test/test.wal")

	// create a new entry
	wal.Log(&WALEntry{
		Operation: "INSERT",
		RecordID:  1,
		Data:      []byte("hello"),
	})
	wal.Log(&WALEntry{
		Operation: "INSERT",
		RecordID:  2,
		Data:      []byte("hello"),
	})

	assert.Equal(t, len(wal.entries), 2)
	wal.Flush()

	assert.Equal(t, len(wal.entries), 0)

	ClearTestFolder()
}
