package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
)

// Create a new csv.Writer for the given file
func create_csv(folder string, filename string) *csv.Writer {
	f, err := os.Create(filepath.Join(folder, filename))
	if err != nil {
		log.Fatal("Failed creating file", filename, "in folder", folder)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'

	return w
}
