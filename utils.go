package main

import (
	"fmt"
	"os"
)

func ClearTestFolder(folder string) {
	err := os.RemoveAll(folder)
	if err != nil {
		fmt.Printf("Error deleting folder '%s': %v\n", folder, err)
	}
}
