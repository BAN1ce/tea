package request

import (
	"tea/src/mqtt/protocol"
	"tea/src/utils"
)

type ConnectPack struct {
	protocol.Pack
	BodyStart     int
	ProtocolName  string
	ProtocolLevel int
	ConnectFlags  byte
	Reserved      int
	CleanSession  int
	WillFlag     int
	WillQos      int
	WillRetain   int
	PasswordFlag int
	UserNameFlag int
	ClientId     string
	WillTopic    string
	WillMessage  string
	UserName     string
	Password     string
	KeepAlive    int
}

func NewConnectPack(pack protocol.Pack) *ConnectPack {

	c := new(ConnectPack)
	c.Pack = pack
	plc := utils.UtfLength(pack.Data[pack.FixHeaderLength:pack.FixHeaderLength+2]) + 2
	c.ProtocolName = string(pack.Data[pack.FixHeaderLength+2 : pack.FixHeaderLength+plc])
	c.ProtocolLevel = int(pack.Data[pack.FixHeaderLength+plc])
	c.ConnectFlags = pack.Data[pack.FixHeaderLength+plc+1]
	c.Reserved = int(c.ConnectFlags & 1)
	c.CleanSession = int((c.ConnectFlags & 2) >> 1)
	c.WillFlag = int((c.ConnectFlags & 4) >> 2)
	c.WillQos = int((c.ConnectFlags & 24) >> 3)
	c.WillRetain = int((c.ConnectFlags & 32) >> 5)
	c.PasswordFlag = int((c.ConnectFlags & 64) >> 6)
	c.UserNameFlag = int((c.ConnectFlags & 128) >> 7)
	c.KeepAlive = utils.UtfLength(pack.Data[pack.FixHeaderLength+plc+2 : pack.FixHeaderLength+plc+4])
	c.BodyStart = plc + pack.FixHeaderLength + 4

	return c
}

