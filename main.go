package main

import (
	"sync"
)

const src = "/run/media/hawk/EOS_DIGITAL/DCIM/100CANON"
const dst = "/home/hawk/data/photos"
const cachePath = "/home/hawk/data/photos/cache.bin"
const workers = 8

const (
	chanSize  = workers * 16
	filePerms = 0644
	dirPerms  = 0755
)

func main() {
	wg := &sync.WaitGroup{}
	dirChan := make(chan dirReq, chanSize)
	go dirCreator(dirChan)

	cache := initCache()
	createWorkers(workers, wg, dst, walker(src), dirChan, cache)
	wg.Wait()
	close(dirChan)
	cache.save()
}
