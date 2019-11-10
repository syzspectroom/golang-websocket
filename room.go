package main

import (
	"fmt"
	"time"
)

const (
	//commands list
	disconnectCommand = 1

	roomTimeout = 60 * 5 * time.Second
	//Must be less than roomTimeout.
	livenessCheckTime = (roomTimeout * 9) / 10
)

type RoomControl struct {
	msg int
}

type Room struct {
	roomID string
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	broadcast  chan []byte
	control    chan *RoomControl
	timeout    time.Time
	clients    map[*Client]bool
}

func newRoom(roomID string) *Room {
	room := &Room{
		roomID:     roomID,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		control:    make(chan *RoomControl),
		timeout:    time.Now().Add(roomTimeout),
		clients:    make(map[*Client]bool),
	}

	go room.run()

	return room
}

func (r *Room) dropClients() {
	for client := range r.clients {
		r.removeClient(client)
	}
}

func (r *Room) removeClient(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		close(c.send)
	}
}

func (r *Room) extendTimeout() {
	r.timeout = time.Now().Add(roomTimeout)
}

func (r *Room) run() {
	defer r.dropClients()
	for {
		select {
		case client := <-r.register:
			fmt.Println("clinet registered")
			r.clients[client] = true
		case client := <-r.unregister:
			fmt.Println("clinet unregistered")
			r.removeClient(client)
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					r.removeClient(client)
				}
			}
		case control := <-r.control:
			switch control.msg {
			case disconnectCommand:
				return
			}
		}
	}
}
