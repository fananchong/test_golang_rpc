package main

import "fmt"

// Greeter RPC 要调用的类（演示）
type Greeter struct {
}

// Hello RPC 要调用的类方法（演示）
func (g *Greeter) Hello(name string, result *string) error {
	fmt.Println("hello", name)
	*result = "OK"
	return nil
}
