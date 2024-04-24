package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
	suite.Suite
}

type TestRecord struct {
	Name   string
	Amount float64
}

func TestHeapCrud(t *testing.T) {
	heap := NewHeap[TestRecord]("products", "storage/test")

	err := heap.Insert(&Record[TestRecord]{ID: 1, Data: TestRecord{Name: "Product 1", Amount: 100.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 2, Data: TestRecord{Name: "Product 2", Amount: 150.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 3, Data: TestRecord{Name: "Product 3", Amount: 200.0}})
	err = heap.Update(&Record[TestRecord]{ID: 2, Data: TestRecord{Name: "Updated Product", Amount: 200.0}})
	heap.Flush()

	insertedRecord, err := heap.FindByID(1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Product 1", insertedRecord.Data.Name)

	updatedRecord, err := heap.FindByID(2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Updated Product", updatedRecord.Data.Name)

	ClearTestFolder()
}

func TestFindByIdLargeDataset(t *testing.T) {
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 1000; i++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	assert.Equal(t, 1000, len(heap.wal.entries))

	for i := 0; i < 4; i++ {
		randomNum := rand.Uint64() % 1001
		record, _ := heap.FindByID(randomNum)
		assert.Equal(t, fmt.Sprintf("Product %d", randomNum), record.Data.Name)
	}

	//ClearTestFolder()
}

func TestFindByIdWithinIndexesLargeDataset(t *testing.T) {
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 1000; i++ {
		_ = heap.InsertWithIndex(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	assert.Equal(t, 1000, len(heap.wal.entries))

	for i := 0; i < 4; i++ {
		randomNum := rand.Uint64() % 1001
		record, _ := heap.FindByIdWithinIndex(randomNum)
		assert.Equal(t, fmt.Sprintf("Product %d", randomNum), record.Data.Name)
	}

	ClearTestFolder()
}

func TestCanRecoverFromWAL(t *testing.T) {
	heap := NewHeap[TestRecord]("products", "storage/test")
	err := heap.Insert(&Record[TestRecord]{ID: 1, Data: TestRecord{Name: "Product 1", Amount: 100.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 2, Data: TestRecord{Name: "Product 2", Amount: 200.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 3, Data: TestRecord{Name: "Product 3", Amount: 220.0}})
	err = heap.Update(&Record[TestRecord]{ID: 2, Data: TestRecord{Name: "Updated Product", Amount: 120.0}})
	err = heap.Delete(3)

	if err != nil {
		t.Error(err)
	}

	heap.wal.Flush()

	recoveredHeap := NewHeap[TestRecord]("products2", "storage/test")
	recoveredHeap.Recover(heap.wal.path)

	firstRecord, _ := recoveredHeap.FindByID(1)
	updatedRecord, _ := recoveredHeap.FindByID(2)
	deletedRecord, _ := recoveredHeap.FindByID(3)

	assert.Equal(t, 0, len(recoveredHeap.wal.entries))
	assert.Equal(t, "Product 1", firstRecord.Data.Name)
	assert.Equal(t, 100.0, firstRecord.Data.Amount)

	assert.Equal(t, "Updated Product", updatedRecord.Data.Name)
	assert.Equal(t, 120.0, updatedRecord.Data.Amount)

	assert.Nil(t, deletedRecord)

	ClearTestFolder()
}

func BenchmarkFindById(b *testing.B) {
	// create 100 records
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 1000; i++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		randomNum := rand.Uint64() % 101
		_, _ = heap.FindByID(randomNum)

		//assert.Equal(b, fmt.Sprintf("Product %d", randomNum), record.Data.Name)
	}
}

func BenchmarkInsertWithIndex(b *testing.B) {
	heap := NewHeap[TestRecord]("products", "storage/test")

	for n := 0; n < b.N; n++ {
		_ = heap.InsertWithIndex(&Record[TestRecord]{ID: uint64(1), Data: TestRecord{Name: fmt.Sprintf("Product %d", 1), Amount: 100.0}})
	}
}

func BenchmarkInsertWithoutIndex(b *testing.B) {
	heap := NewHeap[TestRecord]("products", "storage/test")

	for n := 0; n < b.N; n++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(1), Data: TestRecord{Name: fmt.Sprintf("Product %d", 1), Amount: 100.0}})
	}
}

func BenchmarkFindByIdWithIndexes(b *testing.B) {

	// create 100 records
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 500000; i++ {
		_ = heap.InsertWithIndex(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		randomNum := rand.Uint64() % 101
		_, _ = heap.FindByIdWithinIndex(randomNum)

		//assert.Equal(b, fmt.Sprintf("Product %d", randomNum), record.Data.Name)
	}
}

func BenchmarkFindAll(b *testing.B) {
	// create 100 records
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 100; i++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_, _ = heap.FindAll()

		//assert.Equal(b, fmt.Sprintf("Product %d", randomNum), record.Data.Name)
	}

	ClearTestFolder()
}

func BenchmarkFlush(b *testing.B) {
	// create 100 records
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 100; i++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		heap.Flush()

		ClearTestFolder()
	}
}
