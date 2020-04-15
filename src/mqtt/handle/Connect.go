package handle

import (
	"fmt"
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/response"
	"tea/src/utils"
)

type ConnectPack struct {
	protocol.Pack
	bodyStart     int
	protocolName  string
	protocolLevel int
	connectFlags  byte
	reserved      int
	cleanSession  int
	willFlag      int
	willQos       int
	willRetain    int
	passwordFlag  int
	userNameFlag  int
	ClientId      string
	WillTopic     string
	WillMessage   string
	UserName      string
	Password      string
	keepAlive     int
}

func newConnectPack(pack protocol.Pack) *ConnectPack {

	c := new(ConnectPack)
	c.Pack = pack
	plc := utils.UtfLength(pack.Data[pack.FixHeaderLength:pack.FixHeaderLength+2]) + 2
	fmt.Println(plc)
	c.protocolName = string(pack.Data[pack.FixHeaderLength+2 : pack.FixHeaderLength+plc])
	c.protocolLevel = int(pack.Data[pack.FixHeaderLength+plc+1])
	fmt.Println(pack.FixHeaderLength+plc+1, "plc")
	c.connectFlags = pack.Data[pack.FixHeaderLength+plc+1]
	c.reserved = int(c.connectFlags & 1)
	fmt.Println(c.connectFlags)
	c.cleanSession = int((c.connectFlags & 2) >> 1)
	c.willFlag = int((c.connectFlags & 4) >> 2)
	c.willQos = int((c.connectFlags & 24) >> 3)
	c.willRetain = int((c.connectFlags & 32) >> 5)
	c.passwordFlag = int((c.connectFlags & 64) >> 6)
	fmt.Println(int(c.connectFlags & 128))
	c.userNameFlag = int((c.connectFlags & 128) >> 7)
	c.keepAlive = utils.UtfLength(pack.Data[pack.FixHeaderLength+plc+2 : pack.FixHeaderLength+plc+4])
	c.bodyStart = plc + pack.FixHeaderLength + 4

	return c
}

type Connect struct {
}

func NewConnect() *Connect {
	return new(Connect)
}

func (c *Connect) Handle(pack protocol.Pack, client *manage.Client) {

	connectPack := newConnectPack(pack)
	body := connectPack.Data[connectPack.bodyStart:]
	// reserved 不为0时 断开客户端
	if connectPack.reserved != 0 {
		client.Stop()
	}
	plc := 0
	clientIdLength := utils.UtfLength(body[0:2])
	plc += 2
	connectPack.ClientId = string(body[2 : clientIdLength+plc])
	plc += clientIdLength
	if connectPack.willFlag == 1 {
		willTopicLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.WillTopic = string(body[plc : willTopicLength+plc])
		plc += willTopicLength
		willMessageLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.WillMessage = string(body[plc : willMessageLength+plc])
		plc += willMessageLength
	}
	if connectPack.userNameFlag == 1 {
		userNameLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.UserName = string(body[plc : plc+userNameLength])
		plc += userNameLength
	}
	if connectPack.passwordFlag == 1 {
		passwordLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.Password = string(body[plc : plc+passwordLength])
		plc += passwordLength
	}

	connack := response.NewConnack(response.ACCEPT)

	fmt.Println("connect handle")

	protocol.Encode(connack, client)
	//todo 用户名密码进行验证，为连接设置clientID，为客户端开辟session

}
