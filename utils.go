package main

import (
	"fmt"
	"os"
)

func ClearTestFolder() {
	folder := "storage/test"
	err := os.RemoveAll(folder)
	if err != nil {
		fmt.Printf("Error deleting folder '%s': %v\n", folder, err)
	}
}
