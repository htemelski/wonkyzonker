package main

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
	dirChan := make(chan dirReq, chanSize)
	go dirCreator(dirChan)

	cache := initCache(cachePath)
	startWorkers(workers, dst, walker(src), dirChan, cache)
	close(dirChan)
	cache.save()
}
