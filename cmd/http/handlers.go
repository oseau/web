package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oseau/web"
	"github.com/oseau/web/db"
	"github.com/oseau/web/redis"
)

// Handler holds the db
type Handler struct {
	db    *db.DB
	redis *redis.Redis
	hub   *Hub
}

// NewHandler create the Handler
func NewHandler(db *db.DB, redis *redis.Redis, hub *Hub) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
		hub:   hub,
	}
}

type counter struct {
	CountClick int `json:"click"`
	CountView  int `json:"view"`
}

func (h *Handler) getCounter(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(counter{CountClick: h.db.GetCount().Count, CountView: int(h.redis.AddViewCount())})
}

func (h *Handler) setCounter(w http.ResponseWriter, r *http.Request) {
	var c counter
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	}
	// for demo purposes, we just insert a new count record without any checks here
	// which is probably NOT what you want.
	h.db.SetCount(c.CountClick)
	json.NewEncoder(w).Encode(counter{CountClick: c.CountClick})
}

func getVersion(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(web.VersionString))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == URLFrontend
	},
}

func (h *Handler) ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade websocket connection", "error", err)
		return
	}
	client := &Client{hub: h.hub, conn: conn, send: make(chan []byte), ready: make(chan struct{})}
	go client.writePump()
	go client.readPump()
	<-client.ready // wait for the client to be ready to receive messages
	client.hub.register <- client
}
