package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"sync"
	"time"
)

func main() {
	fmt.Println()
	c, r := 100, 100
	addr := "tcp://127.0.0.1:1883"
	// addr := "tcp://121.89.176.16:1883"
	qos := 0
	clients := make([]mqtt.Client, c)

	for i := 0; i < c; i++ {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(addr)
		opts.SetUsername("client_" + strconv.Itoa(i))

		opts.SetDefaultPublishHandler(handleMsg)

		clients[i] = mqtt.NewClient(opts)
		if token := clients[i].Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			i--
			continue
		}
	}
	defer disconnect(clients)
	fmt.Println("connect success")

	var wg sync.WaitGroup
	wg.Add(10000)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " start publish")
	for j, client := range clients {
		to := c - j - 1
		topic := "say/" + strconv.Itoa(to)
		payload := "hello " + strconv.Itoa(to)
		fmt.Printf("to: %d\ttopic: %s\tpayload: %s\n", to, topic, payload)
		for i := 0; i < r; i++ {
			//wg.Add(1)
			go func() {
				//topic = "say/0"
				//payload = "hello 0"
				token := client.Publish(topic, byte(qos), false, payload)
				if token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " finish publish")
	// time.Sleep(time.Second * 10)
}

func handleMsg(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(msg.Topic(), msg.Payload())
}

func disconnect(clients []mqtt.Client) {
	for _, clients := range clients {
		clients.Disconnect(250)
	}
	fmt.Println("connection close")
}
