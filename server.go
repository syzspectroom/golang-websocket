package main

import (
	"fmt"
	"time"

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

func (s *Server) livenessChecker(r *Room) {
	defer func() {
		r.control <- &RoomControl{msg: disconnectCommand}
		//TODO: check if it always deletes all clients before removing a room
		delete(s.rooms, r.roomID)
	}()
	for {
		time.Sleep(livenessCheckTime)
		if time.Now().After(r.timeout) {
			fmt.Println("livenessChecker - timeout")
			break
		}
	}
}

func (s *Server) addRoom(room *Room) {
	fmt.Printf("adding new room: %+v\n", room.roomID)
	s.rooms[room.roomID] = room
	go s.livenessChecker(room)
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
	room.extendTimeout()
}

func (s *Server) addClientToRoomByID(roomID string, conn *websocket.Conn) {
	room := s.getRoomByID(roomID)
	s.addClientToRoom(room, conn)
	room.extendTimeout()
}

func (s *Server) broadcastToRoomByID(roomID string, msg string) {
	room := server.getRoomByID(roomID)
	room.broadcast <- []byte(msg)
	room.extendTimeout()
}

func (s *Server) pingRoomByID(roomID string) {
	room := server.getRoomByID(roomID)
	room.extendTimeout()
}
