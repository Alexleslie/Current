package main

import (
	"Current/Grpc/codec"
	"Current/tools/logc"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func startListenServer(attr chan string) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logc.Error("[startListenServer] listen TCP connection error, err=[%+v]", err)
		return
	}
	logc.Info("[startListenServer] Start listening, attr=[%+v]", listener.Addr().String())
	attr <- listener.Addr().String()
	codec.Accept(listener)
}

func startClient(attr chan string) {
	conn, err := net.Dial("tcp", <-attr)
	if err != nil {
		logc.Error("Dial TCP error, err=[%+v]", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	_ = json.NewEncoder(conn).Encode(codec.DefaultGobOption)
	c := codec.NewGobCodeC(conn)

	for i := 0; i < 2; i++ {
		h := &codec.Header{
			Service: "Math",
			Method:  "Sum",
			Seq:     uint64(i),
		}
		body := fmt.Sprintf("1+1=?, req=[%+v]", h.Seq)
		if err = c.Write(h, body); err != nil {
			logc.Error("[startClient] write header and body error, err=[%+v]", err)
		}
		time.Sleep(1 * time.Second)
		rspHeader := &codec.Header{}
		var rspBody interface{}
		_ = c.ReadHeader(rspHeader)
		_ = c.ReadBody(rspBody)
		logc.Info("header=[%+v], reply=[%+v]", rspHeader, rspBody)
	}
}

func main() {
	attr := make(chan string)
	go startListenServer(attr)
	time.Sleep(1 * time.Second)
	startClient(attr)

}
