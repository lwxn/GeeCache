package gcache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string,cacheBytes int64,getter Getter)*Group {
	if getter == nil{
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string)*Group {
	mu.RLock()
	defer mu.RUnlock()
	g,ok := groups[name]
	if ok{
		return g;
	}else{
		return nil
	}

}

//if I have the key in the mainCache
func (g *Group)Get(key string)(ByteView,error){
	if key == ""{
		return ByteView{},fmt.Errorf("empty key!")
	}else{
		if v,ok := g.mainCache.Get(key);ok{
			log.Printf("hit GeeCache")
			return v,nil
		}
		return g.load(key)
	}
}

func (g *Group)load(key string)(ByteView,error){
	return g.getLocally(key)
}

func (g *Group)getLocally(key string)(ByteView,error){
	v,err := g.getter.Get(key)
	if err != nil{
		return ByteView{}, err
	}

	value := ByteView{
		b: CloneByte(v),
	}
	g.populateCache(key,value)
	return value,nil;
}

func (g *Group)populateCache(key string,value ByteView){
	g.mainCache.Add(key,value)
}

type Getter interface {
	Get(key string)([]byte,error)
}

type GetterFunc func(key string)([]byte,error)

func(f GetterFunc) Get(key string)([]byte,error){
	return f(key)
}
