package main

import (
	"flag"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"sync"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1:1883", "server IP")
	connections = flag.Int("conn", 1000, "number of tcp connections")
	per         = flag.Int("per", 1000, "number of messages count per connection")
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

	var mu sync.RWMutex

	statistic := make(map[string]int)

	fmt.Println(*ip)

	addr := fmt.Sprintf("tcp://%s", *ip)

	for i := 0; i < len(conns); i++ {

		opts := mqtt.NewClientOptions().AddBroker(addr).SetUsername("hello" + string(i))

		opts.SetKeepAlive(30 * time.Second)
		opts.SetPingTimeout(120 * time.Second)

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
			mu.Lock()
			statistic[message.Topic()] += 1
			mu.Unlock()
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
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			token := c.Publish(fmt.Sprintf("product/%d", i), 0, false, fmt.Sprintf("Hello World %d - %s", i, timeStr))
			token.Wait()
		}
		time.Sleep(1 * time.Second)
		j++

		if j == *per {
			break
		}
	}

	fmt.Println("All message published")

	time.Sleep(60 * time.Second)

	t := time.NewTicker(5 * time.Second)

	for {
		select {

		case <-t.C:
			for k, v := range statistic {
				if v != *per {
					fmt.Println(k, v, "miss receiver")
				}
			}

			fmt.Println("all scan")

		}
	}

}
