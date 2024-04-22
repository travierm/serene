package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	pageSize = 8192
)

type Record[T any] struct {
	ID   uint64
	Data T
}

type Heap[T any] struct {
	tableName string
	pageSize  int
	pages     []*Page
	mutex     sync.Mutex
	dataDir   string
	wal       *WAL
}

type Page struct {
	data []byte
}

func NewHeap[T any](tableName string, dataDir string) *Heap[T] {
	return &Heap[T]{
		tableName: tableName,
		pageSize:  pageSize,
		pages:     make([]*Page, 0),
		dataDir:   dataDir,
		wal:       NewWAL(fmt.Sprintf("%s/%s.wal", dataDir, tableName)),
	}
}

func (h *Heap[T]) Insert(record *Record[T], disableWAL ...bool) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Serialize the record into bytes
	var recordBytes []byte
	recordBytes = append(recordBytes, make([]byte, 8)...)
	binary.BigEndian.PutUint64(recordBytes[0:8], record.ID)

	dataBytes, err := json.Marshal(record.Data)
	if err != nil {
		return err
	}

	if len(disableWAL) == 0 || !disableWAL[0] {
		h.wal.Log(&WALEntry{
			Operation: "INSERT",
			RecordID:  record.ID,
			Data:      dataBytes,
		})
	}

	recordBytes = append(recordBytes, make([]byte, 4)...)
	binary.BigEndian.PutUint32(recordBytes[8:12], uint32(len(dataBytes)))
	recordBytes = append(recordBytes, dataBytes...)

	// Find a page with enough free space
	page := h.findPageWithSpace(len(recordBytes))
	if page == nil {
		// Create a new page if no suitable page found
		page = &Page{data: make([]byte, 0, h.pageSize)}
		h.pages = append(h.pages, page)
	}

	// Append the serialized record to the page
	page.data = append(page.data, recordBytes...)

	return nil
}

func (h *Heap[T]) Update(record *Record[T], disableWAL ...bool) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, page := range h.pages {
		offset := 0
		for offset < len(page.data) {
			recordID := binary.BigEndian.Uint64(page.data[offset : offset+8])
			dataSize := binary.BigEndian.Uint32(page.data[offset+8 : offset+12])

			if recordID == record.ID {
				// Serialize the updated record data
				dataBytes, err := json.Marshal(record.Data)
				if err != nil {
					return err
				}

				if len(disableWAL) == 0 || !disableWAL[0] {
					h.wal.Log(&WALEntry{
						Operation: "UPDATE",
						RecordID:  record.ID,
						Data:      dataBytes,
					})
				}

				// Check if the updated data size is different from the original size
				if len(dataBytes) != int(dataSize) {
					// Remove the original record
					copy(page.data[offset:], page.data[offset+12+int(dataSize):])
					page.data = page.data[:len(page.data)-12-int(dataSize)]

					// Insert the updated record as a new record
					updatedRecordBytes := make([]byte, 12+len(dataBytes))
					binary.BigEndian.PutUint64(updatedRecordBytes[0:8], record.ID)
					binary.BigEndian.PutUint32(updatedRecordBytes[8:12], uint32(len(dataBytes)))
					copy(updatedRecordBytes[12:], dataBytes)

					page.data = append(page.data, updatedRecordBytes...)
				} else {
					// Update the record data in-place
					copy(page.data[offset+12:offset+12+int(dataSize)], dataBytes)
				}

				return nil
			}

			offset += 12 + int(dataSize)
		}
	}

	return fmt.Errorf("record with ID %d not found", record.ID)
}

func (h *Heap[T]) Delete(id uint64, disableWAL ...bool) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, page := range h.pages {
		offset := 0
		for offset < len(page.data) {
			recordID := binary.BigEndian.Uint64(page.data[offset : offset+8])
			dataSize := binary.BigEndian.Uint32(page.data[offset+8 : offset+12])

			if recordID == id {
				if len(disableWAL) == 0 || !disableWAL[0] {
					h.wal.Log(&WALEntry{
						Operation: "DELETE",
						RecordID:  id,
					})
				}

				// Remove the record by shifting the remaining data to the left
				copy(page.data[offset:], page.data[offset+12+int(dataSize):])
				page.data = page.data[:len(page.data)-12-int(dataSize)]

				return nil
			}

			offset += 12 + int(dataSize)
		}
	}

	return fmt.Errorf("record with ID %d not found", id)
}

func (h *Heap[T]) FindByID(id uint64) (*Record[T], error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, page := range h.pages {
		offset := 0
		for offset < len(page.data) {
			recordID := binary.BigEndian.Uint64(page.data[offset : offset+8])
			dataSize := binary.BigEndian.Uint32(page.data[offset+8 : offset+12])

			if recordID == id {
				dataBytes := page.data[offset+12 : offset+12+int(dataSize)]

				var data T
				err := json.Unmarshal(dataBytes, &data)
				if err != nil {
					return nil, err
				}

				record := &Record[T]{
					ID:   recordID,
					Data: data,
				}
				return record, nil
			}

			offset += 12 + int(dataSize)
		}
	}

	return nil, fmt.Errorf("record with ID %d not found", id)
}

func (h *Heap[T]) FindAll() ([]*Record[T], error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var records []*Record[T]

	for _, page := range h.pages {
		offset := 0
		for offset < len(page.data) {
			recordID := binary.BigEndian.Uint64(page.data[offset : offset+8])
			dataSize := binary.BigEndian.Uint32(page.data[offset+8 : offset+12])
			dataBytes := page.data[offset+12 : offset+12+int(dataSize)]

			var data T
			err := json.Unmarshal(dataBytes, &data)
			if err != nil {
				return nil, err
			}

			record := &Record[T]{
				ID:   recordID,
				Data: data,
			}
			records = append(records, record)

			offset += 12 + int(dataSize)
		}
	}

	return records, nil
}

func (h *Heap[T]) findPageWithSpace(size int) *Page {
	for _, page := range h.pages {
		if len(page.data)+size <= h.pageSize {
			return page
		}
	}
	return nil
}

func (h *Heap[T]) Vacuum() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Create a new set of pages
	newPages := make([]*Page, 0)

	// Rewrite the data to new pages
	for _, page := range h.pages {
		if len(page.data) > 0 {
			newPage := &Page{data: make([]byte, len(page.data))}
			copy(newPage.data, page.data)
			newPages = append(newPages, newPage)
		}
	}

	// Replace the old pages with the new pages
	h.pages = newPages

	return nil
}

func (h *Heap[T]) Flush() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Create the data directory if it doesn't exist
	err := os.MkdirAll(h.dataDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Write each page to a separate file within the data directory
	for i, page := range h.pages {
		filename := filepath.Join(h.dataDir, fmt.Sprintf("%s_page%d.bin", h.tableName, i))
		err := writePageToFile(filename, page)
		if err != nil {
			return err
		}
	}

	h.wal.Flush()

	return nil
}

func writePageToFile(filename string, page *Page) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the page size as a header
	headerBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerBytes, uint32(len(page.data)))
	_, err = file.Write(headerBytes)
	if err != nil {
		return err
	}

	// Write the page data
	_, err = file.Write(page.data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Heap[T]) Fill() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Clear the existing pages
	h.pages = make([]*Page, 0)

	// Scan the data directory for page files
	files, err := os.ReadDir(h.dataDir)
	if err != nil {
		return err
	}

	// Read each page file and reconstruct the pages
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bin" {
			filename := filepath.Join(h.dataDir, file.Name())
			page, err := readPageFromFile(filename)
			if err != nil {
				return err
			}
			h.pages = append(h.pages, page)
		}
	}

	return nil
}

func readPageFromFile(filename string) (*Page, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the page size from the header
	headerBytes := make([]byte, 4)
	_, err = file.Read(headerBytes)
	if err != nil {
		return nil, err
	}
	pageSize := binary.BigEndian.Uint32(headerBytes)

	// Read the page data
	pageData := make([]byte, pageSize)
	_, err = file.Read(pageData)
	if err != nil {
		return nil, err
	}

	page := &Page{data: pageData}
	return page, nil
}

func (h *Heap[T]) Recover(walFilePath string) error {
	file, err := os.Open(walFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var entry WALEntry
		err := decoder.Decode(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch entry.Operation {
		case "INSERT", "UPDATE":
			var data T
			err := json.Unmarshal(entry.Data, &data)
			if err != nil {
				return err
			}
			record := &Record[T]{
				ID:   entry.RecordID,
				Data: data,
			}
			if entry.Operation == "INSERT" {
				err = h.Insert(record, true)
			} else {
				err = h.Update(record, true)
			}

			if err != nil {
				return err
			}
		case "DELETE":
			err = h.Delete(entry.RecordID, true)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
