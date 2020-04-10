package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tea/src/client"
	"tea/src/server"
	"time"
)

func main() {

	netAddr2, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8001")
	s2 := server.NewServer(netAddr2)

	s2.SetProtocol(func() bufio.SplitFunc {
		return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF {
				return len(data), data[:len(data)], errors.New("EOF")
			}
			index := bytes.IndexAny(data, "&")
			if index != -1 {
				return index + 1, data[:index], nil
			}
			return 0, nil, nil
		}
	})
	s2.SetOnMessage(func(msg []byte, client *client.Client) error {
		client.Write(msg)
		return nil
	})

	s2.SetHbInterval(20 * time.Second)
	s2.SetHbContent([]byte{'&'})
	s2.Run(context.Background())

	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)

	s := <-c
	//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
	fmt.Println("get signal:", s)

}
