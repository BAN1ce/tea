package handle

import (
	"sync/atomic"
	"tea/src/distributed"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
)

type Hb struct {
}

func NewHb() *Hb {

	return new(Hb)
}

func (h *Hb) Handle(pack *protocol.Pack, client *manage.Client) {

	atomic.AddUint32(&distributed.HBTotalCount, 1)

	protocol.Encode(response.NewHback(), client)
}
