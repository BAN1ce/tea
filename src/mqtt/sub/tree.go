package sub

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

/**
订阅树的根节点
*/
var TreeRoot *Node

func init() {

	TreeRoot = newNode("/")
}

/**
节点的客户端ID map
*/
type nodeClients struct {
	M  map[uuid.UUID]interface{}
	Mu sync.RWMutex
}

/**
订阅树节点的子节点map
*/
type childNodes struct {
	m  map[string]*Node
	mu sync.RWMutex
}

/**
新建一个节点客户端map
*/
func newNodeClients() *nodeClients {
	n := nodeClients{
		M: make(map[uuid.UUID]interface{}),
	}

	return &n
}

/**
新建一个节点子节点map
*/
func newChildNodes() *childNodes {

	m := childNodes{
		m: make(map[string]*Node),
	}
	return &m
}

/**
订阅树中的节点
*/
type Node struct {
	topicSection string
	ChildNodes   *childNodes
	Clients      *nodeClients
}

/**
新建一个节点
*/
func newNode(topicSection string) *Node {
	n := Node{
		topicSection: topicSection,
		ChildNodes:   newChildNodes(),
		Clients:      newNodeClients(),
	}

	return &n
}

/**
topic 切片转话成树，树的叶子节点存入clientID
*/
func topicSliceBeTree(topicsSlice []string, clientId uuid.UUID) (*Node, error) {

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
			last.ChildNodes.m[t] = n
		}
		last = n
	}
	last.Clients.M[clientId] = true

	return first, nil

}

/**
订阅树添加订阅客户端订阅
*/
func AddTreeSub(topicSlice []string, clientId uuid.UUID) {

	queue := make([]*Node, 0)
	queue = append(queue, TreeRoot)
	i := 0
	first := queue[0]
	for len(queue) != 0 && i < len(topicSlice) {
		first = queue[0]
		queue = queue[1:]

		first.ChildNodes.mu.Lock()
		if childNode, ok := first.ChildNodes.m[topicSlice[i]]; ok {
			queue = append(queue, childNode)
			i++
			if i == len(topicSlice) {
				childNode.Clients.Mu.Lock()
				childNode.Clients.M[clientId] = true
				childNode.Clients.Mu.Unlock()
			}
		} else {
			if childTree, err := topicSliceBeTree(topicSlice[i:], clientId); err == nil {
				first.ChildNodes.m[topicSlice[i]] = childTree
			}
		}
		first.ChildNodes.mu.Unlock()
	}

}

/**
获取topic在订阅树中保存的节点，未找到返回false
*/
func GetTreeSub(topicSlice []string) (*Node, bool) {
	queue := make([]*Node, 0)
	queue = append(queue, TreeRoot)
	i := 0
	first := queue[0]
	for len(queue) != 0 {
		first = queue[0]
		queue = queue[1:]

		first.ChildNodes.mu.RLock()
		if childNode, ok := first.ChildNodes.m[topicSlice[i]]; ok {
			queue = append(queue, childNode)
			i++
		}
		first.ChildNodes.mu.RUnlock()
		if i == len(topicSlice) {
			return queue[len(queue)-1], true
		}
	}
	return nil, false
}

/**
订阅树删除订阅节点和客户端ID,如果删除后节点的clientID个数为0则删除从父节点中删去当前节点
*/
func DeleteTreeSub(topicSlice []string, clientId uuid.UUID) {

	queue := make([]*Node, 0)
	queue = append(queue, TreeRoot)
	i := 0
	first := queue[0]

	for len(queue) != 0 && i < len(topicSlice) {
		first = queue[0]
		queue = queue[1:]

		first.ChildNodes.mu.Lock()
		if childNode, ok := first.ChildNodes.m[topicSlice[i]]; ok {
			queue = append(queue, childNode)
			i++
			if i == len(topicSlice) {
				childNode.Clients.Mu.Lock()
				delete(childNode.Clients.M, clientId)
				if len(childNode.Clients.M) == 0 {
					delete(first.ChildNodes.m, topicSlice[i-1])
				}
				childNode.Clients.Mu.Unlock()

			}
		}
		first.ChildNodes.mu.Unlock()
	}

}

/**
订阅树的广度搜索
*/
func Bfs(root *Node, topicSlice []string) (*Node, bool) {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	i := 0
	parent := root
	first := queue[0]
	for len(queue) != 0 {
		first = queue[0]
		queue = queue[1:]

		first.ChildNodes.mu.RLock()
		if childNode, ok := first.ChildNodes.m[topicSlice[i]]; ok {
			queue = append(queue, childNode)
			i++
		}
		first.ChildNodes.mu.RUnlock()
		if i == len(topicSlice)-1 {
			return queue[len(queue)-1], true
		}
	}

	return parent, false
}
