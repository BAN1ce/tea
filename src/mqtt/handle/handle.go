package handle

import "tea/src/client"
import "tea/src/mqtt"

type Handle interface {
	Handle(pack mqtt.Pack, client *client.Client)
}
