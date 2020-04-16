package response

import "tea/src/utils"

type Suback struct {
	cmd        uint8
	identifier uint16
	payload    []byte
}

func NewSuback(identifier uint16, qos []byte) *Suback {

	s := new(Suback)
	s.cmd = 0x90
	s.identifier = identifier
	s.payload = qos
	return s
}

func (s Suback) GetCmd() byte {

	return s.cmd
}

func (s Suback) GetFixedHeaderWithoutLength() byte {
	return s.cmd
}

func (s Suback) GetVariableHeader() ([]byte, bool) {

	return utils.Uint16ToBytes(s.identifier), true
}

func (s Suback) GetPayload() ([]byte, bool) {
	return s.payload, true
}
