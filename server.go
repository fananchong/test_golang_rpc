package main

// golang rpc 已经抽象除了 io 层，因此可以使用 mock 的方式接管网络IO
// 使用 readWriteCloser 类接管 io 层，实际上只是内存中数据过了一遍
// 这种模式可以称之为 mock

import (
	"bytes"
	"net/rpc"
)

func mockServerRecvAndSend(reqData []byte) (repData []byte) {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(reqData),
		wbuf: bytes.NewBuffer(nil),
	}
	c := &rpcServerCodec{rwc}
	if err := rpc.ServeRequest(c); err != nil {
		panic(err)
	}
	repData = rwc.wbuf.Bytes()
	return
}
