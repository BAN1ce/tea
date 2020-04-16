package response

type Response interface {
	GetCmd() byte
	GetFixedHeaderWithoutLength() byte
	GetVariableHeader() ([]byte, bool)
	GetPayload() ([]byte, bool)
}
