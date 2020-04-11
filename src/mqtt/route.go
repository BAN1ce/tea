package mqtt

import (
	"tea/src/client"
	"tea/src/mqtt/handle"
)

var Route []handle.Handle

func Boot() {
	Route = make([]handle.Handle, 17)
	Route[1] = handle.NewConnect()

}

func HandleCmd(pack handle.Pack, client *client.Client) {

	var getHandle handle.Handle
	switch pack.Cmd {
	case 1:
		getHandle = handle.NewConnect()

	}

	if getHandle != nil {
		getHandle.Handle(pack, client)
	}
}
