package main

import (
	"time"

	"github.com/labstack/echo"
)

type roomState int

const (
	ROOM_STATE_INIT roomState = iota - 1
	ROOM_STATE_OPENED
	ROOM_STATE_CLOSED
)

type Room struct {
	// roomId
	roomId string

	// state
	currentState roomState

	// clients
	clients []*client

	// created at
	createdAt time.Time

	// updated at
	updatedAt time.Time
}

func NewRoom(roomId string) *Room {
	return &Room{
		roomId:       roomId,
		currentState: ROOM_STATE_INIT,
		clients:      make([]*client, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func (r *Room) addClient(c echo.Context) error {
	newClient, err := NewClient(c, r)
	if err != nil {
		return err
	}
	r.clients = append(r.clients, newClient)
	newClient.Listen()
	return nil
}

func (r *Room) removeClient(client *client) {

}
