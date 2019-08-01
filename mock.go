package main

// golang rpc 已经抽象除了 io 层，因此可以使用 mock 的方式接管网络IO
// 使用 readWriteCloser 类接管 io 层，实际上只是内存中数据过了一遍
// 这种模式可以称之为 mock

import "bytes"

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
