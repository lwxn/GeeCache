package gcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"zzh":"90",
	"zzh1p":"100",
	"zzhmv":"100",
}

func TestCache_Get(t *testing.T) {
	loadCount := make(map[string]int,len(db))
	gee := NewGroup("cache",2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slowDB] search key",key)
		if v, ok := db[key];ok{
			loadCount[key]++;
			return []byte(v),nil
		}
		return nil,fmt.Errorf("%s is not in the bd",key)
	}))

	for k,v := range db {
		if view,err := gee.Get(k);err != nil || view.String() != v{
			log.Fatal("failed to get ",view)
		}
		if _, err := gee.Get(k);err != nil || loadCount[k] > 1{
			log.Fatal("fail to hit cache : ",k)
		}
	}

	if view,err := gee.Get("unknown");err == nil{
		log.Fatalf("%s shouldn't be here.",view)
	}
}


func TestGetterFunc_Get(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key),nil
	})

	expect := []byte("zzh")
	if v,ok := f.Get("zzh");ok != nil || !reflect.DeepEqual(v,expect){
		log.Fatal("fail")
	}
}
