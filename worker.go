package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"sync"
)

func createWorkers(n int, wg *sync.WaitGroup, dst string, metadata <-chan fileMetadata, dirChan chan<- dirReq) {
	wg.Add(n)
	for range n {
		go worker(wg, dst, metadata, dirChan)
	}
}

func worker(wg *sync.WaitGroup, dst string, metadata <-chan fileMetadata, dirChan chan<- dirReq) {
	defer wg.Done()

	for meta := range metadata {
		fileLocation := path.Join(dst, meta.format, meta.creationDate, meta.fileName)
		_, err := os.Stat(fileLocation)
		if errors.Is(err, fs.ErrNotExist) {
			handleNewFile(meta, fileLocation, dirChan)
			continue
		}

		handleExistingFile(fileLocation, meta, 0)
	}
}
