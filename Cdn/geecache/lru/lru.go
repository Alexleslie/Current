package lru

import (
	"container/list"
)

/*
实现LRU缓存淘汰算法
查找-删除-新增-修改
*/

// Cache 是一个LRU缓存，并发访问不安全
// Cache is LRU cache.It is not safe for concurrent access.
type Cache struct {
	maxBytes       int64                    // 允许使用的最大内存
	nowBytes       int64                    // 当前已使用的内存
	doubleLinkList *list.List               // Go语言标准库实现的双向链表
	cache          map[string]*list.Element // 字典：键是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可为nil
}

// 键值对entry是双向链表节点的数据类型，在链表中仍保留键好处在于淘汰队首节点时可以快速从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

// Value 用于返回值所占的内存大小
// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New 新建一个Cache
// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:       maxBytes,
		doubleLinkList: list.New(),
		cache:          make(map[string]*list.Element),
		OnEvicted:      onEvicted,
	}
}

// Get 从字典中查找对应的双向链表的节点，并将节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//将链表中的节点ele移动到你队尾（双向链表作为队列，队首队尾是相对的，在这里约定front为队尾）
		c.doubleLinkList.MoveToFront(ele)
		kv := ele.Value.(*entry) //断言
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰——删除，移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.doubleLinkList.Back() //取队首节点
	if ele != nil {
		c.doubleLinkList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                  //从字典中删除该节点
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新当前所用缓存
		if c.OnEvicted != nil {                                  //调用回调函数
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增或修改节点
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { //如果键存在，更新节点对应值并将该节点移到队尾
		c.doubleLinkList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { //不存在则新增，队尾添加新节点，并在字典中添加key和节点的映射
		ele := c.doubleLinkList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes { //更新当前所用缓存，如果超过设定的最大缓存值，则移除最近最少访问的节点
		c.RemoveOldest()
	}
}

// Len 缓存cache里的数据数
func (c *Cache) Len() int {
	return c.doubleLinkList.Len()
}
