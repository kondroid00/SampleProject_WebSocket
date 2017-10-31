package main

import (
	"encoding/json"
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
}

func NewRoom(roomId string) *Room {
	return &Room{
		roomId:       roomId,
		currentState: ROOM_STATE_INIT,
		clients:      make([]*Client, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		clientNo:     0,
	}
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

func (r *Room) addClient(c echo.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clientNo++
	newClient, err := NewClient(c, r)
	if err != nil {
		return err
	}
	r.clients = append(r.clients, newClient)
	newClient.run()
	return nil
}

func (r *Room) removeClient(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sendInfo(MSGPREFIX_REMOVED, client)
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

func (r *Room) broadcast(msg []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, c := range r.clients {
		c.write(msg)
	}
}

func (r *Room) broadcastWithPrefix(prefix MsgPrefix, msg []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, c := range r.clients {
		c.writeWithPrefix(prefix, msg)
	}
}

// prefix must be MSGPREFIX_JOINED or MSGPREFIX_REMOVED
func (r *Room) sendInfo(prefix MsgPrefix, client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, sender := range r.clients {
		clients := make([]*dto.Client, 0, len(r.clients))
		for _, c := range r.clients {

			action := false
			if c == client {
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
		sender.writeWithPrefix(prefix, jsonByte)
	}
}
