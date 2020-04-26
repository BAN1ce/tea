package response

import "tea/src/utils"

type UnSuback struct {
	identifier uint16
	payload    []byte
}

func NewUnSuback(identifier uint16) *UnSuback {

	s := new(UnSuback)
	s.identifier = identifier
	return s
}

func (u UnSuback) GetFixedHeaderWithoutLength() byte {
	return 0xb0
}

func (u UnSuback) GetVariableHeader() ([]byte, bool) {

	return utils.Uint16ToBytes(u.identifier), true
}

func (u UnSuback) GetPayload() ([]byte, bool) {

	return nil, false
}
