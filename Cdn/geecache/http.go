package geecache

import (
	"Current/Cdn/geecache/consistenthash"
	pb "Current/Cdn/geecache/geecachepb"
	"Current/Cdn/geecache/utils"
	"Current/Gee/gee"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

// HTTPClient 创建一个结构体HTTPClient作为承载节点间HTTP通信的核心数据结构（包括服务端和客户端）
type HTTPClient struct {
	//this peer's base URL,e.g. "http://example.net:8000"
	self        string                 // 记录自己地址，包括主机名/IP和端口
	basePath    string                 // 作为节点间通讯地址的前缀，默认是/_geecache/
	mu          sync.Mutex             // guards peers and httpGetters，客户端选择和获取缓存时的锁
	peers       *consistenthash.Map    // 一致性哈希算法主数据结构，根据具体的key选择节点
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"，映射远程节点与对应的httpGetter，每一个远程节点对应一个httpGetter(因为httpGetter与远程节点的地址baseURL有关)
}

// http客户端，对应远程节点
type httpGetter struct {
	baseURL string //将要访问的远程节点地址
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string, basePath string) *HTTPClient {
	return &HTTPClient{
		self:     self,
		basePath: basePath,
	}
}

// Log info with server name
func (p *HTTPClient) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// GetValueFromRemotePeer 获取远程节点的缓存值
func (h *httpGetter) GetValueFromRemotePeer(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v/%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()), //URL编码
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u) //使用http.Get()方式获取返回值，并转换成[]bytes类型
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server return:%v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body:%v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body:%v", err)
	}
	return nil
}

// GeneratePeersAndGetter 添加传入的节点，为每一个节点创建一个HTTP客户端httpGetter
func (p *HTTPClient) GeneratePeersAndGetter(replicas int, peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(replicas, nil)
	p.peers.AddPhysicalAndVirtualPeer(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 包装了一致性哈希算法的Get方法，根据具体的key选择节点，返回节点对应的HTTP客户端
func (p *HTTPClient) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.GetRealPeerFromKey(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// NewApiServerCachePeers 生成一致性hash算法的主节点
// 主节点用来分发所有请求
func NewApiServerCachePeers(listenAddr string, cacheAttrs []string, basePath string) *HTTPClient {
	cachePeers := NewHTTPPool(listenAddr, basePath)
	cachePeers.GeneratePeersAndGetter(utils.DefaultReplicas, cacheAttrs...)
	return cachePeers
}

// StartCacheServer 启动缓存服务器的监听
// peerAttrs 代表缓存服务器的地址
// basePath 代表缓存服务器的监听路径
func StartCacheServer(peerAttrs []string, basePath string) {
	cacheServer := gee.Default()
	cacheServer.GET(basePath, GetValueAtCacheServer)
	for _, attr := range peerAttrs {
		u, _ := url.Parse(attr)
		go cacheServer.Run(u.Host)
	}
}

func parseCachePath(path string) []string {
	if path[0] == '/' {
		path = path[1:]
	}
	parts := strings.SplitN(path, "/", 3)
	if len(parts) != 3 {
		return nil
	}
	return parts
}

// GetValueAtCacheServer 处理主节点机器的http请求（找缓存值）
// 解析http请求得到缓存key
// 从group里找缓存key
func GetValueAtCacheServer(ctx *gee.Context) {
	parts := parseCachePath(ctx.Path)
	if parts == nil {
		ctx.ReturnFail(http.StatusInternalServerError, "bad request")
		return
	}
	groupName, key := parts[1], parts[2]
	group := GetGroup(groupName)
	if group == nil {
		ctx.ReturnFail(http.StatusNotFound, "no such group:"+groupName)
		return
	}

	//通过group.Get(key)获取访问数据
	view, err := group.GetKeyAtCacheServer(key)
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		ctx.ReturnFail(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.WriteBytes(http.StatusOK, body)
}
