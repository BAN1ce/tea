package handle

import (
	"fmt"
	"tea/src/client"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/sub"
	"tea/src/utils"
)

type PublishPack struct {
	protocol.Pack
	dup        int
	qos        int
	retain     int
	topicName  string
	identifier int
	payload    []byte
}
type Publish struct {
}

func newPublishPack(pack protocol.Pack) *PublishPack {
	fmt.Println("hello")
	p := new(PublishPack)
	p.Pack = pack
	p.retain = int(pack.FixedHeader[0] & 0x1)
	p.qos = int(pack.FixedHeader[0] & 0x6)
	p.dup = int(pack.FixedHeader[0] & 0x8)

	plc := pack.FixHeaderLength

	topicNameLength := utils.UtfLength(pack.Data[plc : plc+2])
	plc += 2
	p.topicName = string(pack.Data[plc : plc+topicNameLength])
	fmt.Println(p.topicName, "topicName")
	plc += topicNameLength
	if p.qos > 0 {
		p.identifier = utils.UtfLength(pack.Data[plc : plc+2])
		plc += 2
	}
	payloadLength := pack.BodyLength - (plc - pack.FixHeaderLength)

	p.payload = pack.Data[plc : plc+payloadLength]

	fmt.Println(string(p.payload))

	return p

}
func NewPublish() *Publish {

	return new(Publish)
}
func (p *Publish) Handle(pack protocol.Pack, client *client.Client) {

	publishPack := newPublishPack(pack)

	clientList,ok := sub.GetTopicSubClients(publishPack.topicName)

	if ok {
		clients := clientList.GetNode()
		fmt.Println(clients)
	}


	//todo 根据topic 找寻订阅topic的客户端将消息发布给客户端

}
