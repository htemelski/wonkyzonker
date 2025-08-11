package main

import (
	"fmt"
	"os"
	"path"
)

const (
	workers   = 8
	chanSize  = workers * 16
	filePerms = 0644
	dirPerms  = 0755
)

func main() {
	src := os.Getenv("WONKY_SRC")
	dst := os.Getenv("WONKY_DST")

	if src == "" || dst == "" {
		fmt.Printf(
			"empty env variable WONKY_SRC: \"%s\", WONKY_DST: \"%s\"\n",
			src, dst,
		)
		os.Exit(1)
	}

	dirChan := make(chan dirReq, chanSize)
	go dirCreator(dirChan)

	cache := initCache(path.Join(dst, "cache.bin"))
	startWorkers(workers, dst, walker(src), dirChan, cache)
	close(dirChan)
	cache.save()
}
