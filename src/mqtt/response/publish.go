package response

import (
	"encoding/binary"
	"tea/src/utils"
)

type Publish struct {
	dup        int
	qos        int
	retain     int
	topicName  []byte
	Identifier uint16
	payload    []byte
}

func NewPublish(dup, qos, retain int, topicName []byte, payload []byte) Publish {

	return Publish{
		dup:       dup,
		qos:       qos,
		retain:    retain,
		topicName: topicName,
		payload:   payload,
	}
}

func (p Publish) GetFixedHeaderWithoutLength() byte {
	return 0x30
}

func (p Publish) GetVariableHeader() ([]byte, bool) {

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(len(p.topicName)))
	buf = append(buf, p.topicName...)

	if p.qos != 0 {
		buf = append(buf, utils.Uint16ToBytes(p.Identifier)...)
		return buf, true
	} else {
		return buf, true
	}
}

func (p Publish) GetPayload() ([]byte, bool) {
	return p.payload, true
}
