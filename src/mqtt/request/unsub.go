package request

import (
	"tea/src/mqtt/protocol"
	"tea/src/utils"
)

type unSubscribePack struct {
	*protocol.Pack
	Identifier uint16
	Payload    []byte
	Topics     []string
}

func NewUnSubscribePack(pack *protocol.Pack) *unSubscribePack {
	p := new(unSubscribePack)
	p.Pack = pack
	p.Topics = make([]string, 0)

	plc := pack.FixHeaderLength
	p.Identifier = utils.BytesToUint16(pack.Data[plc : plc+2])
	plc += 2
	p.Payload = pack.Data[plc:]
	for p.PackLength > plc+2 {
		l := utils.UtfLength(p.Pack.Data[plc : plc+2])
		plc += 2
		if plc+l <= p.PackLength {
			p.Topics = append(p.Topics, string(p.Data[plc:plc+l]))
			plc += l
		}
	}

	return p
}
