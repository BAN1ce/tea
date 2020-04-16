package response

type Hback struct {
}

func NewHback() *Hback {

	return new(Hback)
}


func (h Hback) GetFixedHeaderWithoutLength() byte {
	return 0xe0
}

func (h Hback) GetVariableHeader() ([]byte, bool) {

	return nil, false
}

func (h Hback) GetPayload() ([]byte, bool) {
	return nil, false
}
