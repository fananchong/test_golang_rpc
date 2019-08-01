package main

// golang rpc 已经抽象除了 io 层，因此可以使用 mock 的方式接管网络IO
// 使用 readWriteCloser 类接管 io 层，实际上只是内存中数据过了一遍
// 这种模式可以称之为 mock

import (
	"bytes"
	"errors"
	"net/rpc"
)

func mockClientSend(method, name string) (reqData []byte) {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(nil),
		wbuf: bytes.NewBuffer(nil),
	}
	c := &rpcClientCodec{rwc}

	req := &rpc.Request{}
	req.ServiceMethod = method
	if err := c.WriteRequest(req, name); err != nil {
		panic(err)
	}
	reqData = rwc.wbuf.Bytes()
	return
}

func mockClientRecv(repData []byte) (result string) {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(repData),
		wbuf: bytes.NewBuffer(nil),
	}
	c := &rpcClientCodec{rwc}

	response := rpc.Response{}
	if err := c.ReadResponseHeader(&response); err != nil {
		panic(err)
	}

	if response.Error != "" {
		c.ReadResponseBody(nil)
	} else {
		if err := c.ReadResponseBody(&result); err != nil {
			panic(errors.New("reading body " + err.Error()))
		}
	}
	return
}
