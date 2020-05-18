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
	ip          = flag.String("ip", "127.0.0.1:1883", "server IP")
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

	satistic := make(map[string]int)

	fmt.Println(*ip)

	addr := fmt.Sprintf("tcp://%s", *ip)

	for i := 0; i < len(conns); i++ {

		opts := mqtt.NewClientOptions().AddBroker(addr).SetUsername("hello" + string(i))

		opts.SetKeepAlive(30 * time.Second)
		// 设置消息回调处理函数
		opts.SetDefaultPublishHandler(f)

		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			i--
			continue
		}

		// 订阅主题
		topic := fmt.Sprintf("product/%d", len(conns)-i-1)
		if token := c.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
			satistic[message.Topic()] += 1
		}); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			i--
		}
		conns[i] = c

	}
	fmt.Println(fmt.Sprintf("create %d clients success", *connections))

	j := 0
	for {
		for i, c := range conns {
			// 发布消息
			token := c.Publish(fmt.Sprintf("product/%d", i), 0, false, fmt.Sprintf("Hello World %d", i))
			token.Wait()
		}
		time.Sleep(1 * time.Second)
		j++

		if j == 10000 {
			break
		}
	}

	fmt.Println("All message published")

	time.Sleep(5 * time.Second)
	for k, v := range satistic {
		if v != 10000 {
			fmt.Println(k, v)
		}
	}
	select {}

}
