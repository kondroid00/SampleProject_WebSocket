package main

import (
	"net/http"
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
	PING_WAIT = 5 * time.Second
	PONG_WAIT = 5 * time.Second
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
}

func NewClient(c echo.Context, room *Room) (*Client, error) {
	socket, err := upgrater.Upgrade(c.Response(), c.Request(), nil)
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
}

func (c *Client) listen() {
	defer func() {
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

func (c *Client) end() {
	c.room.removeClient(c)
}

func (c *Client) receiveMessage(msg []byte) {
	/*
		data := string(msg)
		prefix := string([]rune(data)[:3])
		payload := string([]rune(data)[3:])
		log.Printf("prefix = " + prefix)
		log.Printf("payload = " + payload)

		switch prefix {
		case CONSTANTS_MSGPREFIX_JOINED:

		case CONSTANTS_MSGPREFIX_REMOVED:

		case CONSTANTS_MSGPREFIX_MESSAGE:

		}

		c.write(msg)
	*/
	c.room.broadcast(msg)
}

func (c *Client) write(msg []byte) {
	err := c.socket.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		c.setClose()
	}
}
