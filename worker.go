package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"sync"
)

func createWorkers(
	workersCount int,
	wg *sync.WaitGroup,
	dst string,
	metadata <-chan fileMetadata,
	dirChan chan<- dirReq,
	cache *cache,
) {
	wg.Add(workersCount)
	for range workersCount {
		go worker(wg, dst, metadata, dirChan, cache)
	}
}

func worker(
	wg *sync.WaitGroup,
	dst string,
	metadata <-chan fileMetadata,
	dirChan chan<- dirReq,
	cache *cache,
) {
	defer wg.Done()

	for meta := range metadata {
		if cache.exists(meta.fileName) {
			continue
		}

		cache.store(meta.fileName)
		fileLocation := path.Join(dst, meta.format, meta.creationDate, meta.fileName)
		_, err := os.Stat(fileLocation)
		if errors.Is(err, fs.ErrNotExist) {
			handleNewFile(meta, fileLocation, dirChan)
			continue
		}

		handleExistingFile(fileLocation, meta, 0)
	}
}
