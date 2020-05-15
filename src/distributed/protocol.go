package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type update struct {
	SourceMemberName string
	Action           string // add, del
	Data             map[string]string
}

func newUpdate() *update {

	u := new(update)

	return u
}

func handle(update *update) {

	Cluster.Mutex.RLock()
	member := Cluster.M[update.SourceMemberName]
	Cluster.Mutex.RUnlock()

	switch update.Action {
	case "online":
		fmt.Println(update.Data["uid"])
		fmt.Println("member", member, update.SourceMemberName)
		if uid, err := uuid.Parse(update.Data["uid"]); err == nil {
			clientOnline(uid, member)
		} else {
			fmt.Println(err)
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
func clientOnline(clientId uuid.UUID, member *Node) {

	ClientUidNode.Store(clientId, member.Node)
	member.clientUid.Store(clientId, true)
}

/**
客户端下线,删除客户端id和对应的节点关系
*/
func clientOffline(clientId uuid.UUID, member *Node) {
	ClientUidNode.Delete(clientId)
	member.clientUid.Delete(clientId)
}

func clientSub() {

}

func clientUnsub() {

}

func UpdateOnline(nodeName string, clientId uuid.UUID) {

	clientOnline(clientId, localNode)
	u := update{
		SourceMemberName: nodeName,
		Action:           "online",
		Data:             make(map[string]string),
	}

	u.Data["uid"] = clientId.String()
	if b, err := json.Marshal(u); err == nil {
		broadcasts.QueueBroadcast(&broadcast{
			msg:    b,
			notify: nil,
		})
	} else {
		fmt.Println("json error", err)
	}

}

func UpdateOffline(nodeName string,clientId uuid.UUID)  {
	clientOffline(clientId, localNode)

	u := update{
		SourceMemberName: nodeName,
		Action:           "offline",
		Data:             make(map[string]string),
	}

	u.Data["uid"] = clientId.String()
	if b, err := json.Marshal(u); err == nil {
		broadcasts.QueueBroadcast(&broadcast{
			msg:    b,
			notify: nil,
		})
	} else {
		fmt.Println("json error", err)
	}

}
