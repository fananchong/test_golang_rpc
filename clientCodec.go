package main

// 官方使用 gob 编解码网络数据
// 这里只是个演示，使用的是字符串编解码网络数据
// 你可以改成诸如 protobuf 、 json 、 自定义编解码等等

import (
	"encoding/binary"
	"io"
	"net/rpc"
)

type rpcClientCodec struct {
	rwc io.ReadWriteCloser
}

func (c *rpcClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	// write r
	binary.Write(c.rwc, binary.LittleEndian, uint16(len(r.ServiceMethod)))
	c.rwc.Write([]byte(r.ServiceMethod))
	// write body
	data := []byte(body.(string))
	binary.Write(c.rwc, binary.LittleEndian, uint16(len(data)))
	c.rwc.Write(data)
	return
}

func (c *rpcClientCodec) ReadResponseHeader(r *rpc.Response) error {
	var n uint16
	// r.ServiceMethod
	binary.Read(c.rwc, binary.LittleEndian, &n)
	temp := make([]byte, n)
	c.rwc.Read(temp[:])
	r.ServiceMethod = string(temp)
	// r.Error
	binary.Read(c.rwc, binary.LittleEndian, &n)
	temp = make([]byte, n)
	c.rwc.Read(temp[:])
	r.Error = string(temp)
	return nil
}

func (c *rpcClientCodec) ReadResponseBody(body interface{}) error {
	var n uint16
	binary.Read(c.rwc, binary.LittleEndian, &n)
	temp := make([]byte, n)
	c.rwc.Read(temp[:])
	*body.(*string) = string(temp)
	return nil
}

func (c *rpcClientCodec) Close() error {
	return c.rwc.Close()
}
