package gcache

type ByteView struct {
	b []byte
}

func (v ByteView)Len()int{
	return len(v.b)
}

func (v ByteView)ByteSlices()[]byte{
	return CloneByte(v.b)
}

func (v ByteView)String()string{
	return string(v.b)
}

func CloneByte(b1 []byte)[]byte{
	b := make([]byte,len(b1))
	copy(b,b1)
	return b
}