package request

import (
	"tea/src/mqtt/protocol"
	"tea/src/utils"
)

type subscribePack struct {
	*protocol.Pack
	Identifier uint16
	Payload    []byte
	TopicQos   map[string]uint8
}

func NewSubscribePack(pack *protocol.Pack) *subscribePack {
	p := new(subscribePack)
	p.Pack = pack
	p.TopicQos = make(map[string]uint8)
	plc := pack.FixHeaderLength
	p.Identifier = utils.BytesToUint16(pack.Data[plc : plc+2])
	plc += 2
	p.Payload = pack.Data[plc:]
	for p.PackLength > plc+2 {
		l := utils.UtfLength(p.Pack.Data[plc : plc+2])
		plc += 2
		if plc+l < p.PackLength {
			topic := string(p.Data[plc : plc+l])
			plc += l
			qos := p.Data[plc]
			plc += 1
			p.TopicQos[topic] = qos
		}
	}

	return p
}
