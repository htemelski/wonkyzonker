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
	if len(fileParts) >= 3 {
		fileName = fmt.Sprintf("%s-%s", fileParts[len(fileParts)-2], fileName)
	}

	fileNameParts := strings.Split(fileName, ".")
	format := "misc"
	if len(fileNameParts) == 2 {
		format = fileNameParts[1]
	}

	year, month, day := info.ModTime().Date()
	creationDate := fmt.Sprintf("%d-%d-%d", year, month, day)

	return fileMetadata{
		path:         path,
		fileName:     fileName,
		format:       format,
		creationDate: creationDate,
	}
}

func walker(rootDir string) chan fileMetadata {
	metadataChan := make(chan fileMetadata, chanSize)
	filesCount := 0

	go func() {
		defer close(metadataChan)
		progress := 0
		err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				if path != rootDir {
					filesCount -= 1
				}

				count, err := filesCounter(path)
				if err != nil {
					return err
				}
				filesCount += count
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
			fmt.Printf("failed walk dir %s\n", err.Error())
			os.Exit(1)
		}
	}()

	return metadataChan
}

func filesCounter(dir string) (int, error) {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	files, err := f.Readdirnames(-1)
	if err != nil {
		fmt.Printf("failed read dirnames %s\n", err.Error())
		return 0, err
	}

	return len(files), nil
}
