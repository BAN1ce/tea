package handle

import (
	"tea/src/manage"
	"tea/src/mqtt/protocol"
	"tea/src/mqtt/request"
	"tea/src/mqtt/response"
	"tea/src/utils"
)

type Connect struct {
}

func NewConnect() *Connect {
	return new(Connect)
}

func (c *Connect) Handle(pack protocol.Pack, client *manage.Client) {

	connectPack := request.NewConnectPack(pack)
	body := connectPack.Data[connectPack.BodyStart:]
	// Reserved 不为0时 断开客户端
	if connectPack.Reserved != 0 {
		client.Stop()
	}
	plc := 0
	clientIdLength := utils.UtfLength(body[0:2])
	plc += 2
	connectPack.ClientId = string(body[2 : clientIdLength+plc])
	plc += clientIdLength
	if connectPack.WillFlag == 1 {
		willTopicLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.WillTopic = string(body[plc : willTopicLength+plc])
		plc += willTopicLength
		willMessageLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.WillMessage = string(body[plc : willMessageLength+plc])
		plc += willMessageLength
	}
	if connectPack.UserNameFlag == 1 {
		userNameLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.UserName = string(body[plc : plc+userNameLength])
		plc += userNameLength
	}
	if connectPack.PasswordFlag == 1 {
		passwordLength := utils.UtfLength(body[plc : plc+2])
		plc += 2
		connectPack.Password = string(body[plc : plc+passwordLength])
		plc += passwordLength
	}

	connack := response.NewConnack(response.ACCEPT)

	protocol.Encode(connack, client)
	//todo 用户名密码进行验证，为连接设置clientID，为客户端开辟session

}
