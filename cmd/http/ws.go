package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Hub is the central hub for the websocket connection,
// it manages the active clients list and broadcasts messages to all clients
type Hub struct {
	clients    map[*Client]struct{} // registered clients
	broadcast  chan []byte          // inbound messages from the clients
	register   chan *Client         // register requests from the clients
	unregister chan *Client         // unregister requests from the clients
	done       chan struct{}        // signal to close the hub
}

// Client is a middleman between the websocket connection and the hub.
// it handles the sending and receiving of messages from the clients
type Client struct {
	hub   *Hub            // the hub this client is registered to
	send  chan []byte     // buffered channel of outbound messages to the client
	conn  *websocket.Conn // the websocket connection
	ready chan struct{}   // signal that the client is ready to receive messages
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() { // cleanup when the function returns
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("failed to read message", "error", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	close(c.ready) // we are ready to receive messages
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for range n {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
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
