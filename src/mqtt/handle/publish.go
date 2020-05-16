package handle

import (
	"fmt"
	"strings"
	"tea/src/distributed"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/request"
	"tea/src/mqtt/sub"
	"tea/src/utils"
)

type Publish struct {
}

func newPublishPack(pack protocol.Pack) *request.PublishPack {

	p := new(request.PublishPack)
	p.Pack = pack
	p.Retain = int(pack.FixedHeader[0] & 0x1)
	p.Qos = int(pack.FixedHeader[0] & 0x6)
	p.Dup = int(pack.FixedHeader[0] & 0x8)

	plc := pack.FixHeaderLength

	topicNameLength := utils.UtfLength(pack.Data[plc : plc+2])
	plc += 2
	p.TopicName = string(pack.Data[plc : plc+topicNameLength])
	plc += topicNameLength
	if p.Qos > 0 {
		p.Identifier = utils.BytesToUint16(pack.Data[plc : plc+2])
		plc += 2
	}
	payloadLength := pack.BodyLength - (plc - pack.FixHeaderLength)

	p.Payload = pack.Data[plc : plc+payloadLength]

	return p

}

func NewPublish() *Publish {

	return new(Publish)
}
func (p *Publish) Handle(pack protocol.Pack, client *manage.Client) {

	publishPack := newPublishPack(pack)

	// 根据topic 找寻订阅topic的客户端将消息发布给客户端

	clients := sub.GetSubClientIds(publishPack.TopicName)
	for _, clientId := range clients {
		if c, ok := client.Manage.GetClient(clientId); ok {
			if publishPack.Qos != 0 {
				publishPack.Identifier = c.GetNewIdentifier()
			}
			protocol.Encode(publishPack, c)
		}

	}

	topicSlice := strings.Split(publishPack.TopicName, "/")

	nodes := sub.GetWildCards(topicSlice)

	for _, node := range nodes {

		node.Clients.Mu.RLock()
		for clientId, _ := range node.Clients.M {
			if c, ok := client.Manage.GetClient(clientId); ok {
				if publishPack.Qos != 0 {
					//fixme 当消息不是第一次发送时应使用之前的idf
					publishPack.Identifier = c.GetNewIdentifier()
				}
				protocol.Encode(publishPack, c)
			}
		}
		node.Clients.Mu.RUnlock()
	}

	broadPub, _ := distributed.NewBroadcastPubMessage()

	broadPub.TopicName = publishPack.TopicName
	broadPub.Qos = publishPack.Qos
	broadPub.Dup = publishPack.Dup
	broadPub.Retain = publishPack.Retain
	broadPub.Payload = publishPack.Payload

	fmt.Println("get pub message", string(publishPack.Payload))
	distributed.BroadcastPub(broadPub)

}
