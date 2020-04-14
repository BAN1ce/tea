package handle

import (
	"tea/src/client"
	"tea/src/mqtt/protocol"
)

type Handle interface {
	Handle(pack protocol.Pack, client *client.Client)
}
