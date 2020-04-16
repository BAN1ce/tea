package handle

import (
	"tea/src/manage"
	"tea/src/mqtt/protocol"
)

type Handle interface {
	Handle(pack protocol.Pack, client *manage.Client)
}
