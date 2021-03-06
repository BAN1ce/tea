package response

const ACCEPT = 1
const UNSUPPORT_PROCOTOL = 2
const INVLIAD_CLIENT = 3
const SERVER_FAIL = 4
const UNAUTH = 5

type Connack struct {
	ackSign    uint8
	returnCode uint8
}

func NewConnack(returnCode int) *Connack {

	c := new(Connack)
	c.ackSign = 0 & 0
	c.returnCode = 0 | c.returnCode
	return c
}

func (c Connack) GetFixedHeaderWithoutLength() byte {

	return 0x20
}

func (c Connack) GetVariableHeader() ([]byte, bool) {

	tmp := make([]byte, 2)
	tmp[0] = c.ackSign
	tmp[1] = c.returnCode
	return tmp, true
}

func (c Connack) GetPayload() ([]byte, bool) {

	return nil, false
}
