package handle

import (
	"strings"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
	"tea/src/mqtt/sub"
	"tea/src/utils"
)

type unSubscribePack struct {
	protocol.Pack
	identifier uint16
	payload    []byte
	topics     []string
}
type UnSubscribe struct {
}

func NewUnSubscribe() *UnSubscribe {

	return new(UnSubscribe)
}
func newUnSubscribePack(pack protocol.Pack) *unSubscribePack {
	p := new(unSubscribePack)
	p.Pack = pack
	p.topics = make([]string, 0)

	plc := pack.FixHeaderLength
	p.identifier = utils.BytesToUint16(pack.Data[plc : plc+2])
	plc += 2
	p.payload = pack.Data[plc:]
	for p.PackLength > plc+2 {
		l := utils.UtfLength(p.Pack.Data[plc : plc+2])
		plc += 2
		if plc+l <= p.PackLength {
			p.topics = append(p.topics, string(p.Data[plc:plc+l]))
			plc += l
		}
	}

	return p
}

func (u *UnSubscribe) Handle(pack protocol.Pack, client *manage.Client) {

	p := newUnSubscribePack(pack)

	for _, topic := range p.topics {
		topicSlice := strings.Split(topic, "/")

		//fixme delete wildcards sub 删除模糊订阅
		if utils.HasWildcard(topicSlice) {
			sub.DeleteTreeSub(topicSlice, client.Uid)
		} else {
			sub.DeleteHashSub(topic, client.Uid)
		}

	}

	unSuback := response.NewUnSuback(p.identifier)

	protocol.Encode(unSuback, client)

}
