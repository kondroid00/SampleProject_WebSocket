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
				if err := room.addClient(c); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				return c.String(http.StatusOK, "OK")
			}
		}
	}

	// create a new room if the room does not exists
	newRoom := NewRoom(roomId)
	rm.rooms = append(rm.rooms, newRoom)
	if err := newRoom.addClient(c); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	} else {
		newRoom.setOpen()
		return c.String(http.StatusOK, "OK")
	}

	return c.String(http.StatusInternalServerError, "Error")
}

func (rm *RoomManager) RemoveRoom(room *Room) {
	newRooms := make([]*Room, 0, (len(rm.rooms) - 1))
	for _, r := range rm.rooms {
		if r != room {
			newRooms = append(newRooms, r)
		}
	}
	rm.rooms = newRooms
	
}
