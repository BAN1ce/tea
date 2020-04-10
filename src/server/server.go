package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net"
	"sync"
	"tea/src/client"
	"tea/src/unpack"
	"time"
)

type Server struct {
	name          string
	port          int
	isStop        bool
	clients       sync.Map
	clientDone    chan uuid.UUID
	listener      net.Listener
	protocol      unpack.Protocol
	mutex         *sync.RWMutex
	onServerStart func(s *Server)
	onConnect     client.OnConnect
	onMessage     client.OnMessage
	onClose       client.OnClose
	hbContent     []byte
	hbInterval    time.Duration
	hbTimeout     time.Duration
}

func NewServer(addr net.Addr) *Server {

	var err error
	server := new(Server)
	server.isStop = true
	server.listener, err = net.Listen(addr.Network(), addr.String())

	server.clientDone = make(chan uuid.UUID, 1000)
	server.mutex = new(sync.RWMutex)

	if err != nil {
		panic("Server Listen on " + addr.String() + " FAIL" + err.Error())
	}
	fmt.Println("Server Listen on " + addr.String() + " SUCCESS")
	return server
}

func (s *Server) SetOnConnect(onConnect client.OnConnect) {
	s.onConnect = onConnect
}

func (s *Server) SetOnMessage(onMessage client.OnMessage) {
	s.onMessage = onMessage
}

func (s *Server) SetOnClose(onClose client.OnClose) {
	s.onClose = onClose
}
func (s *Server) SetProtocol(protocol unpack.Protocol) {
	s.protocol = protocol
}
func (s *Server) SetHbInterval(d time.Duration) {
	s.hbInterval = d
}
func (s *Server) SetHbTimeout(d time.Duration) {
	s.hbTimeout = d
}
func (s *Server) SetHbContent(content []byte) {
	s.hbContent = content
}
func (s *Server) Run(ctx context.Context) {

	s.mutex.Lock()

	if s.isStop == false {
		fmt.Println("Server is Running")
	} else {
		if s.onServerStart != nil {
			s.onServerStart(s)
		}
		s.isStop = true
		if s.protocol == nil {
			panic("Set protocol first please")
		}
	}
	s.mutex.Unlock()

	s.mutex.RLock()

	if s.isStop == true {
		go func() {
			for {
				select {

				case uid := <-s.clientDone:
					s.clients.Delete(uid.String())
				}
			}
		}()

		go func() {
			for {
				conn, err := s.listener.Accept()
				if err != nil {
					fmt.Println(err)
					continue
				}
				go s.acceptHandle(ctx, conn)

			}
		}()

	}

	s.mutex.RUnlock()

}

func (s *Server) acceptHandle(ctx context.Context, conn net.Conn) {

	uid := uuid.New()

	c := client.NewClient(conn, uid, s.clientDone, s.protocol, s.onMessage, s.onConnect, s.onClose)

	s.clients.Store(uid, conn)

	if s.hbInterval > 0 && s.hbContent != nil {
		c.SetHbContent(s.hbContent)
		c.SetHbInterval(s.hbInterval)
	}

	if s.hbTimeout > 0 {
		c.SetHbTimeout(s.hbTimeout)
	}

	c.Run(ctx)

}
