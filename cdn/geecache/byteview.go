package geecache

/*缓存值的抽象与封装。ByteView只读数据结构，表示缓存值*/

// A ByteView holds an immutable view of bytes.
// b存储真实的缓存值，byte类型能支持任意的数据类型的存储，如字符串、图片等
type ByteView struct {
	b []byte
}

// Len return the view's length
// LRU.Cache实现中要求缓存对象必须实现Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
// b是只读的，使用ByteSlice()方法返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String return the data as a string,making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
