package main

import (
	"encoding/json"
	"os"
	"sync"
)

type cacheData map[string]struct{}
type cacheCmd func(*cacheData)

type cache struct {
	cmdChan chan cacheCmd
}

func initCache() *cache {
	cache := &cache{
		cmdChan: make(chan cacheCmd, 16),
	}

	data := &cacheData{}
	fileData, err := os.ReadFile(cachePath)
	if err == nil {
		_ = json.Unmarshal(fileData, data)
	}

	go func() {
		for cmd := range cache.cmdChan {
			cmd(data)
		}
	}()

	return cache
}

func (c *cache) exists(key string) bool {
	res := make(chan bool)

	c.cmdChan <- func(cd *cacheData) {
		_, ok := (*cd)[key]
		res <- ok
	}

	return <-res
}

func (c *cache) store(key string) {
	c.cmdChan <- func(cd *cacheData) {
		(*cd)[key] = struct{}{}
	}
}

func (c *cache) save() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	c.cmdChan <- func(cd *cacheData) {
		defer wg.Done()

		data, _ := json.Marshal(&cd)
		_ = os.WriteFile(cachePath, data, 0644)
	}

	wg.Wait()
}
