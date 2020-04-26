package handle

import (
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
)

type Hb struct {
}

func NewHb() *Hb {

	return new(Hb)
}

func (h *Hb) Handle(pack protocol.Pack, client *manage.Client) {

	protocol.Encode(response.NewHback(), client)
}
