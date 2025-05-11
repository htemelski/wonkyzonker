package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func extractMetadata(path string, info fs.FileInfo) fileMetadata {
	fileParts := strings.Split(path, "/")
	fileName := fileParts[len(fileParts)-1]

	format := strings.Split(fileName, ".")[1]

	year, month, day := info.ModTime().Date()
	creationDate := fmt.Sprintf("%d-%d-%d", year, month, day)

	return fileMetadata{
		path:         path,
		fileName:     fileName,
		format:       format,
		creationDate: creationDate,
	}
}

func walker(dir string) chan fileMetadata {
	metadataChan := make(chan fileMetadata, chanSize)
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	filesCount := len(files)

	go func() {
		defer close(metadataChan)
		progress := 0
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {

			if info.IsDir() {
				return nil
			}

			metadataChan <- extractMetadata(path, info)
			if progress%100 == 0 {
				fmt.Printf("Processed %d/%d files\n", progress, filesCount)
			}
			progress++
			return nil
		})
		if err != nil {
			panic(err)
		}
	}()

	return metadataChan
}
