package server

import (
	"fmt"
	"github.com/google/uuid"
	"net"
	"sync"
)

type Server struct {
	name       string
	port       int
	isStop     bool
	clients    sync.Map
	clientDone chan uuid.UUID
	listener   net.Listener
	mutex      *sync.RWMutex
}

func NewServer(addr net.Addr) *Server {

	var err error
	server := new(Server)
	server.listener, err = net.Listen(addr.Network(), addr.String())

	server.clientDone = make(chan uuid.UUID, 1000)
	server.mutex = new(sync.RWMutex)

	if err != nil {
		panic("Server Listen on " + addr.String() + " FAIL" + err.Error())
	}
	fmt.Println("Server Listen on " + addr.String() + " SUCCESS")
	return server
}

func (s *Server) Run() {

	s.mutex.Lock()

	if s.isStop == false {
		fmt.Println("Server is Running")
	} else {
		s.isStop = true
	}
	s.mutex.Unlock()

	s.mutex.RLock()

	if s.isStop == true {

		//todo server 启动
		go func() {
			for {
				select {

				case uid := <-s.clientDone:
					s.clients.Delete(uid.String())
				}
			}
		}()

	}

	s.mutex.RUnlock()

}

func (s *Server) acceptHandle(conn *net.Conn) {

	//todo onConnect 事件
	
}
