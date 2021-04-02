package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte)uint32

type Map struct {
	hash Hash  //hash算法
	replicas int //副本的倍数
	keys []int   //所有的节点的hash值
	hashMap map[int]string    //hash值对应的真实的节点值
}

func New(replicas int,fn Hash)*Map{
	m := &Map{
		replicas: replicas,
		hashMap: make(map[int]string),
	}
	if fn == nil{
		m.hash = crc32.ChecksumIEEE
	}else{
		m.hash = fn
	}
	return m
}

func (m *Map)Add(keys ...string){
	for _,key := range keys{
		for i:=0;i<m.replicas;i++{
			value := int(m.hash([]byte(key + strconv.Itoa(i))))
			m.keys = append(m.keys, value)
			m.hashMap[value] = key
		}
	}
	sort.Ints(m.keys)
}

//根据value来找到节点
func (m *Map)Get(key string)string{
	//如果目前远程服务器都是空的
	if m.keys == nil{
		return ""
	}
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= int(m.hash([]byte(key)))
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}