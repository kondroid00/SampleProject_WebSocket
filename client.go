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
	CLIENT_STATE_INIT clientState = iota - 1
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

type client struct {
	// websocket
	socket *websocket.Conn

	// state
	currentState clientState

	// userId
	userId string

	// userNam
	userName string

	// write chan
	write chan []byte

	// read chan
	read chan []byte

	// room
	room *Room
}

func NewClient(c echo.Context, room *Room) (*client, error) {
	socket, err := upgrater.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, err
	}
	client := &client{
		socket:       socket,
		currentState: CLIENT_STATE_INIT,
		userId:       "",
		userName:     "",
		write:        make(chan []byte, messageBufferSize),
		read:         make(chan []byte, messageBufferSize),
		room:         room,
	}
	return client, nil
}

//state
func (c *client) changeState(nextState clientState) {
	newState := nextState
	oldState := c.currentState

}

func (c *client) onStateEnter(newState clientState) {
	switch newState {
	case CLIENT_STATE_INIT:

	case CLIENT_STATE_OPENED:

	case CLIENT_STATE_CLOSED:

	default:
		return
	}
}

func (c *client) onStateExit(oldState clientState) {
	switch oldState {
	case CLIENT_STATE_INIT:

	case CLIENT_STATE_OPENED:

	case CLIENT_STATE_CLOSED:

	default:
		return
	}
}

func (c *client) Listen() {
	go c.listen()
}

func (c *client) listen() {
	defer func() {
		c.end()
	}()

	ticker := time.NewTicker(time.Second)
loop:
	for {
		// read
		messageType, readMsg, err := c.socket.ReadMessage()
		if err != nil {
			break loop
		}
		switch messageType {
		case websocket.TextMessage:
			c.read <- readMsg
		default:

		}
		// read end

		select {
		// write
		case writeMsg := <-c.write:
			err := c.socket.WriteMessage(websocket.TextMessage, writeMsg)
			if err != nil {
				break loop
			}
		// write end

		// break goroutin if client is closed
		case time := <-ticker.C:
			if c.currentState == CLIENT_STATE_CLOSED {
				break loop
			}
		}
	}
}

func (c *client) end() {

}
