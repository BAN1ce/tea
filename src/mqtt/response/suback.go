package response

import "tea/src/utils"

type Suback struct {
	identifier uint16
	payload    []byte
}

func NewSuback(identifier uint16, qos []byte) *Suback {

	s := new(Suback)
	s.identifier = identifier
	s.payload = qos
	return s
}

func (s Suback) GetFixedHeaderWithoutLength() byte {
	return 0x90
}

func (s Suback) GetVariableHeader() ([]byte, bool) {

	return utils.Uint16ToBytes(s.identifier), true
}

func (s Suback) GetPayload() ([]byte, bool) {
	return s.payload, true
}
