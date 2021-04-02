package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String)Len() int {
	return len(d)
}

func TestCache_Add(t *testing.T) {
	c := New(int64(20),nil)
	c.Add("zzh", String("lp"))
	if v,ok := c.Get("zzh");!ok || string(v.(String)) != "lp"{
		t.Fatal("cache fail to get:","zzh")
	}
}

func TestCache_RemoveOldest(t *testing.T) {

}

func TestOnEvicted(t *testing.T){
	keys := []string{}

	c := New(int64(10), func(key string, value Value) {
		keys  = append(keys,key)
	})

	c.Add("key1", String("123456"))
	c.Add("EMT", String("leimu"))

	fmt.Println(keys)
}
