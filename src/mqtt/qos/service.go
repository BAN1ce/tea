package qos

import (
	"fmt"
	"strings"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
	"tea/src/mqtt/sub"
)

var LocalManage *manage.Manage

func HandleQosZero(topicName string, payload []byte) {

	fmt.Println("handle other node pub message")

	clients := sub.GetSubClientIds(topicName)

	pubResponse := response.NewPublish(0, 0, 0, []byte(topicName), payload)

	for _, id := range clients {
		if c, ok := LocalManage.GetClient(id); ok {
			protocol.Encode(pubResponse, c)

		}
	}

	topicSlice := strings.Split(topicName, "/")

	nodes := sub.GetWildCards(topicSlice)

	for _, node := range nodes {

		node.Clients.Mu.RLock()
		for clientId, _ := range node.Clients.M {
			if c, ok := LocalManage.GetClient(clientId); ok {
				protocol.Encode(pubResponse, c)
			}
		}
		node.Clients.Mu.RUnlock()
	}

}
