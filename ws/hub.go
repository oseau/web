// Package ws provides the websocket hub for the application.
package ws

import "encoding/json"

// Hub is the central hub for the websocket connection,
// it manages the active clients list and broadcasts messages to all clients
type Hub struct {
	clients    map[*Client]struct{} // registered clients
	broadcast  chan []byte          // inbound messages from the clients
	register   chan *Client         // register requests from the clients
	unregister chan *Client         // unregister requests from the clients
	done       chan struct{}        // signal to close the hub
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		done:       make(chan struct{}),
	}
}

// Run is the main loop for the hub, it handles the registration, unregistration, and broadcasting of messages
// this is running in a single seperated goroutine, so we don't need to worry about concurrent access to the hub inside this function
// this is the only place where we:
// read/write h.clients,
// read h.register, h.unregister, h.broadcast and h.done,
// write client.send
// we use this pattern to achieve a lock-free and concurrent-safe design
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			// client counts change, so we need to broadcast the message to all clients
			h.broadcastMsg(onlineCount{Count: len(h.clients)}.Bytes())
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// client counts change, so we need to broadcast the message to all clients
				h.broadcastMsg(onlineCount{Count: len(h.clients)}.Bytes())
			}
		case message := <-h.broadcast:
			h.broadcastMsg(message)
		case <-h.done:
			// the hub is closing, so we need to return
			return
		}
	}
}

func (h *Hub) broadcastMsg(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// this client is likely dead from client side, so we close the connection
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// Close closes the hub
func (h *Hub) Close() error {
	close(h.done) // signal to stop the Run goroutine
	close(h.register)
	close(h.unregister)
	close(h.broadcast)
	return nil
}

type onlineCount struct {
	Count int `json:"count"`
}

func (o onlineCount) Bytes() []byte {
	bytes, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return bytes
}
