package codec

import (
	"Current/Grpc/utils"
	"io"
)

// Header RPC请求的请求头
// 其中包括着RPC请求中的基本信息
// RPC请求的body不固定，但是请求头固定形式
type Header struct {
	Service string // 请求的服务名
	Method  string // 被请求服务的方法/被调用的方法
	Seq     uint64 // 每次请求的唯一序号
	Error   string
}

// CodeC 实现编码与解码的接口规定
type CodeC interface {
	io.Closer
	ReadHeader(*Header) error         // 对header进行解码,同时写入header
	ReadBody(interface{}) error       // 对body进行解码，同时写入
	Write(*Header, interface{}) error // 将header与body编码后写入
}

// EncodeType 定义解码与编码的类型
type EncodeType string

type FuncGetter func(io.ReadWriteCloser) CodeC

var TypeToCodeCMap map[EncodeType]FuncGetter

func init() {
	TypeToCodeCMap = make(map[EncodeType]FuncGetter)
	TypeToCodeCMap[utils.GobType] = NewGobCodeC
}
