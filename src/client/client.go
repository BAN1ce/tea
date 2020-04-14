package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
	"tea/src/unpack"
	"time"
)

type OnMessage func([]byte, *Client) error
type OnClose func(*Client) error
type OnConnect func(*Client) error

type Client struct {
	conn       net.Conn
	isStop     bool
	writeChan  chan []byte
	readChan   chan []byte
	Uid        uuid.UUID
	clientDone chan<- uuid.UUID
	hbTimer    *time.Timer
	hbTimeout  time.Duration
	hbInterval time.Duration
	hbContent  []byte
	mutex      *sync.RWMutex
	cancel     context.CancelFunc
	protocol   unpack.Protocol
	onMessage  OnMessage
	onConnect  OnConnect
	onClose    OnClose
}

/**
create a new client .
*/
func NewClient(conn net.Conn, uuid uuid.UUID, clientDone chan<- uuid.UUID, protocol unpack.Protocol,
	onMessage OnMessage, onConnect OnConnect, onClose OnClose) *Client {

	client := &Client{
		conn:       conn,
		isStop:     true,
		writeChan:  make(chan []byte, 10),
		readChan:   make(chan []byte, 10),
		Uid:        uuid,
		clientDone: clientDone,
		mutex:      new(sync.RWMutex),
		protocol:   protocol,
		onMessage:  onMessage,
		onConnect:  onConnect,
		onClose:    onClose,
	}

	return client
}

func (c *Client) SetHbInterval(d time.Duration) {
	c.hbInterval = d
}

func (c *Client) SetHbTimeout(d time.Duration) {
	c.hbTimeout = d
}

func (c *Client) SetHbContent(content []byte) {
	c.hbContent = content
}

/**
let client run .
*/
func (c *Client) Run(ctx context.Context) {

	c.mutex.Lock()
	if c.isStop == true {
		// client onConnect event
		if c.onConnect != nil {
			if err := c.onConnect(c); err != nil {
				log.Println(err)
			}
		}
		if c.hbInterval > 0 {
			c.hbTimer = time.NewTimer(c.hbInterval)
		}
		ctxChild, cancel := context.WithCancel(ctx)
		c.cancel = cancel

		go c.unPackHandle(ctxChild)
		go c.readHandle(ctxChild)
		go c.writeHandle(ctxChild)

		go c.heartBeatHandle(ctxChild)

		c.isStop = false
	}
	c.mutex.Unlock()
	return
}

/**
stop client , execute cancel function , close read,write chan, send a message to clientDone chan.
*/
func (c *Client) Stop() {

	c.mutex.Lock()
	if c.isStop == false {
		// client onClose event
		if c.onClose != nil {
			if err := c.onClose(c); err != nil {
				log.Println(err)
			}
		}
		c.cancel()
		c.clientDone <- c.Uid
		conn := c.conn
		conn.Close()
		c.isStop = true
	}
	c.mutex.Unlock()

}
func (c *Client) Write(msg []byte) {
	fmt.Println("write", msg)
	c.writeChan <- msg

}

func (c *Client) unPackHandle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			close(c.readChan)
			return

		default:
			scanner := bufio.NewScanner(c.conn)
			scanner.Split(c.protocol())
			for scanner.Scan() {
				c.readChan <- scanner.Bytes()
				if c.hbTimeout > 0 {
					err := c.conn.SetReadDeadline(time.Now().Add(c.hbTimeout))
					if err != nil {
						continue
					}
				}
			}
			if err := scanner.Err(); err != nil {
				c.Stop()
			}
		}
	}
}

func (c *Client) readHandle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-c.readChan:

			if ok {
				// client onMessage event
				if c.onMessage != nil {
					if err := c.onMessage(msg, c); err != nil {
						log.Println(err)
					}
				}
			} else {
				return
			}


		}
	}

}
func (c *Client) writeHandle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			close(c.writeChan)
			return
		case m, ok := <-c.writeChan:

			if ok {
				conn := c.conn

				if _, err := conn.Write(m); err != nil {
					log.Println(err)
				}
			}



		}
	}
}

func (c *Client) heartBeatHandle(ctx context.Context) {

	if c.hbInterval == 0 {
		return
	}
	for {
		select {
		case <-ctx.Done():
			c.hbTimer.Stop()
			return

		case <-c.hbTimer.C:
			//todo 发送心跳包到write chan ，重制定时器

			c.writeChan <- []byte{'&'}
			if c.hbInterval > 0 {
				c.hbTimer.Reset(c.hbInterval)
			}

		}
	}
}
