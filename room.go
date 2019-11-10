package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	//commands list
	disconnectCommand = 1

	livenessCheckTime = 60 * 5 * time.Second
	//Must be less than livenessCheckTime.
	roomTimeout = (livenessCheckTime * 9) / 10
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
	mux        sync.Mutex
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
	r.mux.Lock()
	r.timeout = time.Now().Add(roomTimeout)
	r.mux.Unlock()
}

func (r *Room) readTimeout() time.Time {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.timeout
}

func (r *Room) addClient(conn *websocket.Conn) {
	newClient(r, conn)
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
