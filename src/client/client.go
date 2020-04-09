package client

import (
	"context"
	"github.com/google/uuid"
	"net"
	"sync"
)

type Client struct {
	conn       *net.Conn
	isStop     bool
	writeChan  chan []byte
	readChan   chan []byte
	uid        uuid.UUID
	clientDone chan<- uuid.UUID
	hbTimeout  int
	mutex      sync.RWMutex
}

func NewClient(ctx context.Context, conn *net.Conn, uuid uuid.UUID, clientDone chan<- uuid.UUID) *Client {

	client := &Client{
		conn:       conn,
		isStop:     true,
		writeChan:  make(chan []byte, 10),
		readChan:   make(chan []byte, 10),
		uid:        uuid,
		clientDone: clientDone,
	}
	return client
}

func (c *Client) Run() {

	c.mutex.Lock()

	if c.isStop == true {

		//todo client run

		c.isStop = false
	}

	c.mutex.Unlock()
	return

}

func (c *Client) close() {

}

