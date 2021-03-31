package main

import (
	"fmt"
	"log"
	"net/http"
	"gcache"
)

var db = map[string]string{
	"emt":"amilya",
	"paku":"pak",
	"pandora":"mojuan",
}



func main() {
	gcache.NewGroup("scores",2<<10, gcache.GetterFunc(
		func(key string)([]byte,error){
			log.Println("[slowDB] search key",key)
			if v,ok := db[key];ok{
				return []byte(v),nil
			}else{
				return nil,fmt.Errorf("%s is not in db",key)
			}
		}))

	addr := "localhost:9999"
	peers := gcache.NewHTTPPool(addr)
	log.Fatal(http.ListenAndServe(addr,peers))
}