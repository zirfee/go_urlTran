package main

import (
	"flag"
	"fmt"
	"os"
	"urlServer"
)

var (
	masterAddr = flag.String("masterAddr", "", "master端口")
	port       = flag.String("port", "127.0.0.1:8080", "slave端口")
	enableRpc  = flag.Bool("enableRpc", false, "是否启用RPC")
)

func main() {
	flag.Parse()
	go urlServer.Start(masterAddr, port, enableRpc)
	fmt.Println("服务器正在运行。。。输入exit退出")
	isExit := ""
	for isExit != "exit" {
		fmt.Scanln(&isExit)
	}
	os.Exit(0)

}
