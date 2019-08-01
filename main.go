package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	rpc.Register(&Greeter{})
	reqData := mockClientSend("Greeter.Hello", "world")
	repData := mockServerProcess(reqData)
	result := mockClientRecv(repData)
	fmt.Println("result:", result)
}
