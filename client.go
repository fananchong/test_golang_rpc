package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"net/rpc"
)

func mockClientSend(method, name string) (reqData []byte) {
	rwc := &readWriteCloser{
		rbuf: bytes.NewBuffer(nil),
		wbuf: bytes.NewBuffer(nil),
	}
	encBuf := bufio.NewWriter(rwc)
	c := &gobClientCodec{rwc, gob.NewDecoder(rwc), gob.NewEncoder(encBuf), encBuf}

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
	encBuf := bufio.NewWriter(rwc)
	c := &gobClientCodec{rwc, gob.NewDecoder(rwc), gob.NewEncoder(encBuf), encBuf}

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
