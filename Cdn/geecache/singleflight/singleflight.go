package singleflight

import (
	"sync"
)

/*
缓存雪崩：缓存在同一时刻全部失效，造成瞬间DB请求量大，压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的key设置了相同的过期时间引起。
缓存击穿：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到DB，造成瞬时DB请求量大、压力骤增。
缓存穿透：查询有一个不存在是=的数据，因为不存在则不会写到缓存中，所以每次都会区请求DB，如果瞬间流量过大，穿透到DB，导致宕机。
*/

// 代表正在进行中或者已经结束的请求
type call struct {
	wg  sync.WaitGroup // 锁，避免重入，并发协程之间不需要消息传递
	val interface{}
	err error
}

// Group sigleflight主数据结构，管理不同的key的请求（call）
type Group struct {
	mu sync.Mutex       // 保护m不被并发读写，正常字典都要避免被并发读写
	m  map[string]*call // 管理当前不同key的请求
}

// Do 针对相同的key，无论Do被调用和多少次，fn函数只调用一次，等待fn调用结束，返回返回值和错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	//如果之前没有任何请求，新建一个管理key请求的group，延时初始化
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	//如果请求中已有key请求正在进行中，则等待
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	//如果之前请求没有当前key请求，新建key请求
	c := new(call)
	//发起请求前加锁
	c.wg.Add(1)
	//添加到g.m，表示key请求正在处理
	g.m[key] = c
	g.mu.Unlock()

	//调用fn函数，发起请求获得返回值和错误
	c.val, c.err = fn()
	//请求结束
	c.wg.Done()

	//当前请求已完成，删除group中当前的请求
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
