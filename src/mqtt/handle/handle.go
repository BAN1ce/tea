package handle

import "tea/src/client"

type Pack struct {
	Data            []byte
	FixedHeader     []byte
	FixHeaderLength int
	Cmd             int
}

func NewPack() *Pack {

	pack := new(Pack)

	return pack
}

type Handle interface {
	Handle(pack Pack, client *client.Client)
}
