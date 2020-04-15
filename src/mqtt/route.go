package mqtt

import (
	"tea/src/manage"
	"tea/src/mqtt/handle"
	"tea/src/mqtt/protocol"
)

var Route []handle.Handle

func Boot() {
	Route = make([]handle.Handle, 17)
	Route[1] = handle.NewConnect()
	Route[3] = handle.NewPublish()

}

func HandleCmd(pack protocol.Pack, client *manage.Client) {

	var getHandle handle.Handle
	switch pack.Cmd {
	case 1:
		getHandle = handle.NewConnect()
	case 3:
		getHandle = handle.NewPublish()

	case 8:
		getHandle = handle.NewSubscribe()

	}

	if getHandle != nil {
		getHandle.Handle(pack, client)
	}
}
