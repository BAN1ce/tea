package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tea/src/manage"
	"tea/src/mqtt"
	"tea/src/mqtt/protocol"
	"tea/src/server"
)

var (
	mqPort = flag.Int("mq_port", 1883, "mqtt port")
)

func main() {

	flag.Parse()

	mqtt.Boot()

	address := fmt.Sprintf("0.0.0.0:%d", *mqPort)

	netAddr2, _ := net.ResolveTCPAddr("tcp", address)
	s2 := server.NewServer(netAddr2)

	s2.SetProtocol(func() bufio.SplitFunc {

		return protocol.Input()

	})
	s2.SetOnMessage(func(msg []byte, client *manage.Client) error {
		pack := protocol.Decode(msg)

		mqtt.HandleCmd(*pack, client)

		return nil
	})

	s2.Run(context.Background())

	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)

	s := <-c
	//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
	fmt.Println("get signal:", s)

}
