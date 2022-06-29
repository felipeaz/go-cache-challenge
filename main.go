package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu sync.Mutex
)

type Cache interface {
	Add(key string, value any, expirationInMS int64) Cache
	Unset(key string) Cache
	List()
}

type cache struct {
	data map[string]any
}

func NewCache() Cache {
	c := &cache{
		data: make(map[string]any),
	}
	return c
}

func (c *cache) Add(key string, value any, exp int64) Cache {
	c.data[key] = value
	go func() {
		select {
		case <-time.After(time.Duration(exp) * time.Millisecond):
			fmt.Printf("%s expired\n", key)
			_, ok := c.data[key]
			if ok {
				c.Unset(key)
			}
			return
		}
	}()
	return c
}

func (c *cache) Unset(key string) Cache {
	delete(c.data, key)
	c.List()
	return c
}

func (c *cache) List() {
	for key, value := range c.data {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	c := NewCache()
	c.Add("first", "first cache value", 500).
		Add("second", "second cache value", 2000).
		Add("third", "second cache value", 100)
	c.List()
	time.Sleep(time.Second * 5)
}
