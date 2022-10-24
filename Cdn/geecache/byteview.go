package geecache

/*缓存值的抽象与封装。ByteView只读数据结构，表示缓存值*/

// ByteView是一个只读结构，表示存储真实的缓存值，byte类型能支持任意的数据类型的存储，如字符串、图片等
// A ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte
}

// Len 返回ByteView.b的长度
// Len return the view's length
// LRU.Cache实现中要求缓存对象必须实现Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回ByteView.b的一个拷贝，因为b是只读的，防止缓存值被外部程序修改
// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 将数据作为字符串返回，必要时进行复制
// String return the data as a string,making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 对缓存值进行复制
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
