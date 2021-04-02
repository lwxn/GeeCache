package geecache

import (
	"geecache/lru"
	"sync"
)

//add mutex lock to cache
type cache struct {
	mu sync.Mutex
	c *lru.Cache
	MaxByte int64
}

func (c *cache)Add(key string,value ByteView)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.c == nil{
		c.c = lru.New(c.MaxByte,nil)
	}
	c.c.Add(key,value)
}

func (c *cache)Get(key string)(ByteView,bool){
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.c == nil{
		return ByteView{},false
	}

	v,ok := c.c.Get(key)
	if ok{
		return v.(ByteView),ok
	}
	return ByteView{},ok
}