package geecache

type ByteView struct {
	b []byte
}

func (b ByteView)Len() int  {
	return len(b.b)
}

func (b ByteView)String() string{
	return string(b.b)
}

func cloneBytes(b []byte)[]byte{
	t := make([]byte,len(b))
	copy(t,b)
	return t
}