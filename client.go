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
	//write chan []byte

	// read chan
	//read chan []byte

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
		//write:        make(chan []byte, messageBufferSize),
		//read: make(chan []byte, messageBufferSize),
		room: room,
	}
	return client, nil
}

func (c *client) setClose() {
	c.currentState = CLIENT_STATE_CLOSED
}

func (c *client) run() {
	go c.listen()
}

func (c *client) listen() {
	defer func() {
		c.end()
	}()

loop:
	for {
		/*
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
				println(time.String())
			}
		*/

		// read

		messageType, readMsg, err := c.socket.ReadMessage()
		if err != nil {
			break loop
		}
		switch messageType {
		case websocket.TextMessage:
			//c.read <- readMsg
			c.receiveMessage(readMsg)
		default:

		}

		// read end

	}
}

func (c *client) end() {

}

func (c *client) receiveMessage(msg []byte) {
	c.write(msg)
}

func (c *client) write(msg []byte) {
	err := c.socket.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		c.setClose()
	}
}
