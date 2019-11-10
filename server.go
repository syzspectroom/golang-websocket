package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	rooms map[string]*Room
	mux   sync.Mutex
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
		s.mux.Lock()
		delete(s.rooms, r.roomID)
		s.mux.Unlock()
	}()
	for {
		time.Sleep(livenessCheckTime)
		if time.Now().After(r.readTimeout()) {
			fmt.Println("livenessChecker - timeout")
			break
		}
	}
}

func (s *Server) roomsCount() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.rooms)
}

func (s *Server) addRoom(room *Room) {
	s.mux.Lock()
	fmt.Printf("adding new room: %+v\n", room.roomID)
	s.rooms[room.roomID] = room
	go s.livenessChecker(room)
	s.mux.Unlock()
}

func (s *Server) getRoomByID(roomID string) *Room {
	var room *Room
	s.mux.Lock()
	room = s.rooms[roomID]
	s.mux.Unlock()
	if room == nil {
		room = newRoom(roomID)
		s.addRoom(room)
	}

	return room
}

func (s *Server) addClientToRoom(room *Room, conn *websocket.Conn) {
	room.addClient(conn)
	room.extendTimeout()
}

func (s *Server) addClientToRoomByID(roomID string, conn *websocket.Conn) {
	room := s.getRoomByID(roomID)
	s.addClientToRoom(room, conn)
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
