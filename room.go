package main

import "fmt"

type Room struct {
	roomID string
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]bool
}

func newRoom(roomID string) *Room {
	room := &Room{
		roomID:     roomID,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go room.run()

	return room
}

func (h *Room) run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("clinet registered")
			h.clients[client] = true
		case client := <-h.unregister:
			fmt.Println("clinet unregistered")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					fmt.Println("clinet unregistered default")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
