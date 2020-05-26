package manage

import (
	"context"
	"github.com/google/uuid"
	"net"
	"strings"
	"sync"
	"tea/src/mqtt/sub"
	"tea/src/unpack"
	"tea/src/utils"
	"time"
)

type Manage struct {
	clients            map[uuid.UUID]*Client
	ClientsUsernameUid sync.Map
	clientDone         chan uuid.UUID
	OnConnect          OnConnect
	OnMessage          OnMessage
	OnClose            OnClose
	Protocol           unpack.Protocol
	HbTimer            *time.Timer
	HbTimeout          time.Duration
	HbInterval         time.Duration
	HbContent          []byte
	mu                 sync.RWMutex
}

var LocalManage *Manage

func NewManage(onConnect OnConnect, onMessage OnMessage, onClose OnClose, protocol unpack.Protocol) *Manage {
	m := new(Manage)
	m.clients = make(map[uuid.UUID]*Client, 50000)
	m.clientDone = make(chan uuid.UUID, 10)
	m.Protocol = protocol
	m.OnClose = onClose
	m.OnConnect = onConnect
	m.OnMessage = onMessage
	m.HbTimeout = 60 * time.Second
	LocalManage = m
	return m
}

func (m *Manage) Run() {

	go func() {

		for {

			select {
			case clientId := <-m.clientDone:
				//fixme 删除客户端的两种订阅
				if client, ok := m.GetClient(clientId); ok {
					for topic, _ := range client.Topics {
						topics := strings.Split(topic, "/")
						if utils.HasWildcard(topics) {
							sub.DeleteTreeSub(topics, clientId)
						} else {
							sub.DeleteHashSub(topic, clientId)
						}

					}
					m.ClientsUsernameUid.Delete(client.UserName)
				}
				m.delClient(clientId)

			}
		}
	}()
}

func (m *Manage) AddClient(ctx context.Context, conn net.Conn) {

	uid := uuid.New()

	c := NewClient(conn, uid, m.clientDone, m.Protocol, m.OnMessage, m.OnConnect, m.OnClose, m)

	m.setClient(uid, c)

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

	if c, ok := m.GetClient(clientId); ok {
		c.Write(data)
	}
}

func (m *Manage) GetClient(clientId uuid.UUID) (*Client, bool) {

	m.mu.RLock()
	c, ok := m.clients[clientId]
	defer m.mu.RUnlock()

	return c, ok
}

func (m *Manage) setClient(clientId uuid.UUID, client *Client) {

	m.mu.Lock()
	m.clients[clientId] = client
	m.mu.Unlock()
}

func (m *Manage) delClient(clientId uuid.UUID) {
	m.mu.Lock()
	delete(m.clients, clientId)
	m.mu.Unlock()

}

func (m *Manage) GetClientCount() int {

	m.mu.RLock()
	count := len(m.clients)
	m.mu.RUnlock()
	return count
}
