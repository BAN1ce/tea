package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net"
	"sync"
	"tea/src/manage"
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
	Protocol      unpack.Protocol
	mutex         *sync.RWMutex
	onServerStart func(s *Server)
	OnConnect     manage.OnConnect
	OnMessage     manage.OnMessage
	OnClose       manage.OnClose
	hbContent     []byte
	hbInterval    time.Duration
	hbTimeout     time.Duration
	Manage        *manage.Manage
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

	fmt.Println(`
  _______ ______          
 |__   __|  ____|   /\    
    | |  | |__     /  \   
    | |  |  __|   / /\ \  
    | |  | |____ / ____ \ 
    |_|  |______/_/    \_\
	`)

	fmt.Println("Server Listen on " + addr.String() + " SUCCESS")
	return server
}
func (s *Server) GetHbInterval() time.Duration {

	return s.hbInterval
}

func (s *Server) GetHbTimeout() time.Duration {

	return s.hbTimeout
}
func (s *Server) GetHbContent() []byte {

	return s.hbContent

}
func (s *Server) SetOnConnect(onConnect manage.OnConnect) {
	s.OnConnect = onConnect
}

func (s *Server) SetOnMessage(onMessage manage.OnMessage) {
	s.OnMessage = onMessage
}

func (s *Server) SetOnClose(onClose manage.OnClose) {
	s.OnClose = onClose
}
func (s *Server) SetProtocol(protocol unpack.Protocol) {
	s.Protocol = protocol
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
		s.Manage = manage.NewManage(s.OnConnect, s.OnMessage, s.OnClose, s.Protocol)
		if s.onServerStart != nil {
			s.onServerStart(s)
		}
		s.isStop = true
		if s.Protocol == nil {
			panic("Set Protocol first please")
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
	s.Manage.Run()

	s.mutex.RUnlock()

}

func (s *Server) acceptHandle(ctx context.Context, conn net.Conn) {

	s.Manage.AddClient(ctx, conn)

}
