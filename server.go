package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"net/rpc"
)

func mockServerProcess(reqData []byte) (repData []byte) {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(reqData),
		wbuf: bytes.NewBuffer(nil),
	}
	buf := bufio.NewWriter(rwc)
	c := &gobServerCodec{rwc, gob.NewDecoder(rwc), gob.NewEncoder(buf), buf, false}
	if err := rpc.ServeRequest(c); err != nil {
		panic(err)
	}
	repData = rwc.wbuf.Bytes()
	return
}
