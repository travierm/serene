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

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func BytesToUint64(b []byte) uint64 {
	var result uint64
	for i := 0; i < len(b) && i < 8; i++ {
		result = result<<8 | uint64(b[i])
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
