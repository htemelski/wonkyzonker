package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"sync"
)

type fileMetadata struct {
	path         string
	fileName     string
	format       string
	creationDate string
}

type file []byte

func (f file) sha256() string {
	h := sha256.New()
	h.Write(f)
	return string(h.Sum(nil))
}

func (f file) write(dst string) error {
	return os.WriteFile(dst, f, filePerms)
}

func read(src string) (file, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func handleExistingFile(fileLocation string, meta fileMetadata, iterations int) {
	ogLocation := fileLocation

	if iterations != 0 {
		lastDotIndex := strings.LastIndex(fileLocation, ".")
		firstPart := fileLocation[:lastDotIndex]
		lastPart := fileLocation[lastDotIndex+1:]
		fileLocation = fmt.Sprintf("%s-%d.%s", firstPart, iterations, lastPart)
	}

	newFile, err := read(meta.path)
	if err != nil {
		fmt.Printf("failed to read new file %s, error: %s\n", meta.path, err.Error())
		return
	}

	existingFile, err := read(fileLocation)
	if err != nil {
		err = newFile.write(fileLocation)
		if err != nil {
			fmt.Printf("failed to write new file %s, error: %s", fileLocation, err.Error())
		}
		return
	}

	if existingFile.sha256() == newFile.sha256() {
		return
	}

	handleExistingFile(ogLocation, meta, iterations+1)
}

func handleNewFile(meta fileMetadata, fileLocation string, dirChan chan<- dirReq) {
	file, err := read(meta.path)
	if err != nil {
		fmt.Printf("failed to read file %s, error: %s\n", meta.fileName, err.Error())
		return
	}

	dirWg := &sync.WaitGroup{}
	dirWg.Add(1)
	dirReq := dirReq{
		path: fileLocation,
		wg:   dirWg,
	}
	dirChan <- dirReq
	dirWg.Wait()

	err = file.write(fileLocation)
	if err != nil {
		fmt.Printf("failed to write file %s, error: %s\n", fileLocation, err.Error())
	}
}
