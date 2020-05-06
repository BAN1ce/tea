package request

import (
	"encoding/binary"
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

	return A.FixedHeader[0]
}

func (A PublishPack) GetVariableHeader() ([]byte, bool) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(len(A.TopicName)))
	buf = append(buf, []byte(A.TopicName)...)
	if A.Qos != 0 {
		buf = append(buf, utils.Uint16ToBytes(A.Identifier)...)
	}
	return buf, true
}

func (A PublishPack) GetPayload() ([]byte, bool) {
	return A.Payload, true
}
