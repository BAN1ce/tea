package distributed

import (
	"encoding/json"
	"github.com/google/uuid"
	"sync/atomic"
	"tea/src/utils"
)

var BroadcastedCount uint32

type BroadcastPubMessage struct {
	TopicName string
	Payload   []byte
	Dup       int
	Qos       int
	Retain    int
	Uid       uuid.UUID
}

func NewBroadcastPubMessage() (*BroadcastPubMessage, error) {

	p := new(BroadcastPubMessage)
	if uid, ok := uuid.NewUUID(); ok == nil {
		p.Uid = uid
		return p, nil
	} else {
		return nil, ok
	}
}

func BroadcastPub(p *BroadcastPubMessage) {

	b, _ := json.Marshal(p)

	b = append(utils.Uint16ToBytes(0x01), b...)

	atomic.AddUint32(&BroadcastedCount, 1)
	broadcasts.QueueBroadcast(&broadcast{msg: b, notify: nil})

}
