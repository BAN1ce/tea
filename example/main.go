package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tea/src/server"
)

var (
	mqPort = flag.Int("mq_port", 1883, "mqtt port")
)

func main() {


	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *mqPort)

	netAddr2, _ := net.ResolveTCPAddr("tcp", address)
	server.NewServer(netAddr2)

	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)

	s := <-c
	//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
	fmt.Println("get signal:", s)

}
