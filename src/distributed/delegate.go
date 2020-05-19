package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
	"sync"
	"sync/atomic"
	"tea/src/mqtt/qos"
	"tea/src/utils"
)

var BroadcastHandleTotalCount uint32

var BroadcastTotalCount uint32

var keyMap sync.Map

type delegate struct {
}

type eventDelegate struct {
}

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d delegate) NodeMeta(limit int) []byte {

	return nil
}

func (d delegate) NotifyMsg(msg []byte) {

	//todo 收到其他节点pub消息，去重处理且发布到订阅的客户端

	fmt.Println("notify msg", string(msg))

	cmd := utils.BytesToUint16(msg[0:2])

	switch cmd {

	case 0x01:
		pub := &BroadcastPubMessage{}
		json.Unmarshal(msg[2:], pub)
		atomic.AddUint32(&BroadcastTotalCount, 1)
		if _, ok := keyMap.LoadOrStore(pub.Uid.String(), true); ok {
			fmt.Println("exists message receiver from other node")
		} else {

			atomic.AddUint32(&BroadcastHandleTotalCount, 1)
			go qos.HandleQosZero(pub.TopicName, pub.Payload)
		}

	}
	return
}

func (d delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d delegate) LocalState(join bool) []byte {

	return nil
}

func (d delegate) MergeRemoteState(buf []byte, join bool) {

	fmt.Println("MergeRemoteState", string(buf))
	return
}

func (e eventDelegate) NotifyJoin(n *memberlist.Node) {

}

func (e eventDelegate) NotifyLeave(n *memberlist.Node) {
}

func (e eventDelegate) NotifyUpdate(*memberlist.Node) {
	fmt.Println("Notify Update")
}
