package handle

import (
	"fmt"
	"github.com/google/uuid"
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

func (c *Connect) Handle(pack *protocol.Pack, client *manage.Client) {

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
	
	// 同名客户端连接则关闭上一个连接
	if existedClient, ok := client.Manage.ClientsUsernameUid.LoadOrStore(connectPack.UserName, client.Uid); ok {
		uid, ok := existedClient.(uuid.UUID)
		if ok {
			fmt.Println(connectPack.UserName, client.Uid, "connect to close")
			existedClient, ok := client.Manage.GetClient(uid)
			if ok {
				// 关闭之前连接
				existedClient.Stop()
				// 新连接添加
				client.Manage.ClientsUsernameUid.Store(connectPack.UserName, client.Uid)
				return
			}
		}
	}

	connack := response.NewConnack(response.ACCEPT)
	client.UserName = connectPack.UserName

	protocol.Encode(connack, client)
	//todo 用户名密码进行验证，为连接设置clientID，为客户端开辟session

}
