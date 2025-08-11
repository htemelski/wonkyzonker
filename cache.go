package main

import (
	"encoding/json"
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
		cmds: make(chan cacheCmd, workers*16),
		path: path,
	}

	data := cacheData{}
	fileData, err := os.ReadFile(cache.path)
	if err == nil {
		_ = json.Unmarshal(fileData, &data)
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
	c.cmds <- func(cd cacheData) {
		cd[key] = struct{}{}
	}
}

func (c *cache) save() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	c.cmds <- func(cd cacheData) {
		defer wg.Done()

		data, _ := json.Marshal(cd)
		_ = os.WriteFile(c.path, data, 0644)
	}

	wg.Wait()
}
