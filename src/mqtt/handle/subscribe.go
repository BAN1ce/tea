package handle

import (
	"strings"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
	"tea/src/mqtt/sub"
	"tea/src/utils"
)

type subscribePack struct {
	protocol.Pack
	identifier uint16
	payload    []byte
	topicQos   map[string]uint8
}

type Subscribe struct {
}

func newSubscribePack(pack protocol.Pack) *subscribePack {
	p := new(subscribePack)
	p.Pack = pack
	p.topicQos = make(map[string]uint8)

	plc := pack.FixHeaderLength
	p.identifier = utils.BytesToUint16(pack.Data[plc : plc+2])
	plc += 2
	p.payload = pack.Data[plc:]
	for p.PackLength > plc+2 {
		l := utils.UtfLength(p.Pack.Data[plc : plc+2])
		plc += 2
		if plc+l < p.PackLength {
			topic := string(p.Data[plc : plc+l])
			plc += l
			qos := p.Data[plc]
			plc += 1
			p.topicQos[topic] = qos
		}
	}

	return p
}

func NewSubscribe() *Subscribe {

	return new(Subscribe)
}

func (s *Subscribe) Handle(pack protocol.Pack, client *manage.Client) {
	p := newSubscribePack(pack)

	qosSlice := make([]byte, 0)
	for topic, qos := range p.topicQos {
		topicSlice := strings.Split(topic, "/")
		qosSlice = append(qosSlice, qos)
		client.Topics[topic] = true
		// 客户端模糊订阅和绝对订阅分开记录
		if utils.HasWildcard(topicSlice) {
			sub.AddTreeSub(topicSlice, client.Uid)
		} else {
			sub.AddHashSub(topic, client.Uid)
		}

	}
	suback := response.NewSuback(p.identifier, qosSlice)

	protocol.Encode(suback, client)
}
