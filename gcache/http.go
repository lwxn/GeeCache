package gcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath string = "/_zzh/"

type HTTPPool struct {
	self string
	basePath string
}

func NewHTTPPool(self string)*HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool)Log(format string,v ...interface{}){
	log.Printf("[Server %s] %s",p.self,fmt.Sprintf(format,v...))
}

func (p *HTTPPool)ServeHTTP(w http.ResponseWriter,req *http.Request){
	if !strings.HasPrefix(req.URL.Path,p.basePath){
		panic("HTTPServe serves the error path : " + req.URL.Path)
	}

	//-----/basePath/groupName/key
	parts := strings.SplitN(req.URL.Path[len(p.basePath):],"/",2)
	if len(parts) < 2{
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil{
		http.Error(w,fmt.Sprintf("Group %s not exists ",groupName),http.StatusNotFound)
		return
	}

	view,err := group.Get(key)
	if err != nil{
		http.Error(w,fmt.Sprintf("key %s not exists ",key),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/octet-stream")
	w.Write(view.ByteSlices())

}