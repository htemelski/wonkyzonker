package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
)

type cacheData map[string]struct{}
type cacheCmd func(cacheData)

type cache struct {
	cmds chan cacheCmd
	path string
}

func initCache(path string) *cache {
	cache := &cache{
		cmds: make(chan cacheCmd, chanSize),
		path: path,
	}

	data := cacheData{}
	fileData, err := os.ReadFile(cache.path)
	if err == nil {
		_ = gob.NewDecoder(bytes.NewReader(fileData)).Decode(&data)
	}

	go func() {
		for cmd := range cache.cmds {
			cmd(data)
		}
	}()

	return cache
}

func (c *cache) exists(key string) bool {
	res := make(chan bool)

	c.cmds <- func(cd cacheData) {
		_, ok := cd[key]
		res <- ok
	}

	return <-res
}

func (c *cache) store(key string) {
	shouldSave := make(chan bool)
	c.cmds <- func(cd cacheData) {
		cd[key] = struct{}{}

		shouldSave <- (len(cd)%100 == 0)
	}

	if <-shouldSave {
		c.save()
	}
}

func (c *cache) save() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	c.cmds <- func(cd cacheData) {
		defer wg.Done()
		buf := &bytes.Buffer{}
		_ = gob.NewEncoder(buf).Encode(cd)
		_ = os.WriteFile(c.path, buf.Bytes(), filePerms)
	}

	wg.Wait()
}
