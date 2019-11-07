package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Server struct {
	rooms map[string]*Room
}

func newServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}

func (s *Server) addRoom(room *Room) {
	fmt.Printf("adding new room: %+v\n", room.roomID)
	s.rooms[room.roomID] = room
}

func (s *Server) getRoomByID(roomID string) *Room {
	var room *Room

	room = s.rooms[roomID]

	if room == nil {
		room = newRoom(roomID)
		s.addRoom(room)
	}

	return room
}

func (s *Server) addClientToRoom(room *Room, conn *websocket.Conn) {
	newClient(room, conn)
}

func (s *Server) addClientToRoomByID(roomID string, conn *websocket.Conn) {
	room := s.getRoomByID(roomID)
	s.addClientToRoom(room, conn)
}

func (s *Server) broadcastToRoomByID(roomID string, msg string) {
	room := server.getRoomByID(roomID)
	room.broadcast <- []byte(msg)
}
