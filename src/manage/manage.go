package manage

import (
	"context"
	"github.com/google/uuid"
	"net"
	"sync"
	"tea/src/mqtt/sub"
	"tea/src/unpack"
	"time"
)

type Manage struct {
	clients    sync.Map
	clientDone chan uuid.UUID
	OnConnect  OnConnect
	OnMessage  OnMessage
	OnClose    OnClose
	Protocol   unpack.Protocol
	HbTimer    *time.Timer
	HbTimeout  time.Duration
	HbInterval time.Duration
	HbContent  []byte
}

func NewManage(onConnect OnConnect, onMessage OnMessage, onClose OnClose, protocol unpack.Protocol) *Manage {
	m := new(Manage)
	m.clientDone = make(chan uuid.UUID, 10)
	m.Protocol = protocol
	m.OnClose = onClose
	m.OnConnect = onConnect
	m.OnMessage = onMessage
	return m
}

func (m *Manage) Run() {

	go func() {

		for {

			select {
			case clientId := <-m.clientDone:
				//fixme 删除客户端的订阅信息
				if client, ok := m.clients.Load(clientId); ok {
					if c, ok := client.(*Client); ok {
						for topic, _ := range c.Topics {
							sub.DeleteHashSub(topic, clientId)
						}
					}
				}
				m.clients.Delete(clientId)

			}
		}
	}()
}

func (m *Manage) AddClient(ctx context.Context, conn net.Conn) {

	uid := uuid.New()

	c := NewClient(conn, uid, m.clientDone, m.Protocol, m.OnMessage, m.OnConnect, m.OnClose, m)

	m.clients.Store(uid, c)

	if hbInterval, hbContent := m.HbInterval, m.HbContent;
		hbInterval > 0 && hbContent != nil {
		c.SetHbContent(hbContent)
		c.SetHbInterval(hbInterval)
	}

	if hbTimout := m.HbTimeout; hbTimout > 0 {
		c.SetHbTimeout(hbTimout)
	}
	c.Run(ctx)
}

func (m *Manage) SendToClient(clientId uuid.UUID, data []byte) {

	if c, ok := m.clients.Load(clientId); ok {

		if c, ok := c.(*Client); ok {
			c.Write(data)
		}
	}
}

func (m *Manage) GetClient(clientId uuid.UUID) (*Client, bool) {

	if c, ok := m.clients.Load(clientId); ok {
		if c, ok := c.(*Client); ok {
			return c, true
		}
	}
	return nil, false

}
