package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sync"
)

type dirReq struct {
	path string
	wg   *sync.WaitGroup
}

func dirCreator(dirRequests <-chan dirReq) {
	createdDirs := map[string]struct{}{}

	for dirRequest := range dirRequests {
		func() {
			defer dirRequest.wg.Done()

			dir := path.Dir(dirRequest.path)
			if _, ok := createdDirs[dir]; ok {
				return
			}

			parentDir := path.Dir(dir)

			if _, err := os.Stat(parentDir); errors.Is(err, fs.ErrNotExist) {
				err := os.Mkdir(parentDir, dirPerms)
				if err != nil {
					fmt.Println("failed to make parent dir", parentDir, err.Error())
				}
			}

			if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
				err := os.Mkdir(dir, dirPerms)
				if err != nil {
					fmt.Println("failed to make sub dir", dir, err.Error())
				}
			}
			createdDirs[dir] = struct{}{}
		}()
	}
}
