package main

import (
	"flag"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1000, "number of tcp connections")
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {

	flag.Parse()

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	conns := make([]mqtt.Client, *connections)

	addr := fmt.Sprintf("tcp://%s:1883", *ip)

	for i := 0; i < len(conns); i++ {

		opts := mqtt.NewClientOptions().AddBroker(addr).SetClientID("emqx_test_client")

		opts.SetKeepAlive(30 * time.Second)
		// 设置消息回调处理函数
		opts.SetDefaultPublishHandler(f)
		opts.SetPingTimeout(1 * time.Second)

		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		// 订阅主题
		if token := c.Subscribe("product/#", 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			i--
		}
		conns[i] = c
	}
	fmt.Println(fmt.Sprintf("create %d clients success", *connections))

	for {
		for i, c := range conns {
			// 发布消息
			go func() {
				token := c.Publish(fmt.Sprintf("product/%d", i), 0, false, fmt.Sprintf("Hello World %d", i))
				token.Wait()
			}()
		}
		time.Sleep(1 * time.Second)
	}
}
