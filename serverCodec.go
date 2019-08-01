package main

// 官方使用 gob 编解码网络数据
// 这里只是个演示，使用的是字符串编解码网络数据
// 你可以改成诸如 protobuf 、 json 、 自定义编解码等等

import (
	"encoding/binary"
	"io"
	"net/rpc"
)

type rpcServerCodec struct {
	rwc io.ReadWriteCloser
}

func (c *rpcServerCodec) ReadRequestHeader(r *rpc.Request) error {
	// r.ServiceMethod
	var n uint16
	binary.Read(c.rwc, binary.LittleEndian, &n)
	temp := make([]byte, n)
	c.rwc.Read(temp[:])
	r.ServiceMethod = string(temp)
	return nil
}

func (c *rpcServerCodec) ReadRequestBody(body interface{}) error {
	// body string
	var n uint16
	binary.Read(c.rwc, binary.LittleEndian, &n)
	temp := make([]byte, n)
	c.rwc.Read(temp[:])
	*body.(*string) = string(temp)
	return nil
}

func (c *rpcServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	// write r.ServiceMethod
	binary.Write(c.rwc, binary.LittleEndian, uint16(len(r.ServiceMethod)))
	c.rwc.Write([]byte(r.ServiceMethod))
	// write r.Error
	binary.Write(c.rwc, binary.LittleEndian, uint16(len(r.Error)))
	c.rwc.Write([]byte(r.Error))
	// write body
	data := []byte(*body.(*string))
	binary.Write(c.rwc, binary.LittleEndian, uint16(len(data)))
	c.rwc.Write(data)
	return nil
}

func (c *rpcServerCodec) Close() error {
	return c.rwc.Close()
}
