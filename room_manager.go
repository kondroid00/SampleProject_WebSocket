package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type RoomManager struct {
	rooms []*Room
}

// singleton
var instance *RoomManager = newRoomManager()

func newRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make([]*Room, 0),
	}
}

func RoomManagerInstance() *RoomManager {
	return instance
}

func (rm *RoomManager) Serve(c echo.Context) error {
	roomId := c.Param("id")
	for _, room := range rm.rooms {
		if room.roomId == roomId {
			if room.currentState == ROOM_STATE_OPENED {
				err := room.addClient(c)
				if err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
			}
		}
	}

	return c.String(http.StatusInternalServerError, "Error")
}
