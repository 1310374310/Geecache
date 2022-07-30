package geecache

// 只读数据结构，用于表示缓存值
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

// 将数据作为字节切片的形式返回
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 将数据以string形式返回
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	// 返回缓存的拷贝，防止缓存内容被外部程序修改
	copy(c, b)
	return c
}
