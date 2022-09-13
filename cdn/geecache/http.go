package geecache

import (
	"cdn/geecache/consistenthash"
	pb "cdn/geecache/geecachepb"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

/*提供被其他节点访问的能力（基于http）*/

const (
	defaultBasePath  = "/_geecache/"
	defaultResplicas = 50 // 默认虚拟节点数
)

// HTTPPool 创建一个结构体HTTPPool作为承载节点间HTTP通信的核心数据结构（包括服务端和客户端）
type HTTPPool struct {
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
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断访问路径前缀是否是basePath，不是返回错误
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}
	//打印访问方法和路径
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> requied
	// 安装约定访问路径格式切割得到groupname和key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]
	//获得group实例
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	//通过group.Get(key)获取访问数据
	view, err := group.Get(key)

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//使用w.Write将缓存值作为httpResponse的body返回
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(body)
	if err != nil {
		return
	}
}

// Get 获取远程节点
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v%v",
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

// Set 更新http，实例化一致性哈希算法，添加传入的节点，为每一个节点创建一个HTTP客户端httpGetter
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultResplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 包装了一致性哈希算法的Get方法，根据具体的key选择节点，返回节点对应的HTTP客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}
