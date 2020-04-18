package sub

import (
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
		last.clients.Store(clientId, true)
		if root, ok := TreeMap.LoadOrStore(topics[0], first); ok {
			// 第一级topic已存在时
			if root, ok := root.(*Node); ok {
				bfsWithAdd(root, first, treeHigh, clientId)
			}
		}
	}

}

func DeleteClient(topic string, clientId uuid.UUID) {

}

func bfsWithAdd(root *Node, newTree *Node, newTreeHigh int, clientId uuid.UUID) {

	queue := make([]*Node, 0)
	queue = append(queue, root)
	i := 1
	parent := root
	for len(queue) != 0 {
		first := queue[0]
		queue = queue[1:]
		first.childes.Range(func(key, value interface{}) bool {
			if value, ok := key.(*Node); ok {
				if value.trace == newTree.trace {
					queue = append(queue, value)
					parent = root
					i++
					newTree.childes.Range(func(key, value interface{}) bool {
						if node, ok := key.(*Node); ok {
							newTree = node
							return false
						}
						return true
					})
					return false
				}
			}
			return true
		})
		if i == newTreeHigh {
			queue[len(queue)-1].clients.Store(clientId, true)
			return
		} else {
			parent.childes.Store(newTree, true)
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
