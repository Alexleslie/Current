package main

import (
	"Current/Cdn/geecache"
	"Current/Cdn/geecache/utils"
	"Current/Gee/gee"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

var ApiServerGroup *geecache.Group

func getFromDb(key string) ([]byte, error) {
	log.Println("[SlowDB] search key", key)
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("%s not exist", key)
}

func getKeyInApi(ctx *gee.Context) {
	key := ctx.Query("key")
	view, err := ApiServerGroup.Get(key)
	if err != nil {
		ctx.ReturnFail(http.StatusInternalServerError, err.Error())
	}
	ctx.WriteBytes(http.StatusOK, view.ByteSlice())

}

func main() {
	var basePath string

	basePath = utils.DefaultBasePath

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	var addrsForListen []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
		addrsForListen = append(addrsForListen, v[7:])
	}

	geecache.StartCacheServer(addrsForListen, basePath)

	apiServer := gee.Default()
	ApiServerGroup = geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(getFromDb))
	ApiServerGroup.RegisterPeers(geecache.NewApiServerCachePeers(apiAddr, addrs, basePath))

	apiServer.GET("/api", getKeyInApi)
	apiServer.Run(apiAddr[7:])

}
