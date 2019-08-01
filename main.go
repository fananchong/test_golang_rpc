package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	rpc.Register(&Greeter{})
	reqData := mockClientSend("Greeter.Hello", "world") // 模拟客户端发送数据
	repData := mockServerRecvAndSend(reqData)           // 模拟服务器接收数据、处理、并发送数据
	result := mockClientRecv(repData)                   // 模拟客户端收到数据
	fmt.Println("result:", result)
}
