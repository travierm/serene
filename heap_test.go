package main

import (
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

	ClearTestFolder("storage/test")
}
