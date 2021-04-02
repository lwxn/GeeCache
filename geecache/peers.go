package geecache

//选取服务器
type PeerPicker interface {
	PeerPick(key string)(PeerGetter,bool)
}

//服务器获取值
type PeerGetter interface {
	Get(group string,key string)([]byte,error)
}
