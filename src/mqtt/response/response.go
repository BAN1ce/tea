package response

type Response interface {
	GetCmd() byte
	GetFixedHeaderWithLength() byte
	GetVariableHeader() ([]byte, bool)
	GetPayload() ([]byte, bool)
}
