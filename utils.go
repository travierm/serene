package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveKey[T any](slice []T, key T, equals func(T, T) bool) []T {
	for i, v := range slice {
		if equals(v, key) {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func Filter[T any](s []T, fn func(T) bool) []T {
	var result []T
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func ClearTestFolder() {
	folder := "storage/test"

	// Delete .bin files
	binFiles, err := filepath.Glob(filepath.Join(folder, "*.bin"))
	if err != nil {
		fmt.Printf("Error finding .bin files in '%s': %v\n", folder, err)
		return
	}
	for _, file := range binFiles {
		err = os.Remove(file)
		if err != nil {
			fmt.Printf("Error deleting file '%s': %v\n", file, err)
		}
	}

	// Delete .wal files
	walFiles, err := filepath.Glob(filepath.Join(folder, "*.wal"))
	if err != nil {
		fmt.Printf("Error finding .wal files in '%s': %v\n", folder, err)
		return
	}
	for _, file := range walFiles {
		err = os.Remove(file)
		if err != nil {
			fmt.Printf("Error deleting file '%s': %v\n", file, err)
		}
	}
}
