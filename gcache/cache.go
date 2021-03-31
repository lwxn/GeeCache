package gcache

import (
	"gcache/lru"
	"sync"
)

type cache struct {
	m sync.Mutex

	lru *lru.Cache
	cacheBytes int64
}

func (c *cache)Add(key string,value ByteView){
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil{
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key, value)
}

func (c *cache)Get(key string)(value ByteView,ok bool){
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil{
		return
	}
	if v,ok := c.lru.Get(key);ok{
		return v.(ByteView),ok
	}
	return
}


