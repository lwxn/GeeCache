package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db =  map[string]string{
	"emt":"amily",
	"lemu":"486",
	"lamu":"Rozwaer",
}

func createGroup()*geecache.Group{
	return geecache.NewGroup("lp",2<<10,geecache.Getterfunc(
		func(key string) ([]byte, error) {
			log.Println("slowly cache begins...")
			if v,ok := db[key];ok{
				return []byte(v),nil
			}
			return nil,fmt.Errorf("Not found key %s",key)
		}))
}

func startCacheServer(addr string,addrs []string,gee *geecache.Group){
	peers := geecache.NewHTTPPOOL(addr)
	peers.Set(addrs...)

	gee.RegisterGroup(peers)
	log.Println("geecache is running at",addr)
	log.Fatal(http.ListenAndServe(addr[7:],peers))
}

func startAPIServer(apiAddr string,gee *geecache.Group){
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte(view.String()))

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}


func main() {
	var port int
	var api bool
	flag.IntVar(&port,"port",8001,"Geecache server port")
	flag.BoolVar(&api,"api",false,"start a pi server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001 : "http://localhost:8001",
		8002 : "http://localhost:8002",
		8003 : "http://localhost:8003",
	}

	var addrs []string
	for _,v := range addrMap{
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api{
		go startAPIServer(apiAddr,gee)
	}
	startCacheServer(addrMap[port],[]string(addrs),gee)


	addr := "localhost:8080"
	handler := geecache.NewHTTPPOOL(addr)
	log.Println("Now the server is running at ",addr)
	http.ListenAndServe(addr,handler)
}



