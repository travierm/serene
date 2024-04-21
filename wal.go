package main

import (
	"encoding/json"
	"os"
)

type WALEntry struct {
	Data []byte
}

type WAL struct {
	file *os.File
}

func (w *WAL) AppendEntry(data []byte) error {
	entry := WALEntry{Data: data}
	encodedEntry, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = w.file.Write(append(encodedEntry, '\n'))
	return err
}
