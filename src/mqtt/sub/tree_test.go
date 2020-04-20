package sub

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"testing"
)

/**
空订阅树添加一个订阅和搜索
*/
func TestTreeAddASub(t *testing.T) {

	topic := "product1/device1/get"
	clientId := uuid.New()
	topics := strings.Split(topic, "/")

	AddTreeSub(topics, clientId)

	flag := false
	node, ok := GetTreeSub(topics)
	if ok {
		node.clients.Range(func(key, value interface{}) bool {
			fmt.Println(key, value)
			if key == clientId {
				flag = true
				t.Log("Pass")
				return false
			}
			return true
		})
	}

	if flag == false {
		t.Error("can no find node")
		t.FailNow()

	}
}

/**
多个客户端订阅不同的topic
*/
func TestTreeAddPluralSub(t *testing.T) {

	clientIds := make([]uuid.UUID, 0)
	topics := make([]string, 0)
	for i := 0; i < 10; i++ {
		topic := fmt.Sprintf("product%d/device%d/get", i, i)
		clientId := uuid.New()
		clientIds = append(clientIds, clientId)
		topics = append(topics, topic)
		topicSplit := strings.Split(topic, "/")
		AddTreeSub(topicSplit, clientId)
	}

	for i := 0; i < 10; i++ {
		flag := false
		node, ok := GetTreeSub(strings.Split(topics[i], "/"))
		if ok {
			node.clients.Range(func(key, value interface{}) bool {
				if key == clientIds[i] {
					flag = true
					return false
				}
				return true
			})
		}
		if flag == false {
			t.Error("can not find sub", clientIds[i], topics[i])
		}

	}

}

func TestTreeDiffClientSub(t *testing.T) {

	clientIds := make([][]uuid.UUID, 0)
	topics := make([]string, 0)
	for i := 0; i < 10; i++ {
		topic := fmt.Sprintf("product%d/device%d/get", i, i)
		cs := make([]uuid.UUID, 0)
		for j := 0; j < i+1; j++ {
			clientId := uuid.New()
			cs = append(cs, clientId)
		}
		clientIds = append(clientIds, cs)
		topics = append(topics, topic)
		topicSplit := strings.Split(topic, "/")
		for _, v := range cs {
			AddTreeSub(topicSplit, v)
		}
	}

	for i := 0; i < 10; i++ {
		topic := fmt.Sprintf("product%d/device%d/get", i, i)
		n, ok := GetTreeSub(strings.Split(topic, "/"))
		if ok {

			for _, v := range clientIds[i] {

				if _, ok := n.clients.Load(v); ok == false {
					t.Error("topic can not find a client sub ", topic, v)
				}
			}

		} else {
			t.Error("topic can not find sub ", topic)
		}

	}
}
