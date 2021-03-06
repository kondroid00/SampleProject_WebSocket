package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"

	"github.com/gorilla/websocket"
)

type clientState int

const (
	// state
	CLIENT_STATE_INIT clientState = iota
	CLIENT_STATE_OPENED
	CLIENT_STATE_CLOSED

	// ping pong
	PING_INTERVAL = time.Second
	PING_WAIT     = 5 * time.Second
	PONG_WAIT     = 5 * time.Second
)

// websocket
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrater = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

type Client struct {
	// websocket
	socket *websocket.Conn

	// mutex
	mutex sync.Mutex

	// state
	currentState clientState

	// userId
	userId string

	// userNam
	userName string

	// room
	room *Room

	// client no
	clientNo int

	// end sync once
	endSyncOnce sync.Once
}

func NewClient(c *echo.Context, room *Room) (*Client, error) {
	socket, err := upgrater.Upgrade((*c).Response(), (*c).Request(), nil)
	if err != nil {
		return nil, err
	}
	client := &Client{
		socket:       socket,
		currentState: CLIENT_STATE_INIT,
		userId:       "",
		userName:     "",
		room:         room,
		clientNo:     room.clientNo,
	}
	return client, nil
}

func (c *Client) setOpen() {
	c.currentState = CLIENT_STATE_OPENED
}

func (c *Client) setClose() {
	c.currentState = CLIENT_STATE_CLOSED
}

func (c *Client) run() {
	go c.listen()
	go c.pingPong()
}

func (c *Client) sendError(errorCode ErrorCode) {
	jsonByte, _ := json.Marshal(struct {
		ErrorCode int    `json:"errorCode"`
		Message   string `json:"message"`
	}{
		ErrorCode: int(errorCode),
		Message:   getErrorCodeMessage(errorCode),
	})
	c.writeWithPrefix(MSGPREFIX_ERROR, jsonByte)
}

func (c *Client) writeWithPrefix(prefix MsgPrefix, msg []byte) {
	data := append([]byte(prefix), msg...)
	c.write(data)
}

func (c *Client) write(msg []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	err := c.socket.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		c.setClose()
	}
}

//------------------------------------------------------------------------------------------------------
// listen goroutin
//------------------------------------------------------------------------------------------------------
func (c *Client) listen() {
	defer func() {
		log.Printf("clientNo %d listen goroutin break", c.clientNo)
		c.setClose()
		c.end()
	}()

loop:
	for {
		// read
		messageType, readMsg, err := c.socket.ReadMessage()
		if err != nil {
			break loop
		}
		switch messageType {
		case websocket.TextMessage:
			c.receiveMessage(readMsg)
		default:
		}
	}
}

func (c *Client) receiveMessage(msg []byte) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("Recover %s clientNo %d", "receiveMessage", c.clientNo)
		}
	}()

	data := string(msg)
	prefix := string([]rune(data)[:3])
	payload := string([]rune(data)[3:])
	log.Printf("prefix = " + prefix)
	log.Printf("payload = " + payload)

	switch MsgPrefix(prefix) {
	case MSGPREFIX_JOINED:
		type Data struct {
			Name   string `json:"name"`
			UserId string `json:"userId"`
		}
		data := &Data{}
		err := json.Unmarshal([]byte(payload), &data)
		if err != nil {
			c.sendError(ERRORCODE_JOINED)
		}
		c.userName = data.Name
		c.userId = data.UserId
		c.setOpen()
		c.room.adapter.sendInfoChan <- &InfoData{
			prefix: MSGPREFIX_JOINED,
			client: c,
		}

	case MSGPREFIX_REMOVED:

	case MSGPREFIX_MESSAGE:
		c.room.adapter.sendMessageChan <- msg
	case MSGPREFIX_ERROR:

	}

}

//------------------------------------------------------------------------------------------------------
// listen goroutin
//------------------------------------------------------------------------------------------------------
func (c *Client) pingPong() {
	defer func() {
		log.Printf("clientNo %d pingPong goroutin break", c.clientNo)
		err := recover()
		if err != nil {
			log.Printf("Recover %s clientNo %d", "pingPong", c.clientNo)
		}
	}()

	ticker := time.NewTicker(PING_INTERVAL)
	c.socket.SetReadDeadline(time.Now().Add(PONG_WAIT))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
loop:
	for {
		select {
		case <-ticker.C:
			if c.currentState == CLIENT_STATE_CLOSED {
				break loop
			}
			if err := c.socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				break loop
			}
			c.socket.SetWriteDeadline(time.Now().Add(PING_WAIT))
		}
	}
}

// May be called in all goroutin
func (c *Client) end() {
	log.Printf("end clientNo = %d", c.clientNo)
	c.endSyncOnce.Do(func() {
		c.room.adapter.removeClientChan <- c
	})
}
