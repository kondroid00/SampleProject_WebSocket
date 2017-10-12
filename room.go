package main

import (
	"encoding/json"
	"time"

	"github.com/kondroid00/SampleProject_Socket/dto"
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
	r.currentState = ROOM_STATE_OPENED
}

func (r *Room) setClose() {
	r.currentState = ROOM_STATE_CLOSED
}

func (r *Room) addClient(c echo.Context) error {
	r.clientNo++
	newClient, err := NewClient(c, r)
	if err != nil {
		return err
	}
	r.clients = append(r.clients, newClient)
	newClient.run()
	newClient.setOpen()
	r.sendInfo(CONSTANTS_MSGPREFIX_JOINED, newClient)
	return nil
}

func (r *Room) removeClient(client *Client) {
	r.sendInfo(CONSTANTS_MSGPREFIX_REMOVED, client)
	newClients := make([]*Client, 0, (len(r.clients) - 1))
	for _, c := range r.clients {
		if c != client {
			newClients = append(newClients, c)
		}
	}
	r.clients = newClients

	if len(r.clients) == 0 {
		RoomManagerInstance().RemoveRoom(r)
	}
}

func (r *Room) broadcast(msg []byte) {
	for _, c := range r.clients {
		c.write(msg)
	}
}

// prefix must be CONSTANTS_MSGPREFIX_JOINED or CONSTANTS_MSGPREFIX_REMOVED
func (r *Room) sendInfo(prefix msgPrefix, client *Client) {
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
		data := append([]byte(prefix), jsonByte...)
		sender.write(data)
	}
}
