package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"sync"
)

func startWorkers(
	workersCount int,
	dst string,
	metadata <-chan fileMetadata,
	dirChan chan<- dirReq,
	cache *cache,
) {
	wg := &sync.WaitGroup{}
	wg.Add(workersCount)
	for range workersCount {
		go worker(wg, dst, metadata, dirChan, cache)
	}
	wg.Wait()
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
