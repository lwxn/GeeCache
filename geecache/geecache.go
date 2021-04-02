package geecache

import (
	"fmt"
	"log"
	"sync"
	"geecache/singleflight"
)

//--------------------group的设计（可以理解为一个远程服务器）
type Group struct{
	name string
	getter Getter
	groupcache cache
	peers PeerPicker

	//
	loader *singleflight.Group
}

var(
	rwmu sync.RWMutex
	groups =  make(map[string]*Group)
)

func NewGroup(name string,maxByte int64,getter Getter)*Group{
	if getter == nil{
		panic("The group getter is empty...")
	}
	rwmu.Lock()
	defer rwmu.Unlock()
	g := &Group{
		name: name,
		groupcache: cache{
			MaxByte: maxByte,
		},
		getter: getter,
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string)*Group{
	rwmu.Lock()
	defer rwmu.Unlock()
	if g,ok := groups[name];ok{
		return g
	}
	return nil
}

func (g *Group)RegisterGroup(peers PeerPicker){
	g.peers = peers
}

func (g* Group)load(key string)(value ByteView,err error){
	viewi,err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil{
			if peer,ok := g.peers.PeerPick(key);ok{
				if value,err := g.getFromPeer(peer,key);err == nil{
					return value,nil
				}else{
					log.Println("Fail to get from peer",err)
				}

			}
		}
		return g.getLocally(key)
	})
	if err == nil{
		return viewi.(ByteView),nil
	}
	return
}

func (g* Group)getFromPeer(peer PeerGetter,key string)(ByteView,error){
	bytes,err := peer.Get(g.name,key)
	if err != nil{
		return ByteView{},err
	}
	return ByteView{bytes},nil
}

//根据key从group之中获取value
func (g *Group)Get(key string)(ByteView,error){
	//如果key是空的？
	if key == ""{
		return ByteView{},fmt.Errorf("geecache.go: the key is empty")
	}
	if v,ok := g.groupcache.Get(key);ok{
		return v,nil
	}
	fmt.Println("--------------------not found------------------")
	return g.getLocally(key)
}

//从本地进行查找
func (g *Group)getLocally(key string)(ByteView,error){
	v,err := g.getter.Get(key)
	if err != nil{
		return ByteView{},err
	}else{
		bv := ByteView{
			cloneBytes(v),
		}
		g.addCache(key,bv)
		return bv,nil
	}
}


func (g *Group) addCache(key string, v ByteView){
	g.groupcache.Add(key,v)
}




//----------------------------回调函数的设计
type Getter interface {
	Get(key string)([]byte,error)
}

type Getterfunc func(key string)([]byte,error)


func (f Getterfunc)Get(key string)([]byte,error){
	return f(key)
}


