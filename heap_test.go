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

func TestHeap(t *testing.T) {
	heap := NewHeap[TestRecord]("products", "storage/test")

	err := heap.Insert(&Record[TestRecord]{ID: 1, Data: TestRecord{Name: "Product 1", Amount: 100.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 2, Data: TestRecord{Name: "Product 2", Amount: 150.0}})
	err = heap.Insert(&Record[TestRecord]{ID: 3, Data: TestRecord{Name: "Product 3", Amount: 150.0}})
	heap.Flush()

	newHeap := NewHeap[TestRecord]("products", "storage/test")
	newHeap.Fill()
	record, err := newHeap.FindByID(2)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Product 2", record.Data.Name)

	ClearTestFolder()
}

func BenchmarkFindById(b *testing.B) {
	// create 100 records
	heap := NewHeap[TestRecord]("products", "storage/test")
	for i := 0; i < 100; i++ {
		_ = heap.Insert(&Record[TestRecord]{ID: uint64(i), Data: TestRecord{Name: fmt.Sprintf("Product %d", i), Amount: 100.0}})
	}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		randomNum := rand.Uint64() % 101
		_, _ = heap.FindByID(randomNum)

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

	//ClearTestFolder()
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
