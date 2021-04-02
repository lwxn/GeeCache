package lru

import "container/list"

type Cache struct {
	maxByte int64   //最大缓存值
	nByte int64		//当前缓存值
	l *list.List    //双向链表
 	cache map[string]*list.Element     //哈希
	onEnvicted func(key string,value Value)
}

type node struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

//构造函数
func New(maxByte int64,onEvicted func(key string,value Value)) *Cache {
	return &Cache{
		maxByte: maxByte,
		l : list.New(),
		cache: make(map[string]*list.Element),
		onEnvicted: onEvicted,
	}
}

//需要使用一个节点
func (c *Cache)Get(key string)(Value,bool){
	//如果存在的话，移动到最前面
	if e,ok := c.cache[key];ok{
		c.l.MoveToFront(e)
		kv := e.Value.(*node)
		return kv.value,true
	}
	return nil,false

}

//添加新的一条进去
func (c *Cache)Add(key string,value Value){
	// 如果在cache之中已经有这一条了
	if e,ok := c.cache[key];ok{
		kv := e.Value.(*node)
		c.nByte += int64(value.Len()) -int64(kv.value.Len())
		c.Get(key)
		kv.value = value
	}else{  //否则就临时加进去
		c.l.PushFront(&node{key: key,value: value})
		c.cache[key] = c.l.Front()
		c.nByte += int64(len(key)) + int64(value.Len())
	}
	for c.nByte > c.maxByte{
		c.RemoveOldest()
	}
}

//删除缓存，删除最少用的那个节点
func (c *Cache)RemoveOldest()  {
	e := c.l.Back()
	if e != nil{
		c.l.Remove(e)
		kv := e.Value.(*node)
		c.nByte -= int64(kv.value.Len()) + int64(len(kv.key))
		delete(c.cache,kv.key)
		if c.onEnvicted != nil{
			c.onEnvicted(kv.key,kv.value)
		}
	}

}

//获取缓存的长度
func (c *Cache)Len() int  {
	return c.l.Len()
}



