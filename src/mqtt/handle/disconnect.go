package handle

import (
	"tea/src/manage"
	"tea/src/mqtt/protocol"
)

type Disconnect struct {
}

func NewDisconnect() *Disconnect {

	return new(Disconnect)
}

func (c *Disconnect) Handle(pack protocol.Pack, client *manage.Client) {

	client.Stop()

}
