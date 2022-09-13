package geecache

import (
	pb "cdn/geecache/geecachepb"
	"cdn/geecache/singleflight"
	"fmt"
	"log"
	"sync"
)

/*负责与外部交互，控制缓存存储和获取的主流程
							是
接受 key --> 检查是否被被缓存 -----> 返回缓存值
				|  否					   是
				|-----> 是否应从远程节点获取 -----> 与远程节点交互 --> 返回缓存值
							|  否
							|----->调用“回调函数”，获取值并添加到缓存 --> 返回缓存值
*/

// Getter 从一个key中获取数据
// 作为缓存未命中时候的回调函数
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 是实现Getter接口的函数类型，简称为接口型函数
type GetterFunc func(key string) ([]byte, error)

// Get 是Getter接口的Get方法实现
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 是缓存的一个命名空间
// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string              // 每个Group拥有唯一的名称name
	getter    Getter              // 缓存未命中时获取源数据的回调(callback)
	mainCache cache               // 并发缓存
	peers     PeerPicker          // 客户端/远程节点
	loader    *singleflight.Group // 确保每一个key只被请求一次（避免并发多次请求）
}

var (
	mu     sync.RWMutex              // 只读锁
	groups = make(map[string]*Group) // 全局变量
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// RegisterPeers 将PeerPicker接口的HTTPPool注入Group中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if peers == nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// GetGroup 搜索特定名称的Group，用了只读锁，不涉及任何冲突变量的写操作
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[GeeCache] hit")
		return v, nil
	}

	// 缓存不存在，调用load函数
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 无论并发调用放的数量如何，每个key只请求一次（本地或远程）,封装singleflight.Group.Do()
	// 分布式场景下会调用getFromPeer从其他节点获取缓存
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			// 使用PickPeer()方法调用getFromPeer()从远程获取节点
			// 若是本机节点或失败，则调用getLocally()。
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		// 回调函数，单机获取源数据
		return g.getLocally(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return view.(ByteView), nil
}

// 使用实现了PeerGetter接口的httpGetter访问其他HTTP客户端（远程节点）获取缓存
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

// getLocally调用用户回调函数g.getter.Get()获取源数据，并且将源数据添加到缓存mainCache中（通过populateCache方法）
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
