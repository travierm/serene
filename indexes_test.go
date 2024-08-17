package main

import (
	"encoding/binary"
	"testing"

	"github.com/spy16/kiwi/index/bptree"
	"github.com/stretchr/testify/assert"
)

func TestIndexCrud(t *testing.T) {
	index, err := NewIndex("storage/test/test_index.bin", File)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var count uint32 = 5
	for i := uint32(0); i < count; i++ {
		k, v := genKV(i)

		index.Put(k, v)
	}

	bptree.Print(index.tree)
	return

	sizeBeforeDelete := index.tree.Size()
	_, err = index.Del(uint64ToByte(20))
	if err != nil {
		t.Fatalf(err.Error())
	}

	existsValue, err := index.Get(uint64ToByte(30))

	if err != nil {
		t.Fatalf(err.Error())
	}
	// notExistsValue, err := index.Get([]byte{20})
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }

	assert.Equal(t, uint64(30), existsValue)
	//assert.Equal(t, 0, notExistsValue)

	index.tree.Close()
	index, err = NewIndex("storage/test/test_index.bin", File)
	if err != nil {
		t.Fatalf(err.Error())
	}

	bptree.Print(index.tree)

	assert.Equal(t, sizeBeforeDelete, index.tree.Size())

	ClearTestFolder()
}

func uint64ToByte(i uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], i)
	return b[:]
}

func genKV(i uint32) ([]byte, uint64) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], i)
	return b[:], uint64(i)
}
