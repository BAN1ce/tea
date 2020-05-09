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

func DeleteHashSub( topic string,clientId uuid.UUID) {

	m, ok := subTopicHashTable.Load(topic)
	if ok {
		m, ok := m.(sync.Map)
		if ok {
			m.Delete(clientId)
		}
	}
}

