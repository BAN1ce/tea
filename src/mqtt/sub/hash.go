package sub

import (
	"github.com/google/uuid"
	"sync"
)

/**
topic的绝对订阅哈希表
*/
var subTopicHashTable sync.Map

func Sub(topic string, clientId uuid.UUID) {
	var subMap sync.Map
	subMap.Store(clientId, true) //fix 节点IP
	m, ok := subTopicHashTable.LoadOrStore(topic, subMap)
	if ok {
		m, ok := m.(sync.Map)
		if ok {
			m.Store(clientId, true)
		}
	}
}

func GetSubClientIds(topic string) []uuid.UUID {

	clients := make([]uuid.UUID, 0)
	m, ok := subTopicHashTable.Load(topic)
	if ok {
		m, ok := m.(sync.Map)
		if ok {
			m.Range(func(key, value interface{}) bool {
				clients = append(clients, key.(uuid.UUID))
				return true
			})
		}
	}

	return clients

}

func UnSub(clientId uuid.UUID, topic string) {

	m, ok := subTopicHashTable.Load(topic)
	if ok {
		m, ok := m.(sync.Map)
		if ok {
			m.Delete(clientId)
		}
	}
}

/**
订阅客户端链表
*/
type subList struct {
	ListLength int
	first      *subListNode
	last       *subListNode
	mutex      *sync.RWMutex
}

/**
链表节点
*/
type subListNode struct {
	clientId uuid.UUID
	next     *subListNode
	pre      *subListNode
}

/**
获取topic订阅的所有客户端ID
*/
func (l *subList) GetNode() []uuid.UUID {
	l.mutex.RLock()
	nodes := make([]uuid.UUID, l.ListLength)
	tmp := l.first
	i := 0
	for tmp != nil {

		nodes[i] = tmp.clientId
		i++
		tmp = tmp.next
	}
	l.mutex.RUnlock()
	return nodes
}

/**
新建一个订阅链表
*/
func newSubList() *subList {

	l := new(subList)

	l.mutex = new(sync.RWMutex)

	return l
}

/**
新建一个订阅链表节点
*/
func newSubListNode(clientId uuid.UUID) *subListNode {

	s := new(subListNode)

	s.clientId = clientId
	s.next = nil

	return s
}

/**
获取topic订阅的list
*/
func GetTopicSubClients(topic string) (*subList, bool) {
	sl, ok := subTopicHashTable.Load(topic)
	if ok {
		sl, ok := sl.(*subList)

		return sl, ok
	}

	return nil, false
}

/**
删除topic订阅中的一个客户端
*/
func DeleteHashSub(topic string, clientId uuid.UUID) {

	sl, ok := subTopicHashTable.Load(topic)
	if ok {
		sl, ok := sl.(*subList)

		if ok {
			sl.mutex.Lock()
			tmp := sl.first
			length := 0
			if sl.ListLength > 0 {
				for tmp != nil {
					if tmp.clientId == clientId {
						if tmp.pre != nil {
							tmp.pre.next = tmp.next
						}
						if tmp.next != nil {
							tmp.next.pre = tmp.pre
						}
						sl.ListLength--
					}
					tmp = tmp.next

				}
				length = sl.ListLength
			}
			sl.mutex.Unlock()
			if length == 0 {
				subTopicHashTable.Delete(topic)
			}
		}
	}
}

/**
topic订阅的链表里添加一个新的客户端
*/
func AddHashSub(topic string, clientId uuid.UUID) {

	l := newSubList()
	n := newSubListNode(clientId)
	l.first = n
	l.last = n
	l.ListLength = 1
	tmp := l.first
	for tmp != nil {
		tmp = tmp.next
	}

	sl, loaded := subTopicHashTable.LoadOrStore(topic, l)
	if loaded == true {
		value, ok := sl.(*subList)
		if ok {
			value.mutex.Lock()
			value.last.next = n
			n.pre = value.last
			value.last = n
			value.ListLength++
			value.mutex.Unlock()

		}

	}

}
