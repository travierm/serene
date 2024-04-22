package main

import (
	"encoding/json"
	"os"
	"sync"
)

type WALEntry struct {
	Operation string
	RecordID  uint64
	Data      []byte
}

type WAL struct {
	path    string
	entries []*WALEntry
	mutex   sync.Mutex
}

func NewWAL(path string) *WAL {
	return &WAL{
		path:    path,
		entries: make([]*WALEntry, 0),
	}
}

func (w *WAL) Log(entry *WALEntry) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.entries = append(w.entries, entry)
}

func (w *WAL) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	file, err := os.Create(w.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, entry := range w.entries {
		err := encoder.Encode(entry)
		if err != nil {
			return err
		}
	}

	// Clear the WAL entries after flushing
	w.entries = make([]*WALEntry, 0)

	return nil
}
