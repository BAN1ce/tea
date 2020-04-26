package main

import (
	"bufio"
	"context"
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

func main() {

	mqtt.Boot()
	netAddr2, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:1883")
	s2 := server.NewServer(netAddr2)

	s2.SetProtocol(func() bufio.SplitFunc {

		return protocol.Input()

	})
	s2.SetOnMessage(func(msg []byte, client *manage.Client) error {

		//fmt.Println("receiver: ")
		//for _, v := range msg {
		//	fmt.Printf("%08s ", utils.ConvertToBin(int(v)))
		//}

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
