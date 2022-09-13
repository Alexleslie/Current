package main

import (
	"cdn/geecache"
	"fmt"
	"log"
	"net/http"
)

/*
type server int

//http.ListenAndServe 接收 2 个参数，第一个参数是服务启动的地址;
//第二个参数是 Handler，任何实现了 ServeHTTP 方法的对象都可以作为 HTTP 的 Handler。
func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}
*/

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(getFromDb))
}

func getFromDb(key string) ([]byte, error) {
	log.Println("[SlowDB] search key", key)
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("%s not exist", key)
}

// 启动缓存服务器，创建HTTPPool，添加节点信息，注册到gee中，启动HTTP服务（main中共三个端口，8001/8002/8003），用户不感知。
func startCacheServer(addr string, addrs []string, geeCache *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	geeCache.RegisterPeers(peers)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 启动API服务(main中为端口9999)，与用户进行交互，用户感知
func startAPIServer(apiAddr string, geeCache *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := geeCache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, err = w.Write(view.ByteSlice())
			if err != nil {
				return
			}
		}))
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

// main()函数需要命令行传入port和api2 个参数，用来在指定端口启动 HTTP 服务
func main() {
	var port int
	var api bool

	api = true
	port = 8001
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
