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

	node, ok := GetTreeSub(topics)
	if ok {
		flag := true
		for k, _ := range node.clients.m {
			if k == clientId {
				flag = false
			}
		}
		if flag {
			t.Error("Get Node But Can Find ClientID")
			t.FailNow()
		}
	} else {
		t.Error("Can Not Find Node")
		t.FailNow()
	}

	DeleteTreeSub(topics, clientId)

	node, ok = GetTreeSub(topics)
	if !ok {
		t.Log("Delete Sub Success")
	} else {
		t.Error("Delete Sub Fail")
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
		node, ok := GetTreeSub(strings.Split(topics[i], "/"))
		if ok {

			flag := true
			for k, _ := range node.clients.m {
				if k == clientIds[i] {
					flag = false
				}
			}
			if flag {
				t.Error("Get Node But Can Find ClientID")
			}

		} else {
			t.Error("Can Not Find Node")
		}

	}
}

/**
  不同客户端订阅不同的topic，每个topic订阅的客户端数量不同
*/
func TestTreeDiffClientSub(t *testing.T) {
	clientIds := make([]map[uuid.UUID]bool, 0)
	topics := make([]string, 0)
	for i := 0; i < 10; i++ {
		topic := fmt.Sprintf("product%d/device%d/get", i, i)
		cs := make(map[uuid.UUID]bool)
		for j := 0; j < i+1; j++ {
			clientId := uuid.New()
			cs[clientId] = true
		}
		clientIds = append(clientIds, cs)
		topics = append(topics, topic)
		topicSplit := strings.Split(topic, "/")
		for k, _ := range cs {
			AddTreeSub(topicSplit, k)
		}
	}

	for i := 0; i < 10; i++ {
		topic := fmt.Sprintf("product%d/device%d/get", i, i)
		node, ok := GetTreeSub(strings.Split(topic, "/"))
		if ok {

			flag := true
			for k, _ := range node.clients.m {
				if _, ok := clientIds[i][k]; ok {
				} else {
					flag = false
				}
			}
			if flag {
				t.Error("Get Node But Can Find ClientID")
			}

		} else {
			t.Error("topic can not find sub ", topic)
		}

	}

}
