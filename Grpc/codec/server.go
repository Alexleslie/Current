package codec

import (
	"Current/tools/logc"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

type Option struct {
	MagicNumber int
	EncodeType
}

var DefaultGobOption = &Option{
	MagicNumber,
	GobType,
}

var DefaultServer = NewServer()

type Server struct {
}

type Request struct {
	header *Header
	argv   reflect.Value
	reply  reflect.Value
}

var InvalidRequest = struct{}{}

func NewServer() *Server {
	return &Server{}
}

func Accept(listen net.Listener) {
	_ = DefaultServer.Accept(listen)
}

func (server *Server) Accept(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			logc.Error("[Server.Accept] net listen accept error, err=[%+v]", err)
			return err
		}
		logc.Info("[Server.Accept] conn=[%+v]", conn)
		go server.ServeConnIncludeOption(conn)
	}
}

func (server *Server) ServeConnIncludeOption(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()

	var opt *Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		logc.Error("[Server.ServeConnIncludeOption] Json decode conn error, err=[%+v]", err)
	}
	logc.Info("[Server.ServeConnIncludeOption] Received option=[%+v]", opt)
	if opt.MagicNumber != MagicNumber {
		logc.Error("[Server.ServeConnIncludeOption] Received request is not a rpc request, "+
			"request magicNumber=[%+v]", opt.MagicNumber)
		return
	}
	getCodecFunc := TypeToCodeCMap[opt.EncodeType]
	if getCodecFunc == nil {
		logc.Error("[Server.ServeConnIncludeOption] Unsupported given type=[%+v]", opt.EncodeType)
		return
	}
	server.ServeCodeC(getCodecFunc(conn))
}

func (server *Server) ServeCodeC(c CodeC) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(c)
		if err != nil {
			if req == nil {
				logc.Error("[Server.ServeCodeC] Server readRequest error, req is nil, err=[%+v]", err)
				break
			}
			req.header.Error = err.Error()
			if err := server.sendResponse(c, req.header, InvalidRequest, sending); err != nil {
				logc.Error("[Server.ServeCodeC] sendResponse error, err=[%+v]", err)
			}
		}
		wg.Add(1)
		go server.handlerRequest(c, req, sending, wg)
	}
	wg.Wait()
	_ = c.Close()
}

func (server *Server) readRequestHeader(c CodeC) (*Header, error) {
	header := &Header{}
	if err := c.ReadHeader(header); err != nil {
		logc.Error("[Server.readRequestHeader] Read header error, err=[%+v]", err)
		return nil, err
	}
	return header, nil

}

func (server *Server) readRequest(c CodeC) (*Request, error) {
	header, err := server.readRequestHeader(c)
	if err != nil {
		return nil, err
	}
	logc.Info("[Server.readRequest] readRequestHeader success, header=[%+v]", header)
	req := &Request{header: header}

	req.argv = reflect.New(reflect.TypeOf(""))
	if err = c.ReadBody(req.argv.Interface()); err != nil {
		log.Println(err)
		logc.Error("[Server.readRequest] CodeC decode body is error,err=[%+v]", err)
		return req, err
	}
	return req, nil
}

func (server *Server) handlerRequest(c CodeC, req *Request, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	logc.Info("[Server.handlerRequest] request header=[%+v], request body=[%+v]", req.header, req.argv.Elem())
	req.reply = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.header.Seq))
	if err := server.sendResponse(c, req.header, req.reply.Interface(), mutex); err != nil {
		logc.Error("[Server.handlerRequest] sendResponse error, err=[%+v]", err)
	}
}

func (server *Server) sendResponse(c CodeC, header *Header, body interface{}, mutex *sync.Mutex) error {
	mutex.Lock()
	defer mutex.Unlock()
	if err := c.Write(header, body); err != nil {
		logc.Error("[Server.sendResponse] write body error, err=[%+v]", err)
		return err
	}
	return nil
}
