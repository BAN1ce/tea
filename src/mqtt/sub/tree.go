package sub

import (
	"github.com/google/uuid"
	"strings"
	"sync"
)

/**
数的子节点双向链表
*/
type childList struct {
	ListLength int
	first      *subListNode
	last       *subListNode
	mutex      *sync.RWMutex
}

type Node struct {
	trace   string
	clients sync.Map //节点订阅的客户端clientId
	childes sync.Map //节点的子节点map
}

var TreeMap sync.Map

func newNode(trace string) *Node {
	return &Node{
		trace: trace,
	}
}
func AddClient(topic string, clientId uuid.UUID) {

	if len(topic) < 1 {
		return
	}
	topic = topic[1:]
	topics := strings.Split(topic, "/")
	var first *Node
	var last *Node
	for i, t := range topics {
		n := newNode(t)
		if i == 0 {
			first = n
		}
		last = n
	}
	if last != nil && first != nil {
		last.clients.Store(clientId, true)
		TreeMap.LoadOrStore(topics[0], first)
		//todo 第一级topic已存在时

	}

}

func DeleteClient(topic string, clientId uuid.UUID) {

}
