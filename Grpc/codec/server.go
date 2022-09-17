package codec

import (
	"Current/Grpc/utils"
	"Current/tools/logc"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
)

type Option struct {
	MagicNumber int
	EncodeType
}

var DefaultGobOption = &Option{
	utils.MagicNumber,
	utils.GobType,
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

func (server *Server) Accept(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			logc.Error("[Server.Accept] net listen accept error, err=[%+v]", err)
			return err
		}
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
	if opt.MagicNumber != utils.MagicNumber {
		logc.Error("[Server.ServeConnIncludeOption] Received request is not a rpc request, "+
			"request magicNumber=[%+v]", opt.MagicNumber)
		return
	}
	codecFunc := TypeToCodeCMap[opt.EncodeType]
	if codecFunc == nil {
		logc.Error("[Server.ServeConnIncludeOption] Unsupported given type=[%+v]", opt.EncodeType)
		return
	}
	go server.ServeCodeC(codecFunc(conn))
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
	var header *Header
	if err := c.ReadHeader(header); err != io.EOF && err != io.ErrUnexpectedEOF {
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
	req := &Request{header: header}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err := c.ReadBody(req.argv); err != nil {
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
