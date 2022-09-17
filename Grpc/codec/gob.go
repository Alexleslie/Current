package codec

import (
	"Current/tools/logc"
	"bufio"
	"encoding/gob"
	"io"
)

type GobCodeC struct {
	conn    io.ReadWriteCloser
	buf     *bufio.Writer
	decoder *gob.Decoder
	encoder *gob.Encoder
}

// 检测GobCodeC是否实现了CodeC的所有接口
var _ CodeC = (*GobCodeC)(nil)

func (g GobCodeC) Close() error {
	return g.conn.Close()
}

func (g GobCodeC) ReadHeader(header *Header) error {
	return g.decoder.Decode(header)
}

func (g GobCodeC) ReadBody(body interface{}) error {
	return g.decoder.Decode(body)
}

func (g GobCodeC) Write(header *Header, body interface{}) error {
	defer func() {
		// Write结束时需要将缓冲区的内容全部写入
		err := g.buf.Flush()
		if err != nil {
			logc.Error("[GobCodeC.Write] Buf Flush error, err=[%+v]", err)
			_ = g.Close()
			return
		}
	}()
	if err := g.encoder.Encode(header); err != nil {
		logc.Error("[GobCodeC.Write] Encode header error, err=[%+v]", err)
		return err
	}
	if err := g.encoder.Encode(body); err != nil {
		logc.Error("[GobCodeC.Write] Encode body error, err=[%+v]", err)
		return err
	}
	return nil
}

func NewGobCodeC(conn io.ReadWriteCloser) CodeC {
	buf := bufio.NewWriter(conn)
	return &GobCodeC{
		conn:    conn,
		buf:     buf,
		decoder: gob.NewDecoder(conn),
		encoder: gob.NewEncoder(conn),
	}

}
