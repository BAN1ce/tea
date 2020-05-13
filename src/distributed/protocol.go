package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type update struct {
	sourceMemberName string
	Action           string // add, del
	Data             map[string]string
}

func newUpdate() *update {

	u := new(update)

	return u
}

func handle(update *update) {

	Cluster.Mutex.RLock()
	member := Cluster.M[update.sourceMemberName]
	Cluster.Mutex.RUnlock()

	switch update.Action {
	case "online":
		if uid, err := uuid.FromBytes([]byte(update.Data["uid"])); err == nil {
			clientOnline(uid, member)
		}
	case "offline":
		if uid, err := uuid.FromBytes([]byte(update.Data["uid"])); err == nil {
			clientOffline(uid, member)
		}

	}

}

/**
客户端上线，记录客户端的id和对应的节点关系
*/
func clientOnline(clientId uuid.UUID, member *Member) {

	ClientUidNode.Store(clientId, member.Node)
	member.clientUid.Store(clientId, true)
}

/**
客户端下线,删除客户端id和对应的节点关系
*/
func clientOffline(clientId uuid.UUID, member *Member) {
	ClientUidNode.Delete(clientId)
	member.clientUid.Delete(clientId)
}

func clientSub() {

}

func clientUnsub() {

}

func UpdateOnline(memberName string, clientId uuid.UUID) {

	u := update{
		sourceMemberName: memberName,
		Action:           "online",
		Data:             make(map[string]string),
	}
	u.Data["uid"] = clientId.String()
	if b, err := json.Marshal(u); err == nil {
		fmt.Println("send connect")
		broadcasts.QueueBroadcast(&broadcast{
			msg:    b,
			notify: nil,
		})
	} else {
		fmt.Println("json error", err)
	}

}
