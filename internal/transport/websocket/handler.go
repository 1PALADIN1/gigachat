package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", h.setupWsConnection)
}

func (h *Handler) setupWsConnection(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in connection:", err)
		return
	}

	go func() {
		defer connection.Close()
		h.handleUserMessages(connection)
	}()
}
