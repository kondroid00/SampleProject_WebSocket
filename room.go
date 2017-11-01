package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/kondroid00/SampleProject_WebSocket/dto"
	"github.com/labstack/echo"
)

type roomState int

const (
	ROOM_STATE_INIT roomState = iota
	ROOM_STATE_OPENED
	ROOM_STATE_CLOSED
)

type Room struct {
	// roomId
	roomId string

	// state
	currentState roomState

	// clients
	clients []*Client

	// created at
	createdAt time.Time

	// updated at
	updatedAt time.Time

	// client no
	clientNo int

	// mutex
	mutex sync.Mutex

	// adapter
	adapter *RoomAdapter

	// returner
	returner *RoomReturner
}

type RoomAdapter struct {
	addClientChan    chan *echo.Context
	removeClientChan chan *Client
	sendMessageChan  chan []byte
	sendInfoChan     chan *InfoData
}

type RoomReturner struct {
	addClientChan chan error
}

type MessageData struct {
	prefix  MsgPrefix
	message []byte
}

type InfoData struct {
	prefix MsgPrefix
	client *Client
}

func NewRoom(roomId string) *Room {
	return &Room{
		roomId:       roomId,
		currentState: ROOM_STATE_INIT,
		clients:      make([]*Client, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		clientNo:     0,
		adapter: &RoomAdapter{
			addClientChan:    make(chan *echo.Context),
			removeClientChan: make(chan *Client),
			sendMessageChan:  make(chan []byte),
			sendInfoChan:     make(chan *InfoData),
		},
		returner: &RoomReturner{
			addClientChan: make(chan error),
		},
	}
}

func (r *Room) run() {
	go r.update()
}

func (r *Room) setOpen() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.currentState = ROOM_STATE_OPENED
}

func (r *Room) setClose() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.currentState = ROOM_STATE_CLOSED
}

//------------------------------------------------------------------------------------------------------
// main goroutin
//------------------------------------------------------------------------------------------------------
func (r *Room) update() {
	defer func() {
		log.Printf("room %s update goroutin break", r.roomId)
	}()

	ticker := time.NewTicker(time.Second)
loop:
	for {
		select {
		case context := <-r.adapter.addClientChan:
			r.addClientImpl(context)
		case client := <-r.adapter.removeClientChan:
			r.removeClientImpl(client)
		case data := <-r.adapter.sendMessageChan:
			r.sendMessageImpl(data)
		case data := <-r.adapter.sendInfoChan:
			r.sendInfoImpl(data)
		case <-ticker.C:
			if r.currentState == ROOM_STATE_CLOSED {
				break loop
			}
		}
	}
}

func (r *Room) addClientImpl(c *echo.Context) {
	r.clientNo++
	newClient, err := NewClient(c, r)
	if err != nil {
		r.returner.addClientChan <- err
	}
	r.clients = append(r.clients, newClient)
	newClient.run()
	r.returner.addClientChan <- nil
}

func (r *Room) removeClientImpl(client *Client) {
	r.sendInfoImpl(&InfoData{
		prefix: MSGPREFIX_REMOVED,
		client: client})
	newClients := make([]*Client, 0, (len(r.clients) - 1))
	for _, c := range r.clients {
		if c != client {
			newClients = append(newClients, c)
		}
	}
	r.clients = newClients

	if len(r.clients) == 0 {
		r.setClose()
		RoomManagerInstance().RemoveRoom(r)
	}
}

/*
func (r *Room) broadcastImpl(data *MessageData) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, c := range r.clients {
		c.writeWithPrefix(data.prefix, data.message)
	}
}
*/

func (r *Room) sendMessageImpl(message []byte) {
	for _, c := range r.clients {
		c.write(message)
	}
}

// prefix must be MSGPREFIX_JOINED or MSGPREFIX_REMOVED
func (r *Room) sendInfoImpl(data *InfoData) {
	for _, sender := range r.clients {
		clients := make([]*dto.Client, 0, len(r.clients))
		for _, c := range r.clients {

			action := false
			if c == data.client {
				action = true
			}

			self := false
			if c == sender {
				self = true
			}

			clientDto := &dto.Client{
				ClientNo: c.clientNo,
				Name:     c.userName,
				Action:   action,
				Self:     self,
			}
			clients = append(clients, clientDto)
		}
		jsonByte, _ := json.Marshal(struct {
			Clients []*dto.Client `json:"clients"`
		}{
			Clients: clients,
		})
		sender.writeWithPrefix(data.prefix, jsonByte)
	}
}
