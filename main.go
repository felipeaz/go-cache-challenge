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

type data struct {
	key        string
	data       any
	expiration int64
	createdAt  int64
}

type cache struct {
	data []data
}

func NewCache() Cache {
	c := &cache{}
	go c.sync()
	return c
}

func (c *cache) sync() {
	for {
		for _, v := range c.data {
			now := time.Now().UnixNano() / int64(time.Millisecond)
			if (now - v.createdAt) >= v.expiration {
				fmt.Printf("removing %s due to expiration time.\n\n", v.key)
				mu.Lock()
				c.Unset(v.key)
				mu.Unlock()
			}
		}
	}
}

func (c *cache) Add(key string, value any, exp int64) Cache {
	c.data = append(c.data, data{
		key:        key,
		data:       value,
		expiration: exp,
		createdAt:  time.Now().UnixNano() / int64(time.Millisecond),
	})
	return c
}

func (c *cache) Unset(key string) Cache {
	var newData []data
	for _, v := range c.data {
		if v.key != key {
			newData = append(newData, v)
		}
	}
	c.data = newData
	c.List()
	return c
}

func (c *cache) List() {
	for _, d := range c.data {
		fmt.Printf("%s: %s\n", d.key, d.data)
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
