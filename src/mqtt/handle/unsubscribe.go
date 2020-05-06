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

type UnSubscribe struct {
}

func NewUnSubscribe() *UnSubscribe {

	return new(UnSubscribe)
}

func (u *UnSubscribe) Handle(pack protocol.Pack, client *manage.Client) {

	p := request.NewUnSubscribePack(pack)

	for _, topic := range p.Topics {
		topicSlice := strings.Split(topic, "/")

		//fixme delete wildcards sub 删除模糊订阅
		if utils.HasWildcard(topicSlice) {
			sub.DeleteTreeSub(topicSlice, client.Uid)
		} else {
			sub.DeleteHashSub(topic, client.Uid)
		}

	}

	unSuback := response.NewUnSuback(p.Identifier)

	protocol.Encode(unSuback, client)

}
