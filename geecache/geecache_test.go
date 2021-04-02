package geecache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"emt":"amily",
	"lemu":"eeee",
	"lamu":"Rozwaer",
}

func TestGetterfunc_Get(t *testing.T) {
	loadCounts := make(map[string]int,len(db))
	gee := NewGroup("boki",1024,Getterfunc(func(key string) ([]byte, error) {
		log.Println("begin slow cache hit...")
		if v,ok := db[key];ok{
			if _,ok := loadCounts[key];!ok{
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v),nil
		}else{
			return nil,fmt.Errorf("%s is not fount in db\n",key)
		}
	}))
	for k,v := range db{
		value,err := gee.Get(k)
		fmt.Printf("%s : %s\n",k,value.String())
		if err != nil || value.String() != v{

			t.Fatal("fail",k,value.String())
		}
		_,err = gee.Get(k)
		if err != nil || loadCounts[k] > 1{
			t.Fatal("You don't use the envicted function!")
		}
	}
}
