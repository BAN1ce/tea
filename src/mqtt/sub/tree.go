package sub

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

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

func (n *Node) TraceClients(course func(key, value interface{}) bool) {

	n.clients.Range(course)

}

func AddTreeSub(topics []string, clientId uuid.UUID) {

	var first *Node
	var last *Node
	for i, t := range topics {
		n := newNode(t)
		if i == 0 {
			first = n
		} else {
			last.childes.Store(n, true)
		}
		last = n
	}
	treeHigh := len(topics) - 1

	if treeHigh < 1 {
		return
	}
	last.clients.Store(clientId, true)
	if last != nil && first != nil {
		if root, ok := TreeMap.LoadOrStore(topics[0], first); ok {
			// 第一级topic已存在时
			if root, ok := root.(*Node); ok {
				bfsWithAdd(root, clientId, topics)
			}
		}
	}

}
func GetTreeSub(topics []string) (*Node, bool) {

	if len(topics) > 0 {
		topicFirst := topics[0]
		if root, ok := TreeMap.Load(topicFirst); ok {
			if root, ok := root.(*Node); ok {
				return bfs(root, topics)
			}
		}

	}

	return nil, false

}

func DeleteClient(topic string, clientId uuid.UUID) {

}

func bfsWithAdd(root *Node, clientId uuid.UUID, topicSlice []string) {

	queue := make([]*Node, 0)
	queue = append(queue, root)
	i := 1
	first := queue[0]
	for len(queue) != 0 {
		first = queue[0]
		queue = queue[1:]
		first.childes.Range(func(key, value interface{}) bool {
			if value, ok := key.(*Node); ok {
				if i < len(topicSlice) {
					if value.trace == topicSlice[i] {
						queue = append(queue, value)
						i++
						return false
					}
				}
			}
			return true
		})

	}
	if i == len(topicSlice) {
		fmt.Println("i", i)
		first.clients.Store(clientId, true)
	} else {
		//todo 添加一个子树在first节点下
		childTree, err := topicSliceBeTree(topicSlice[i:])
		if err != nil {
			first.childes.Store(childTree, true)
		}

	}

}

func bfs(root *Node, topics []string) (*Node, bool) {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	i := 1
	parent := root
	for len(queue) != 0 {
		first := queue[0]
		queue = queue[1:]
		first.childes.Range(func(key, value interface{}) bool {
			if value, ok := key.(*Node); ok {
				if value.trace == topics[i] {
					queue = append(queue, value)
					parent = root
					i++
					return false
				}
			}
			return true
		})
		if i == len(topics) {
			return queue[len(queue)-1], true
		}
	}

	return parent, false

}
func topicSliceBeTree(topicsSlice []string) (*Node, error) {

	if len(topicsSlice) == 0 {
		return nil, errors.New("topicSlice length can not be 0")
	}

	var first *Node
	var last *Node
	for i, t := range topicsSlice {
		n := newNode(t)
		if i == 0 {
			first = n
		} else {
			last.childes.Store(n, true)
		}
		last = n
	}

	return first, nil

}
