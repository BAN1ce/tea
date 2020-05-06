package request

import (
	"tea/src/mqtt/protocol"
	"tea/src/utils"
)

type PublishPack struct {
	protocol.Pack
	Dup        int
	Qos        int
	Retain     int
	TopicName  string
	Identifier uint16
	Payload    []byte
}

func (A PublishPack) GetCmd() byte {
	return 0x3
}

func (A PublishPack) GetFixedHeaderWithoutLength() byte {

	return 0x3 << 4
}

func (A PublishPack) GetVariableHeader() ([]byte, bool) {
	return append([]byte(A.TopicName), utils.Uint16ToBytes(A.Identifier)...), true
}

func (A PublishPack) GetPayload() ([]byte, bool) {
	return A.Payload, true
}
