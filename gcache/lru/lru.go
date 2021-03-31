package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nBytes int64
	ll *list.List
	cache map[string]*list.Element
	OnEvicted func(key string,value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64,onEvicted func(string, Value))*Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//get key
func (c *Cache)Get(key string)(value Value,ok bool){
	if ele,ok := c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value,true
	}
	return
}

//移出旧的缓存
func (c*Cache)RemoveOldest(){
	ele := c.ll.Back()
	if ele != nil{
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil{
			c.OnEvicted(kv.key,kv.value)
		}
	}
}

//增加cache
func (c *Cache)Add(key string,value Value){
	if ele,ok := c.cache[key];ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}else{
		ele := c.ll.PushFront(&entry{key: key,value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes{
		c.RemoveOldest()
	}
}

func (c *Cache)len()int{
	return c.ll.Len()
}
