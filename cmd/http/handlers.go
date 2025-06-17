package http

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/oseau/web"
	"github.com/oseau/web/db"
)

// Handler holds the db
type Handler struct {
	db *db.DB
}

var once sync.Once

// NewHandler create the Handler
func NewHandler() *Handler {
	h := &Handler{}
	once.Do(func() { h.db = db.NewDB() })
	return h
}

type counter struct {
	Count int `json:"count"`
}

func (h *Handler) getCounter(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(counter{Count: h.db.GetCount().Count})
}

func (h *Handler) setCounter(w http.ResponseWriter, r *http.Request) {
	var c counter
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	}
	h.db.SetCount(c.Count)
}

func getVersion(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(web.VersionString))
}
