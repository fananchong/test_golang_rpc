package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net/rpc"
)

type readWriteCloser struct {
	wbuf *bytes.Buffer
	rbuf *bytes.Buffer
}

func (rwc *readWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.rbuf.Read(p)
}

func (rwc *readWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.wbuf.Write(p)
}

func (rwc *readWriteCloser) Close() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()
	return nil
}

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
		c.ReadResponseBody(nil) // 丢弃 body 数据
	} else {
		if err := c.ReadResponseBody(&result); err != nil {
			panic(errors.New("reading body " + err.Error()))
		}
	}
	return
}

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

func main() {
	rpc.Register(&Greeter{})
	reqData := mockClientSend("Greeter.Hello", "world")
	repData := mockServerProcess(reqData)
	result := mockClientRecv(repData)
	fmt.Println("result:", result)
}
