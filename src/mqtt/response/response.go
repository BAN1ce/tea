package response

type Response interface {
	GetFixedHeaderWithoutLength() byte
	GetVariableHeader() ([]byte, bool)
	GetPayload() ([]byte, bool)
}
