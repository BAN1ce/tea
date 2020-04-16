package response

import "tea/src/utils"

type Publish struct {
	dup        int
	qos        int
	retain     int
	topicName  []byte
	identifier uint16
	payload    []byte
}

func NewPublish(dup, qos, retain int, topicName []byte, identifier uint16, payload []byte) Publish {

	return Publish{
		dup:        dup,
		qos:        qos,
		retain:     retain,
		topicName:  topicName,
		identifier: identifier,
		payload:    payload,
	}
}

func (p Publish) GetFixedHeaderWithoutLength() byte {
	return 0x3
}

func (p Publish) GetVariableHeader() ([]byte, bool) {

	return append(p.topicName, utils.Uint16ToBytes(p.identifier)...), true
}

func (p Publish) GetPayload() ([]byte, bool) {
	return p.payload, true
}
