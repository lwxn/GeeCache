package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte)uint32

type Map struct {
	hash Hash      //hash函数
	replicas int   //虚拟节点倍数
	keys []int
	hashMap map[int]string    //虚假节点--->真实节点
}


func New(replicas int,fn Hash) *Map {
	m := &Map{
		hash: fn,
		replicas: replicas,
		hashMap: make(map[int]string),
	}

	//如果hash的算法是空的话
	if m.hash == nil{
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//加入真实节点
func (m *Map) Add (keys ...string) {
	for _,key := range keys{
		for i := 0;i < m.replicas;i++{
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.hashMap[hash] = key
			m.keys = append(m.keys, hash)
		}
	}
	sort.Ints(m.keys) //排序
}

//获取离该key最近的节点
func (m *Map)Get(key string)string{
	//如果是空的
	if len(m.keys) == 0{
		return ""
	}

	hash := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool{
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[index%len(m.keys)]]  //如果找不到的话，返回的是n,取mode为0
}


