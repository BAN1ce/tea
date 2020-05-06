package handle

import (
	"strings"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/request"
	"tea/src/mqtt/response"
	"tea/src/mqtt/sub"
	"tea/src/utils"
)

type Subscribe struct {
}

func NewSubscribe() *Subscribe {

	return new(Subscribe)
}

func (s *Subscribe) Handle(pack protocol.Pack, client *manage.Client) {
	p := request.NewSubscribePack(pack)

	qosSlice := make([]byte, 0)
	for topic, qos := range p.TopicQos {
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
	suback := response.NewSuback(p.Identifier, qosSlice)

	protocol.Encode(suback, client)
}
