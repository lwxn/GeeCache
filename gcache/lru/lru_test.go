package lru

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

type String string

func (s String)Len() int {
	return len(s)
}

func TestAdd(t *testing.T){
	c := New(int64(4444),nil)
	c.Add("H", String("d"))
	if v,ok := c.Get("H");!ok && v != String("d"){
		log.Fatal("fail")
	}
	log.Println("success")
	//if _,ok := c.Get("po");!ok{
	//	log.Fatal("fail to get k2")
	//}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1,k2,k3 := "da","feng","zi"
	v1,v2,v3 := "so","a","pity"
	cap := len(k1+k2+v1+v2)

	c := New(int64(cap),nil)
	c.Add(k1, String(v1))
	c.Add(k2, String(v2))
	c.Add(k3, String(v3))

	//if _,ok := c.Get(k1); !ok{
	//	//log.Fatal("can not get k1")
	//}
}


func TestOnEvicted(t *testing.T){
	keys := make([]string,0)
	callback := func(key string,value Value) {
		keys = append(keys, key)
	}

	c := New(int64(10),callback)
	c.Add("k1", String("2"))
	c.Add("k2", String("1"))

	expect := []string{"k1","k2"}

	fmt.Println(keys)
	if !reflect.DeepEqual(expect,keys){
		t.Fatal("call fails")
	}
}