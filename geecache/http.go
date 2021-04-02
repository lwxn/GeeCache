package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//--------------------
type httpGetter struct{
 	baseURL string
}

func (h *httpGetter)Get(group string,key string)([]byte,error)  {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}


//---------------------------------HTTPPOOL
const (
	defaultBasePath = "/_lwxnlovezzhlp/"
	defaultReplicas = 50
)


type HTTPPOOL struct{
	self string
	basePath string

	mu sync.Mutex
	peers *consistenthash.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPOOL (self string)*HTTPPOOL{
	return &HTTPPOOL{
		self: self,
		basePath: defaultBasePath,
	}
}

//将服务器节点插入
func (h *HTTPPOOL)Set(peers ...string){
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaultReplicas,nil)
	h.peers.Add(peers...)

	h.httpGetters = make(map[string]*httpGetter,len(peers))
	for _,peer := range peers{
		h.httpGetters[peer] = &httpGetter{
			baseURL: peer+h.basePath,
		}
	}
}

//根据key来选择正确的服务器
func (h *HTTPPOOL)PeerPick(key string)(PeerGetter,bool){
	h.mu.Lock()
	defer h.mu.Unlock()

	if peer := h.peers.Get(key);peer != "" || peer != h.self{
		return h.httpGetters[peer],true
	}
	return nil,false
}

var _ PeerPicker = (*HTTPPOOL)(nil)
var _ PeerGetter = (*httpGetter)(nil)

//处理方法     handle url such as /_lwxnlovezzhlp/lp/zzh
func (h *HTTPPOOL)ServeHTTP(w http.ResponseWriter,req *http.Request){
	if strings.HasPrefix(req.URL.Path,h.basePath) != true{
		panic("It is a wrong path url...")
	}

	pathSplit := strings.SplitN(req.URL.Path[len(h.basePath):],"/",2)
	groupname := pathSplit[0]
	key := pathSplit[1]

	group := GetGroup(groupname)
	if group == nil{
		http.Error(w,"no such group name: " + groupname,http.StatusNotFound)
		return
	}
	value,err := group.Get(key)
	if err != nil{
		http.Error(w,fmt.Sprintf("group %s can't find %s\n", groupname,key),http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(cloneBytes(value.b))
}
